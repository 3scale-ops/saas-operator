package redis

import (
	"context"
	"fmt"
	"strings"

	redistypes "github.com/3scale/saas-operator/pkg/redis/types"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
)

// Server represent a redis server and its characteristics
type Server struct {
	Role         redistypes.Role
	ReadOnly     bool
	ClientConfig *redis.Options
}

// GetIP returns the IP of the redis server
func (srv *Server) GetIP() string {
	parts := strings.Split(srv.ClientConfig.Addr, ":")
	return parts[0]
}

// GetPort returns the port of the redis server
func (srv *Server) GetPort() string {
	parts := strings.Split(srv.ClientConfig.Addr, ":")
	return parts[1]
}

// GetClient returns a redis client, ready to talk to this
// redis server
func (srv *Server) GetClient() *redis.Client {
	return redis.NewClient(srv.ClientConfig)
}

// getRole retrieves the current role of a redis Server within the shard
func (srv *Server) getRole(ctx context.Context) (redistypes.Role, error) {
	val, err := srv.GetClient().Do(ctx, "role").Result()
	if err != nil {
		return redistypes.Unknown, err
	}
	return redistypes.Role(val.([]interface{})[0].(string)), nil
}

// isReadOnly checks whether the redis Server has ReadOnly flag or not
func (srv *Server) isReadOnly(ctx context.Context) (bool, error) {
	val, err := srv.GetClient().ConfigGet(ctx, "slave-read-only").Result()
	if err != nil {
		return false, err
	}

	return (val[1].(string)) == "yes", nil
}

// Discover returns the Role and the IsReadOnly flag for a given
// redis Server
func (srv *Server) Discover(ctx context.Context) error {

	role, err := srv.getRole(ctx)
	if err != nil {
		return util.WrapError("[redis-autodiscovery/Server.GetRole]", err)
	}
	srv.Role = role

	if srv.Role == redistypes.Slave {
		ro, err := srv.isReadOnly(ctx)
		if err != nil {
			return util.WrapError("[redis-autodiscovery/Server.IsReadOnly]", err)
		}
		srv.ReadOnly = ro
	}
	return nil
}

// Shard is a list of the redis Server objects that compose a redis shard
type Shard []Server

// NewShard returns a Shard object given the passed redis server URLs
func NewShard(connectionStrings []string) (Shard, error) {
	shard := Shard{}
	for _, cs := range connectionStrings {
		opt, err := redis.ParseURL(cs)
		if err != nil {
			return nil, util.WrapError("[redis-autodiscovery/NewShard]", err)
		}
		shard = append(shard, Server{ClientConfig: opt, Role: redistypes.Unknown, ReadOnly: false})
	}
	return shard, nil
}

// Discover retrieves the role and read-only flag for all the server in the shard
func (s Shard) Discover(ctx context.Context, log logr.Logger) error {

	for idx := range s {
		if err := s[idx].Discover(ctx); err != nil {
			return err
		}
	}

	masters := 0
	for _, server := range s {
		if server.Role == redistypes.Master {
			masters++
		}
	}

	if masters != 1 {
		err := fmt.Errorf("[redis-autodiscovery/Shard.Discover] expected 1 master but got %d", masters)
		log.Error(err, "error discovering shard server roles")
		return err
	}

	return nil
}

// GetMasterAddr returns the URL of the master server in a shard or error if zero
// or more than one master is found
func (s Shard) GetMasterAddr() (string, string, error) {
	for _, server := range s {
		if server.Role == redistypes.Master {
			parts := strings.Split(server.ClientConfig.Addr, ":")
			return parts[0], parts[1], nil
		}
	}
	return "", "", fmt.Errorf("[redis-autodiscovery/Shard.GetMasterAddr] master not found")
}

// Init initializes this shard if not already initialized
func (s Shard) Init(ctx context.Context, masterIndex int32, log logr.Logger) error {

	for idx, srv := range s {
		val, err := srv.GetClient().Do(ctx, "role").Result()
		if err != nil {
			return err
		}

		role := val.([]interface{})[0].(string)
		if role == string(redistypes.Slave) {

			slaveof := val.([]interface{})[1].(string)
			if slaveof == "127.0.0.1" {

				if idx == int(masterIndex) {
					_, err := srv.GetClient().SlaveOf(ctx, "NO", "ONE").Result()
					if err != nil {
						return err
					}
					log.Info(fmt.Sprintf("[@redis-setup] Configured %s as master", srv.GetIP()))
				} else {
					_, err := srv.GetClient().SlaveOf(ctx, s[masterIndex].GetIP(), s[masterIndex].GetPort()).Result()
					if err != nil {
						return err
					}
					log.Info(fmt.Sprintf("[@redis-setup] Configured %s as slave", srv.GetIP()))
				}

			} else {
				s[idx].Role = redistypes.Slave
			}

		} else if role == string(redistypes.Master) {
			s[idx].Role = redistypes.Master
		} else {
			return fmt.Errorf("[@redis-setup] unable to get role for server %s:%s", srv.GetIP(), srv.GetPort())
		}
	}

	return nil
}

// ShardedCluster represents a sharded redis cluster, composed by several Shards
type ShardedCluster map[string]Shard

// NewShardedCluster returns a new ShardedCluster given the shard structure passed as a map[string][]string
func NewShardedCluster(ctx context.Context, serverList map[string][]string, log logr.Logger) (ShardedCluster, error) {

	sc := make(map[string]Shard, len(serverList))

	for shardName, shardServers := range serverList {

		shard, err := NewShard(shardServers)
		if err != nil {
			return nil, err
		}

		if err := shard.Discover(ctx, log.WithValues("ShardName", shardName)); err != nil {
			return nil, err
		}

		sc[shardName] = shard
	}

	return sc, nil
}

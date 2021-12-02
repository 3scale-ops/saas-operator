package redis

import (
	"context"
	"fmt"

	"github.com/3scale/saas-operator/pkg/redis/crud"
	"github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
)

// RedisServer represent a redis server and its characteristics
type RedisServer struct {
	Name     string
	Role     client.Role
	ReadOnly bool
	CRUD     *crud.CRUD
}

func NewRedisServer(name, connectionString string) (*RedisServer, error) {

	crud, err := crud.NewRedisCRUD(connectionString)
	if err != nil {
		return nil, err
	}

	return &RedisServer{Name: connectionString, CRUD: crud, ReadOnly: false}, nil
}

// Discover returns the Role and the IsReadOnly flag for a given
// redis Server
func (srv *RedisServer) Discover(ctx context.Context) error {

	role, _, err := srv.CRUD.RedisRole(ctx)
	if err != nil {
		return util.WrapError("[redis-autodiscovery]", err)
	}
	srv.Role = role

	if srv.Role == client.Slave {
		ro, err := srv.CRUD.RedisConfigGet(ctx, "slave-read-only")
		if err != nil {
			return util.WrapError("[redis-autodiscovery]", err)
		}
		if ro == "yes" {
			srv.ReadOnly = true
		}
	}
	return nil
}

// Shard is a list of the redis Server objects that compose a redis shard
type Shard []RedisServer

// NewShard returns a Shard object given the passed redis server URLs
func NewShard(connectionStrings []string) (Shard, error) {
	shard := make([]RedisServer, len(connectionStrings))
	for i, cs := range connectionStrings {
		rs, err := NewRedisServer(cs, cs)
		if err != nil {
			return nil, err
		}
		shard[i] = *rs
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
		if server.Role == client.Master {
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
	for _, srv := range s {
		if srv.Role == client.Master {
			return srv.CRUD.GetIP(), srv.CRUD.GetPort(), nil
		}
	}
	return "", "", fmt.Errorf("[redis-autodiscovery/Shard.GetMasterAddr] master not found")
}

// Init initializes this shard if not already initialized
func (s Shard) Init(ctx context.Context, masterIndex int32, log logr.Logger) error {

	for idx, srv := range s {
		role, slaveof, err := srv.CRUD.RedisRole(ctx)
		if err != nil {
			return err
		}

		if role == client.Slave {

			if slaveof == "127.0.0.1" {

				if idx == int(masterIndex) {
					if err := srv.CRUD.RedisSlaveOf(ctx, "NO", "ONE"); err != nil {
						return err
					}
					log.Info(fmt.Sprintf("[@redis-setup] Configured %s as master", srv.Name))
				} else {
					if err := srv.CRUD.RedisSlaveOf(ctx, s[masterIndex].CRUD.GetIP(), s[masterIndex].CRUD.GetPort()); err != nil {
						return err
					}
					log.Info(fmt.Sprintf("[@redis-setup] Configured %s as slave", srv.Name))
				}

			} else {
				s[idx].Role = client.Slave
			}

		} else if role == client.Master {
			s[idx].Role = client.Master
		} else {
			return fmt.Errorf("[@redis-setup] unable to get role for server %s", srv.Name)
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

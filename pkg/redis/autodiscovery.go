package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
)

type RedisRole string

const (
	Master  RedisRole = "master"
	Slave   RedisRole = "slave"
	Unknown RedisRole = "unknown"
)

type Server struct {
	Role         RedisRole
	ReadOnly     bool
	ClientConfig *redis.Options
}

func (srv *Server) GetRole(ctx context.Context) (RedisRole, error) {
	rdb := redis.NewClient(srv.ClientConfig)
	val, err := rdb.Do(ctx, "role").Result()
	if err != nil {
		return Unknown, err
	}
	return RedisRole(val.([]interface{})[0].(string)), nil
}

func (srv *Server) IsReadOnly(ctx context.Context) (bool, error) {
	rdb := redis.NewClient(srv.ClientConfig)
	val, err := rdb.ConfigGet(ctx, "slave-read-only").Result()
	if err != nil {
		return false, err
	}

	return (val[1].(string)) == "yes", nil
}

func (srv *Server) Discover(ctx context.Context) error {

	role, err := srv.GetRole(ctx)
	if err != nil {
		return util.WrapError("[redis-autodiscovery/Server.GetRole]", err)
	}
	srv.Role = role

	if srv.Role == Slave {
		ro, err := srv.IsReadOnly(ctx)
		if err != nil {
			return util.WrapError("[redis-autodiscovery/Server.IsReadOnly]", err)
		}
		srv.ReadOnly = ro
	}
	return nil
}

type Shard []Server

func NewShard(connectionStrings []string) (Shard, error) {
	shard := Shard{}
	for _, cs := range connectionStrings {
		opt, err := redis.ParseURL(cs)
		if err != nil {
			return nil, util.WrapError("[redis-autodiscovery/NewShard]", err)
		}
		shard = append(shard, Server{ClientConfig: opt, Role: Unknown, ReadOnly: false})
	}
	return shard, nil
}

func (s Shard) Discover(ctx context.Context, log logr.Logger) error {

	for idx := range s {
		if err := s[idx].Discover(ctx); err != nil {
			return err
		}
	}

	masters := 0
	for _, server := range s {
		if server.Role == Master {
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

func (s Shard) GetMasterAddr() (string, string, error) {
	for _, server := range s {
		if server.Role == Master {
			parts := strings.Split(server.ClientConfig.Addr, ":")
			return parts[0], parts[1], nil
		}
	}
	return "", "", fmt.Errorf("[redis-autodiscovery/Shard.GetMasterAddr] master not found")
}

type ShardedCluster map[string]Shard

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

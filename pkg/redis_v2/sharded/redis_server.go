package sharded

import (
	"github.com/3scale/saas-operator/pkg/redis_v2/client"
	redis "github.com/3scale/saas-operator/pkg/redis_v2/server"
)

type RedisServer struct {
	*redis.Server
	Role   client.Role
	Config map[string]string
}

func NewRedisServerFromPool(connectionString string, alias *string, pool *redis.ServerPool) (*RedisServer, error) {
	srv, err := pool.GetServer(connectionString, alias)
	if err != nil {
		return nil, err
	}

	return &RedisServer{
		Server: srv,
		Role:   client.Unknown,
		Config: map[string]string{},
	}, nil
}

func NewRedisServerFromParams(srv *redis.Server, role client.Role, config map[string]string) *RedisServer {
	return &RedisServer{
		Server: srv,
		Role:   role,
		Config: config,
	}
}

package sharded

import (
	"context"
	"fmt"

	"github.com/3scale/saas-operator/pkg/redis/client"
	redis "github.com/3scale/saas-operator/pkg/redis/server"
	"sigs.k8s.io/controller-runtime/pkg/log"
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

func (srv *RedisServer) InitMaster(ctx context.Context) (bool, error) {
	logger := log.FromContext(ctx, "function", "(*RedisServer).InitMaster")

	role, slaveof, err := srv.RedisRole(ctx)
	if err != nil {
		return false, err
	}

	switch role {
	case client.Slave:

		if slaveof == "127.0.0.1" {
			// needs initialization
			if err := srv.RedisSlaveOf(ctx, "NO", "ONE"); err != nil {
				return false, err
			}
			logger.Info(fmt.Sprintf("configured %s|%s as master", srv.GetAlias(), srv.ID()))
			return true, nil

		} else {
			srv.Role = client.Slave
		}

	case client.Master:
		srv.Role = client.Master
	}

	return false, nil
}

func (srv *RedisServer) InitSlave(ctx context.Context, master *RedisServer) (bool, error) {

	logger := log.FromContext(ctx, "function", "(*RedisServer).InitSlave")

	role, slaveof, err := srv.RedisRole(ctx)
	if err != nil {
		return false, err
	}

	switch role {
	case client.Slave:

		// needs initialization
		if slaveof == "127.0.0.1" {
			// validate first that the master is ready
			role, _, err := master.RedisRole(ctx)
			if err != nil || role != client.Master {
				err := fmt.Errorf("shard master %s|%s is not ready", master.GetAlias(), master.ID())
				logger.Error(err, "slave init failed")
				return false, err

			} else {
				// if master ok, init slave
				if err := srv.RedisSlaveOf(ctx, master.GetHost(), master.GetPort()); err != nil {
					return false, err
				}
				logger.Info(fmt.Sprintf("configured %s|%s as slave", srv.GetAlias(), srv.ID()))
				return true, nil
			}

		} else {
			srv.Role = client.Slave
			// FOR DEBUGGING
			// val, err := srv.GetClient().RedisDo(ctx, "info", "replication")
			// if err != nil {
			// 	logger.Error(err, "unable to get info")
			// } else {
			// 	logger.Info("dump replication status", "Slave", srv.GetAlias())
			// 	logger.Info(fmt.Sprintf("%s", redis.InfoStringToMap(val.(string))))
			// }
		}

	case client.Master:
		srv.Role = client.Master
	}

	return false, nil
}

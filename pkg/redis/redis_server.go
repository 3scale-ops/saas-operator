package redis

import (
	"context"
	"fmt"
	"net"

	"github.com/3scale/saas-operator/pkg/redis/crud"
	"github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
)

// RedisServer represent a redis server and its characteristics
type RedisServer struct {
	Name string
	Host string
	Port string
	Role client.Role
	CRUD *crud.CRUD
}

func (rs *RedisServer) IP() (string, error) {
	var ip string
	if r := net.ParseIP(rs.Host); r != nil {
		ip = r.String()
	} else {
		// if it is not an IP, try to resolve a DNS
		ips, err := net.LookupIP(rs.Host)
		if err != nil {
			return "", err
		}
		if len(ips) > 1 {
			return "", fmt.Errorf("dns resolves to more than 1 IP")
		}
		ip = ips[0].String()
	}

	return ip, nil
}

func NewRedisServerFromConnectionString(name, connectionString string) (*RedisServer, error) {

	crud, err := crud.NewRedisCRUDFromConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	return &RedisServer{Name: name, Host: crud.GetIP(), Port: crud.GetPort(), CRUD: crud, Role: client.Unknown}, nil
}

// Cleanup closes all Redis clients opened during the RedisServer object creation
func (srv *RedisServer) Cleanup(log logr.Logger) error {
	log.V(2).Info("[@redis-server-cleanup] closing client",
		"server", srv.Name, "host", srv.Host,
	)
	if err := srv.CRUD.CloseClient(); err != nil {
		log.Error(err, "[@redis-server-cleanup] error closing server client",
			"server", srv.Name, "host", srv.Host,
		)
		return err
	}
	return nil
}

// Discover returns the Role for a given
// redis Server
func (srv *RedisServer) Discover(ctx context.Context) error {

	role, _, err := srv.CRUD.RedisRole(ctx)
	if err != nil {
		srv.Role = client.Unknown
		return util.WrapError("redis-autodiscovery", err)
	}
	srv.Role = role

	return nil
}

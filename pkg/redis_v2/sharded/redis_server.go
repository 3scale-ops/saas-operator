package sharded

import (
	"context"
	"fmt"

	"github.com/3scale/saas-operator/pkg/redis_v2/client"
	redis "github.com/3scale/saas-operator/pkg/redis_v2/server"
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

type DiscoveryOptionSet []DiscoveryOption

type DiscoveryOption int

const (
	SlaveReadOnlyDiscoveryOpt DiscoveryOption = iota
	SaveConfigDiscoveryOpt
)

func (set DiscoveryOptionSet) Has(opt DiscoveryOption) bool {
	for _, o := range set {
		if opt == o {
			return true
		}
	}
	return false
}

// Discover returns the characteristincs for a given
// redis Server
// It always gets the role first
func (srv *RedisServer) Discover(ctx context.Context, opts ...DiscoveryOption) error {
	logger := log.FromContext(ctx, "function", "(*RedisServer).DiscoverWithOptions()")

	role, _, err := srv.RedisRole(ctx)
	if err != nil {
		srv.Role = client.Unknown
		logger.Error(err, fmt.Sprintf("unable to get %s|%s|%s role", srv.GetAlias(), srv.Role, srv.ID()))
		return err
	}
	srv.Role = role

	if srv.Config == nil {
		srv.Config = map[string]string{}
	}

	if DiscoveryOptionSet(opts).Has(SaveConfigDiscoveryOpt) {

		save, err := srv.RedisConfigGet(ctx, "save")
		if err != nil {
			logger.Error(err, fmt.Sprintf("unable to get %s|%s|%s 'save' option", srv.GetAlias(), srv.Role, srv.ID()))
			return err
		}
		srv.Config["save"] = save
	}

	if DiscoveryOptionSet(opts).Has(SlaveReadOnlyDiscoveryOpt) && role != client.Master {
		slaveReadOnly, err := srv.RedisConfigGet(ctx, "slave-read-only")
		if err != nil {
			logger.Error(err, fmt.Sprintf("unable to get %s|%s|%s 'slave-read-only' option", srv.GetAlias(), srv.Role, srv.ID()))
			return err
		}
		srv.Config["slave-read-only"] = slaveReadOnly
	}

	return nil
}

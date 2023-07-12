package server

import (
	"context"
	"fmt"

	"github.com/3scale/saas-operator/pkg/redis_v2/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DiscoveryOptionSet []DiscoveryOption

type DiscoveryOption int

const (
	SlaveReadOnlyDiscoveryOpt DiscoveryOption = iota
	SaveConfigDiscoveryOpt
	RoleDiscoveryOpt
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
func (srv *Server) Discover(ctx context.Context, opts ...DiscoveryOption) error {
	logger := log.FromContext(ctx, "function", "(*RedisServer).DiscoverWithOptions()")

	if DiscoveryOptionSet(opts).Has(RoleDiscoveryOpt) {
		role, _, err := srv.RedisRole(ctx)
		if err != nil {
			srv.Role = client.Unknown
			logger.Error(err, fmt.Sprintf("unable to get %s|%s|%s:%s role", srv.Alias, srv.Role, srv.Host, srv.Port))
			return err
		}
		srv.Role = role
	}

	if len(opts) > 0 && srv.Config == nil {
		srv.Config = map[string]string{}
	}

	if DiscoveryOptionSet(opts).Has(SaveConfigDiscoveryOpt) {

		save, err := srv.RedisConfigGet(ctx, "save")
		if err != nil {
			logger.Error(err, fmt.Sprintf("unable to get %s|%s|%s:%s 'save' option", srv.Alias, srv.Role, srv.Host, srv.Port))
			return err
		}
		srv.Config["save"] = save
	}

	if DiscoveryOptionSet(opts).Has(SlaveReadOnlyDiscoveryOpt) {
		slaveReadOnly, err := srv.RedisConfigGet(ctx, "slave-read-only")
		if err != nil {
			logger.Error(err, fmt.Sprintf("unable to get %s|%s|%s:%s 'slave-read-only' option", srv.Alias, srv.Role, srv.Host, srv.Port))
			return err
		}
		srv.Config["slave-read-only"] = slaveReadOnly
	}

	return nil
}

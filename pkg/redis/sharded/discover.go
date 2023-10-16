package sharded

import (
	"context"
	"fmt"

	"github.com/3scale/saas-operator/pkg/redis/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DiscoveryOptionSet []DiscoveryOption

type DiscoveryOption int

const (
	SlaveReadOnlyDiscoveryOpt DiscoveryOption = iota
	SaveConfigDiscoveryOpt
	OnlyMasterDiscoveryOpt
	SlavePriorityDiscoveryOpt
	ReplicationInfoDiscoveryOpt
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
	logger := log.FromContext(ctx, "function", "(*RedisServer).Discover()")

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

	if DiscoveryOptionSet(opts).Has(SlavePriorityDiscoveryOpt) {
		slavePriority, err := srv.RedisConfigGet(ctx, "slave-priority")
		if err != nil {
			logger.Error(err, fmt.Sprintf("unable to get %s|%s|%s 'slave-priority' option", srv.GetAlias(), srv.Role, srv.ID()))
			return err
		}
		srv.Config["slave-priority"] = slavePriority
	}

	if DiscoveryOptionSet(opts).Has(ReplicationInfoDiscoveryOpt) && role != client.Master {
		repinfo, err := srv.RedisInfo(ctx, "replication")
		if err != nil {
			logger.Error(err, fmt.Sprintf("unable to get %s|%s|%s replication info", srv.GetAlias(), srv.Role, srv.ID()))
			return err
		}

		var syncInProgress string
		switch flag := repinfo["master_sync_in_progress"]; flag {
		case "0":
			syncInProgress = "no"
		case "1":
			syncInProgress = "yes"
		default:
			logger.Error(err, fmt.Sprintf("unexpected value '%s' for 'master_sync_in_progress' %s|%s|%s", flag, srv.GetAlias(), srv.Role, srv.ID()))
			syncInProgress = ""
		}

		if srv.Info == nil {
			srv.Info = map[string]string{}
		}
		srv.Info["replication"] = fmt.Sprintf("master-link: %s, sync-in-progress: %s", repinfo["master_link_status"], syncInProgress)
	}

	return nil
}

// Discovery errors
type DiscoveryError_Sentinel_Failure struct{ error }
type DiscoveryError_Master_SingleServerFailure struct{ error }
type DiscoveryError_Slave_FailoverInProgress struct{ error }
type DiscoveryError_Slave_SingleServerFailure struct{ error }
type DiscoveryError_UnknownRole_SingleServerFailure struct{ error }

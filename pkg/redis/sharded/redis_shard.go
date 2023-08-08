package sharded

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/3scale/saas-operator/pkg/redis/client"
	redis "github.com/3scale/saas-operator/pkg/redis/server"
	"github.com/3scale/saas-operator/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Shard is a list of the redis Server objects that compose a redis shard
type Shard struct {
	Name    string
	Servers []*RedisServer
	pool    *redis.ServerPool
}

func NewShardFromServers(name string, pool *redis.ServerPool, servers ...*RedisServer) *Shard {
	shard := &Shard{}
	shard.Name = name
	shard.Servers = append(shard.Servers, servers...)
	shard.pool = pool
	return shard
}

// NewShardFromTopology returns a Shard object given the passed redis server URLs
func NewShardFromTopology(name string, servers map[string]string, pool *redis.ServerPool) (*Shard, error) {
	var merr util.MultiError
	shard := &Shard{Name: name, pool: pool}
	shard.Servers = make([]*RedisServer, 0, len(servers))

	for key, connectionString := range servers {
		var alias *string = nil
		if key != connectionString {
			alias = &key
		}
		srv, err := NewRedisServerFromPool(connectionString, alias, pool)
		if err != nil {
			merr = append(merr, err)
			continue
		}
		shard.Servers = append(shard.Servers, srv)
	}

	// sort the slice to obtain consistent results
	sort.Slice(shard.Servers, func(i, j int) bool {
		return shard.Servers[i].ID() < shard.Servers[j].ID()
	})

	return shard, merr.ErrorOrNil()
}

// Discover retrieves the options for all the servers in the shard
// If a SentinelServer is provided, it will be used to autodiscover servers and roles in the shard
func (shard *Shard) Discover(ctx context.Context, sentinel *SentinelServer, options ...DiscoveryOption) error {
	var merr util.MultiError
	logger := log.FromContext(ctx, "function", "(*Shard).Discover", "shard", shard.Name)

	switch sentinel {

	// no sentinel provided
	case nil:
		for idx := range shard.Servers {
			if err := shard.Servers[idx].Discover(ctx, options...); err != nil {
				logger.Error(err, fmt.Sprintf("unable to discover redis server %s", shard.Servers[idx].ID()))
				merr = append(merr, DiscoveryError_UnknownRole_SingleServerFailure{err})
				continue
			}
		}

	// sentinel provided
	default:
		sentinelMasterResult, err := sentinel.SentinelMaster(ctx, shard.Name)
		if err != nil {
			return append(merr, DiscoveryError_Sentinel_Failure{err})
		}

		// Get the corresponding server or add a new one if not found
		srv, err := shard.GetServerByID(fmt.Sprintf("%s:%d", sentinelMasterResult.IP, sentinelMasterResult.Port))
		if err != nil {
			return append(merr, DiscoveryError_Master_SingleServerFailure{err})
		}

		// do not try to discover a master flagged as "s_down" or "o_down"
		if strings.Contains(sentinelMasterResult.Flags, "s_down") || strings.Contains(sentinelMasterResult.Flags, "o_down") {
			err := fmt.Errorf("master %s is s_down/o_down", srv.GetAlias())
			logger.Error(err, "master down")
			return append(merr, DiscoveryError_Master_SingleServerFailure{err})
		}

		// Confirm the server role
		if err = srv.Discover(ctx, options...); err != nil {
			srv.Role = client.Role(client.Unknown)
			return append(merr, DiscoveryError_Master_SingleServerFailure{err})
		} else if srv.Role != client.Master {
			// the role that the server reports is different from the role that
			// sentinel sees. Probably the sentinel configuration hasn't converged yet
			// this is an error and should be retried
			srv.Role = client.Role(client.Unknown)
			return append(merr, DiscoveryError_Master_SingleServerFailure{fmt.Errorf("sentinel config has not yet converged for %s", srv.GetAlias())})
		}

		if DiscoveryOptionSet(options).Has(OnlyMasterDiscoveryOpt) {
			return merr.ErrorOrNil()
		}

		// discover slaves
		sentinelSlavesResult, err := sentinel.SentinelSlaves(ctx, shard.Name)
		if err != nil {
			return append(merr, DiscoveryError_Sentinel_Failure{err})
		}
		for _, slave := range sentinelSlavesResult {

			// Get the corresponding server or add a new one if not found
			srv, err := shard.GetServerByID(fmt.Sprintf("%s:%d", slave.IP, slave.Port))
			if err != nil {
				merr = append(merr, DiscoveryError_Slave_SingleServerFailure{err})
				continue
			}

			// do not try to discover a slave flagged as "s_down" or "o_down"
			if strings.Contains(slave.Flags, "s_down") || strings.Contains(slave.Flags, "o_down") {
				err := fmt.Errorf("slave %s is s_down/o_down", srv.GetAlias())
				log.Log.Error(err, "slave is down")
				merr = append(merr, DiscoveryError_Slave_SingleServerFailure{err})
				continue

			} else {
				if err := srv.Discover(ctx, options...); err != nil {
					srv.Role = client.Role(client.Unknown)
					logger.Error(err, fmt.Sprintf("unable to discover redis server %s", srv.GetAlias()))
					merr = append(merr, DiscoveryError_Slave_SingleServerFailure{err})
					continue
				}
				if srv.Role != client.Slave {
					// the role that the server reports is different from the role that
					// sentinel sees. Probably the sentinel configuration hasn't converged yet
					// this is an error and should be retried
					srv.Role = client.Role(client.Unknown)
					merr = append(merr, DiscoveryError_Slave_SingleServerFailure{fmt.Errorf("sentinel config has not yet converged for %s", srv.GetAlias())})
					continue
				}
			}
		}
	}

	return merr.ErrorOrNil()
}

// GetMaster returns the host:port of the master server
// in a shard or error if zero or more than one master is found
func (shard *Shard) GetMaster() (*RedisServer, error) {
	master := []*RedisServer{}

	for _, srv := range shard.Servers {
		if srv.Role == client.Master {
			master = append(master, srv)
		}
	}

	if len(master) != 1 {
		return nil, util.WrapError("(*Shard).GetMasterAddr", fmt.Errorf("wrong number of masters: %d != 1", len(master)))
	}

	return master[0], nil
}

func (shard *Shard) GetSlavesRW() []*RedisServer {
	servers := []*RedisServer{}
	for _, srv := range shard.Servers {
		if srv.Role == client.Slave {
			if val, ok := srv.Config["slave-read-only"]; ok && val == "no" {
				servers = append(servers, srv)
			}
		}
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].ID() < servers[j].ID()
	})
	return servers
}

func (shard *Shard) GetSlavesRO() []*RedisServer {
	servers := []*RedisServer{}
	for _, srv := range shard.Servers {
		if srv.Role == client.Slave {
			if val, ok := srv.Config["slave-read-only"]; ok && val == "yes" {
				servers = append(servers, srv)
			}
		}
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].ID() < servers[j].ID()
	})
	return servers
}

func (shard *Shard) GetServerByID(hostport string) (*RedisServer, error) {
	var rs *RedisServer
	var err error

	for _, srv := range shard.Servers {
		if srv.ID() == hostport {
			rs = srv
			break
		}
	}

	// If the server is not in the list, add a new one
	if rs == nil {
		rs, err = NewRedisServerFromPool("redis://"+hostport, nil, shard.pool)
		if err != nil {
			return nil, err
		}
		shard.Servers = append(shard.Servers, rs)
	}

	return rs, nil
}

// Init initializes the shard if not already initialized
func (shard *Shard) Init(ctx context.Context, masterHostPort string) ([]string, error) {
	merr := util.MultiError{}
	listChanged := []string{}
	var master *RedisServer

	// Init the master
	for _, srv := range shard.Servers {
		if srv.ID() == masterHostPort {
			master = srv
			changed, err := master.InitMaster(ctx)
			if err != nil {
				return listChanged, append(merr, err)
			}
			if changed {
				listChanged = append(listChanged, master.ID())
			}
		}
	}

	// Init the slaves
	for _, srv := range shard.Servers {
		if srv.ID() != masterHostPort {
			changed, err := srv.InitSlave(ctx, master)
			if err != nil {
				merr = append(merr, err)
			}
			if changed {
				listChanged = append(listChanged, srv.ID())
			}
		}
	}

	return listChanged, merr.ErrorOrNil()
}

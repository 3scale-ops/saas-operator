package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/redis/crud"
	redis_client "github.com/3scale/saas-operator/pkg/redis/crud/client"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	shardNotInitializedError = "ERR No such master with that name"
	maxInfoCacheAge          = 10 * time.Second
)

// SentinelServer represents a sentinel Pod
type SentinelServer struct {
	Name string
	IP   string
	Port string
	CRUD *crud.CRUD
}

func NewSentinelServerFromConnectionString(name, connectionString string) (*SentinelServer, error) {

	crud, err := crud.NewRedisCRUDFromConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	return &SentinelServer{Name: name, IP: crud.GetIP(), Port: crud.GetPort(), CRUD: crud}, nil
}

// Cleanup closes all Redis clients opened during the SentinelServer object creation
func (ss *SentinelServer) Cleanup(log logr.Logger) error {
	log.V(2).Info("[@sentinel-server-cleanup] closing client",
		"server", ss.Name, "host", ss.IP,
	)
	if err := ss.CRUD.CloseClient(); err != nil {
		log.Error(err, "[@sentinel-server-cleanup] error closing server client",
			"server", ss.Name, "host", ss.IP,
		)
		return err
	}
	return nil
}

// IsMonitoringShards checks whether or all the shards in the passed list are being monitored by the SentinelServer
func (ss *SentinelServer) IsMonitoringShards(ctx context.Context, shards []string) (bool, error) {

	monitoredShards, err := ss.CRUD.SentinelMasters(ctx)
	if err != nil {
		return false, err
	}

	if len(monitoredShards) == 0 {
		return false, nil
	}

	for _, name := range shards {
		found := false
		for _, monitored := range monitoredShards {
			if monitored.Name == name {
				found = true
			}
		}
		if !found {
			return false, nil
		}
	}

	return true, nil
}

// Monitor ensures that all the shards in the ShardedCluster object are monitored by the SentinelServer
func (ss *SentinelServer) Monitor(ctx context.Context, shards ShardedCluster) ([]string, error) {
	changed := []string{}

	// Initialize unmonitored shards
	shardNames := shards.GetShardNames()
	for _, name := range shardNames {

		_, err := ss.CRUD.SentinelMaster(ctx, name)
		if err != nil {
			if err.Error() == shardNotInitializedError {

				shard := shards.GetShardByName(name)
				host, port, err := shard.GetMasterAddr()
				if err != nil {
					return changed, err
				}

				err = ss.CRUD.SentinelMonitor(ctx, name, host, port, saasv1alpha1.SentinelDefaultQuorum)
				if err != nil {
					return changed, util.WrapError("redis-sentinel/SentinelServer.Monitor", err)
				}
				// even if the next call fails, there has already been a write operation to sentinel
				changed = append(changed, name)

				err = ss.CRUD.SentinelSet(ctx, name, "down-after-milliseconds", "5000")
				if err != nil {
					return changed, util.WrapError("redis-sentinel/SentinelServer.Monitor", err)
				}
				// TODO: change the default failover timeout.
				// TODO: maybe add a generic mechanism to set/modify parameters

			} else {
				return changed, err
			}
		}
	}

	return changed, nil
}

type ShardDiscoveryOption int

const (
	OnlyMasterDiscoveryOpt ShardDiscoveryOption = iota
	SlaveReadOnlyDiscoveryOpt
	SaveConfigDiscoveryOpt
)

func (sdos ShardDiscoveryOptions) Has(sdo ShardDiscoveryOption) bool {
	for _, opt := range sdos {
		if opt == sdo {
			return true
		}
	}
	return false
}

type ShardDiscoveryOptions []ShardDiscoveryOption

// MonitoredShards returns the list of monitored shards of this SentinelServer
func (ss *SentinelServer) MonitoredShards(ctx context.Context, options ...ShardDiscoveryOption) (saasv1alpha1.MonitoredShards, error) {
	opts := ShardDiscoveryOptions(options)

	sm, err := ss.CRUD.SentinelMasters(ctx)
	if err != nil {
		return nil, err
	}

	monitoredShards := make([]saasv1alpha1.MonitoredShard, 0, len(sm))
	for _, s := range sm {

		var servers map[string]saasv1alpha1.RedisServerDetails
		servers, err = ss.DiscoverShard(ctx, s.Name, maxInfoCacheAge, opts)
		if err != nil {
			return nil, err
		}
		monitoredShards = append(monitoredShards,
			saasv1alpha1.MonitoredShard{
				Name:    s.Name,
				Servers: servers,
			},
		)
	}
	return monitoredShards, nil
}

func serverName(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

func connectionString(ip string, port int) string {
	return fmt.Sprintf("redis://%s:%d", ip, port)
}

func (ss *SentinelServer) DiscoverShard(ctx context.Context, shard string, maxInfoCacheAge time.Duration,
	opts ShardDiscoveryOptions) (map[string]saasv1alpha1.RedisServerDetails, error) {

	logger := log.FromContext(ctx, "function", "(*SentinelServer).DiscoverShard()")

	/////////////////////////////////
	// discover the shard's master //
	/////////////////////////////////

	master, err := ss.CRUD.SentinelMaster(ctx, shard)
	if err != nil {
		logger.Error(err, fmt.Sprintf("unable to get master for shard %s", shard))
		return nil, err
	}

	sn := serverName(master.IP, master.Port)

	// do not try to discover a master flagged as "s_down" or "o_down"
	if strings.Contains(master.Flags, "s_down") && strings.Contains(master.Flags, "o_down") {
		return nil, fmt.Errorf("%s master %s is s_down/o_down", shard, sn)
	}

	result := map[string]saasv1alpha1.RedisServerDetails{
		sn: {
			Role:   redis_client.Master,
			Config: map[string]string{},
		},
	}

	if opts.Has(SaveConfigDiscoveryOpt) {

		// open a client to the redis server
		rs, err := NewRedisServerFromConnectionString(sn, connectionString(master.IP, master.Port))
		defer rs.Cleanup(log.FromContext(ctx))
		if err != nil {
			logger.Error(err, fmt.Sprintf("unable to open client to master %s", sn))
			return nil, err
		}

		save, err := rs.CRUD.RedisConfigGet(ctx, "save")
		if err != nil {
			logger.Error(err, fmt.Sprintf("unable to get master %s 'save' option", sn))
			return nil, err
		}
		result[sn].Config["save"] = save
	}

	/////////////////////////////////
	// discover the shard's slaves //
	/////////////////////////////////

	if !opts.Has(OnlyMasterDiscoveryOpt) {
		slaves, err := ss.CRUD.SentinelSlaves(ctx, shard)
		if err != nil {
			logger.Error(err, fmt.Sprintf("unable to get slaves for shard %s", shard))
			return nil, err
		}

		for _, slave := range slaves {

			// do not try to discover slaves flagged as "s_down" or "o_down"
			if !strings.Contains(slave.Flags, "s_down") && !strings.Contains(slave.Flags, "o_down") {

				sn := serverName(slave.IP, slave.Port)
				result[sn] = saasv1alpha1.RedisServerDetails{
					Role:   redis_client.Slave,
					Config: map[string]string{},
				}

				if opts.Has(SaveConfigDiscoveryOpt) || opts.Has(SlaveReadOnlyDiscoveryOpt) {

					// open a client to the redis server
					rs, err := NewRedisServerFromConnectionString(sn, connectionString(slave.IP, slave.Port))
					defer rs.Cleanup(log.FromContext(ctx))
					if err != nil {
						logger.Error(err, fmt.Sprintf("unable to open client to slave %s", sn))
						return nil, err
					}

					if opts.Has(SaveConfigDiscoveryOpt) {
						save, err := rs.CRUD.RedisConfigGet(ctx, "save")
						if err != nil {
							logger.Error(err, fmt.Sprintf("unable to get slave %s 'save' option", sn))
							return nil, err
						}
						result[sn].Config["save"] = save
					}

					if opts.Has(SlaveReadOnlyDiscoveryOpt) {
						slaveReadOnly, err := rs.CRUD.RedisConfigGet(ctx, "slave-read-only")
						if err != nil {
							logger.Error(err, fmt.Sprintf("unable to get slave %s 'slave-read-only' option", sn))
							return nil, err
						}
						result[sn].Config["slave-read-only"] = slaveReadOnly
					}
				}

			}
		}
	}

	return result, nil
}

//

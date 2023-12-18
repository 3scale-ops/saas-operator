package sharded

import (
	"context"
	"sort"

	redis "github.com/3scale-ops/saas-operator/pkg/redis/server"
	operatorutils "github.com/3scale-ops/saas-operator/pkg/util"
)

const (
	shardNotInitializedError = "ERR No such master with that name"
)

// SentinelServer represents a sentinel Pod
type SentinelServer struct {
	*redis.Server
}

func NewSentinelServerFromPool(connectionString string, alias *string, pool *redis.ServerPool) (*SentinelServer, error) {
	srv, err := pool.GetServer(connectionString, alias)
	if err != nil {
		return nil, err
	}

	return &SentinelServer{
		Server: srv,
	}, nil
}

func NewSentinelServerFromParams(srv *redis.Server) *SentinelServer {
	return &SentinelServer{
		Server: srv,
	}
}

func NewHighAvailableSentinel(servers map[string]string, pool *redis.ServerPool) ([]*SentinelServer, error) {
	var merr operatorutils.MultiError
	sentinels := make([]*SentinelServer, 0, len(servers))

	for key, connectionString := range servers {
		var alias *string = nil
		if key != connectionString {
			alias = &key
		}
		srv, err := NewSentinelServerFromPool(connectionString, alias, pool)
		if err != nil {
			merr = append(merr, err)
			continue
		}
		sentinels = append(sentinels, srv)
	}

	sort.Slice(sentinels, func(i, j int) bool {
		return sentinels[i].ID() < sentinels[j].ID()
	})

	return sentinels, merr.ErrorOrNil()
}

// IsMonitoringShards checks whether or all the shards in the passed list are being monitored by the SentinelServer
func (sentinel *SentinelServer) IsMonitoringShards(ctx context.Context, shards []string) (bool, error) {

	monitoredShards, err := sentinel.SentinelMasters(ctx)
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
func (sentinel *SentinelServer) Monitor(ctx context.Context, cluster *Cluster, quorum int) ([]string, error) {
	changed := []string{}

	// Initialize unmonitored shards
	shardNames := cluster.GetShardNames()
	for _, name := range shardNames {

		_, err := sentinel.SentinelMaster(ctx, name)
		if err != nil {
			if err.Error() == shardNotInitializedError {

				shard := cluster.LookupShardByName(name)
				master, err := shard.GetMaster()
				if err != nil {
					return changed, err
				}

				err = sentinel.SentinelMonitor(ctx, name, master.GetHost(), master.GetPort(), quorum)
				if err != nil {
					return changed, operatorutils.WrapError("redis-sentinel/SentinelServer.Monitor", err)
				}
				// even if the next call fails, there has already been a write operation to sentinel
				changed = append(changed, name)

				err = sentinel.SentinelSet(ctx, name, "down-after-milliseconds", "5000")
				if err != nil {
					return changed, operatorutils.WrapError("redis-sentinel/SentinelServer.Monitor", err)
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

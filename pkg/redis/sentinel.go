package redis

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/redis/crud"
	"github.com/3scale/saas-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

// MonitoredShards returns the list of monitored shards of this SentinelServer
func (ss *SentinelServer) MonitoredShards(ctx context.Context, discoverSlaves bool) (saasv1alpha1.MonitoredShards, error) {

	sm, err := ss.CRUD.SentinelMasters(ctx)
	if err != nil {
		return nil, err
	}

	monitoredShards := make([]saasv1alpha1.MonitoredShard, 0, len(sm))
	for _, s := range sm {
		shard := saasv1alpha1.MonitoredShard{Name: s.Name, Master: fmt.Sprintf("%s:%d", s.IP, s.Port)}
		if discoverSlaves {
			// ignore errors trying to discover slaves so at least the master info
			// is returned
			shard.SlavesRO, shard.SlavesRW, _ = ss.DiscoverSlaves(ctx, s.Name, maxInfoCacheAge)
		}
		monitoredShards = append(monitoredShards, shard)
	}
	return monitoredShards, nil
}

func (ss *SentinelServer) DiscoverSlaves(ctx context.Context, shard string, maxInfoCacheAge time.Duration) ([]string, []string, error) {
	slavesRO := []string{}
	slavesRW := []string{}

	infoCache, err := ss.CRUD.SentinelInfoCache(ctx)
	if err != nil {
		return nil, nil, err
	}

	slaves, err := ss.CRUD.SentinelSlaves(ctx, shard)
	if err != nil {
		return nil, nil, err
	}

	for _, slave := range slaves {
		if !strings.Contains(slave.Flags, "s_down") {
			isRO, err := infoCache.GetValue(shard, slave.RunID, "slave_read_only", maxInfoCacheAge)
			if err != nil {
				// ignore this slave as we have no reliable information
				// to determine if its RO or RW
				continue
			}

			if isRO == "0" {
				slavesRW = append(slavesRW, fmt.Sprintf("%s:%d", slave.IP, slave.Port))
			} else {
				slavesRO = append(slavesRO, fmt.Sprintf("%s:%d", slave.IP, slave.Port))
			}
		}
	}

	// sort arrays before returning so sentinel discovery responses
	// can be easily compared
	sort.Strings(slavesRO)
	sort.Strings(slavesRW)

	return slavesRO, slavesRW, nil
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

// SentinelPool represents a pool of SentinelServers that monitor the same
// group of redis shards
type SentinelPool []SentinelServer

// NewSentinelPool creates a new SentinelPool object given a key and a number of replicas by calling the k8s API
// to discover sentinel Pods. The kye es the Name/Namespace of the StatefulSet that owns the sentinel Pods.
func NewSentinelPool(ctx context.Context, cl client.Client, key types.NamespacedName, replicas int) (SentinelPool, error) {

	spool := make([]SentinelServer, replicas)
	for i := 0; i < replicas; i++ {
		pod := &corev1.Pod{}
		key := types.NamespacedName{Name: fmt.Sprintf("%s-%d", key.Name, i), Namespace: key.Namespace}
		err := cl.Get(ctx, key, pod)
		if err != nil {
			return nil, err
		}

		ss, err := NewSentinelServerFromConnectionString(pod.GetName(), fmt.Sprintf("redis://%s:%d", pod.Status.PodIP, saasv1alpha1.SentinelPort))
		if err != nil {
			return nil, err
		}
		spool[i] = *ss
	}
	return spool, nil
}

// IsMonitoringShards checks whether or all the shards in the passed list are being monitored by all
// sentinel servers in the SentinelPool
func (sp SentinelPool) IsMonitoringShards(ctx context.Context, shards []string) (bool, error) {

	for _, ss := range sp {
		ok, err := ss.IsMonitoringShards(ctx, shards)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

// MonitoredShards returns the list of monitored shards of this SentinelServer
func (sp SentinelPool) MonitoredShards(ctx context.Context, quorum int, discoverSlaves bool) (saasv1alpha1.MonitoredShards, error) {

	responses := make([]saasv1alpha1.MonitoredShards, 0, len(sp))

	for _, srv := range sp {

		resp, err := srv.MonitoredShards(ctx, discoverSlaves)
		if err != nil {
			// jump to next sentinel if error occurs
			continue
		}
		responses = append(responses, resp)
	}

	monitoredShards, err := applyQuorum(responses, saasv1alpha1.SentinelDefaultQuorum)
	if err != nil {
		return nil, err
	}

	return monitoredShards, nil
}

// Monitor ensures that all the shards in the ShardedCluster object are monitored by
// all sentinel servers in the SentinelPool
func (sp SentinelPool) Monitor(ctx context.Context, shards ShardedCluster) (map[string][]string, error) {
	changes := map[string][]string{}
	for _, ss := range sp {
		ssChanges, err := ss.Monitor(ctx, shards)
		if err != nil {
			return changes, err
		}
		if len(ssChanges) > 0 {
			changes[ss.Name] = ssChanges
		}
	}
	return changes, nil
}

func applyQuorum(responses []saasv1alpha1.MonitoredShards, quorum int) (saasv1alpha1.MonitoredShards, error) {

	for _, r := range responses {
		// Sort each of the MonitoredShards responses to
		// avoid diffs due to unordered responses from redis
		sort.Sort(r)
	}

	for idx, a := range responses {
		count := 0
		for _, b := range responses {
			if reflect.DeepEqual(a, b) {
				count++
			}
		}

		// check if this response has quorum
		if count >= quorum {
			return responses[idx], nil
		}
	}

	return nil, fmt.Errorf("unable to get monitored shards from sentinel")
}

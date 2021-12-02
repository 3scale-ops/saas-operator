package redis

import (
	"context"
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/redis/crud"
	"github.com/3scale/saas-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	shardNotInitializedError = "ERR No such master with that name"
)

// SentinelServer represents a sentinel Pod
type SentinelServer struct {
	Name string
	CRUD *crud.CRUD
}

func NewSentinelServer(name, connectionString string) (*SentinelServer, error) {

	crud, err := crud.NewRedisCRUD(connectionString)
	if err != nil {
		return nil, err
	}

	return &SentinelServer{Name: name, CRUD: crud}, nil
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
func (ss *SentinelServer) Monitor(ctx context.Context, shards ShardedCluster) error {

	// Initialize unmonitored shards
	for name, shard := range shards {

		_, err := ss.CRUD.SentinelMaster(ctx, name)
		if err != nil {
			if err.Error() == shardNotInitializedError {
				host, port, err := shard.GetMasterAddr()
				if err != nil {
					return err
				}
				err = ss.CRUD.SentinelMonitor(ctx, name, host, port, 2)
				if err != nil {
					return util.WrapError("[redis-sentinel/SentinelServer.Monitor]", err)
				}
				err = ss.CRUD.SentinelSet(ctx, name, "down-after-milliseconds", "5000")
				if err != nil {
					return util.WrapError("[redis-sentinel/SentinelServer.Monitor]", err)
				}

				// TODO: change the default failover timeout.
				// TODO: maybe add a generic mechanism to set/modify parameters

			} else {
				return err
			}
		}
	}

	return nil
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

		ss, err := NewSentinelServer(pod.GetName(), fmt.Sprintf("redis://%s:%d", pod.Status.PodIP, saasv1alpha1.SentinelPort))
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

	for _, connString := range sp {
		ss := SentinelServer(connString)
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

// Monitor ensures that all the shards in the ShardedCluster object are monitored by
// all sentinel servers in the SentinelPool
func (sp SentinelPool) Monitor(ctx context.Context, shards ShardedCluster) error {

	for _, connString := range sp {
		ss := SentinelServer(connString)
		err := ss.Monitor(ctx, shards)
		if err != nil {
			return err
		}
	}
	return nil
}

package redis

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

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

// Cleanup closes all Redis clients opened during the SentinelPool object creation
func (sp SentinelPool) Cleanup(log logr.Logger) []error {
	log.V(1).Info("[@sentinel-pool-cleanup] closing clients")
	var closeErrors []error
	for _, ss := range sp {
		if err := ss.Cleanup(log); err != nil {
			closeErrors = append(closeErrors, err)
		}
	}
	return closeErrors
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

// MonitoredShards returns the list of monitored shards of this SentinelServer
func (sp SentinelPool) MonitoredShards(ctx context.Context, quorum int, options ...ShardDiscoveryOption) (saasv1alpha1.MonitoredShards, error) {
	logger := log.FromContext(ctx, "function", "(SentinelPool).MonitoredShards")
	responses := make([]saasv1alpha1.MonitoredShards, 0, len(sp))

	for _, srv := range sp {

		resp, err := srv.MonitoredShards(ctx, options...)
		if err != nil {
			logger.Error(err, "error getting monitored shards from sentinel", "SentinelServer", srv.Name)
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

	return nil, fmt.Errorf("no quorum of %d sentinels when getting monitored shards", saasv1alpha1.SentinelDefaultQuorum)
}

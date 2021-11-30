package redis

import (
	"context"
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	redistypes "github.com/3scale/saas-operator/pkg/redis/types"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-redis/redis/v8"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	shardNotInitializedError = "ERR No such master with that name"
)

type SentinelServer string

func (ss *SentinelServer) IsMonitoringShards(ctx context.Context, shards []string) (bool, error) {

	monitoredShards, err := ss.Masters(ctx)
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

func (ss *SentinelServer) Monitor(ctx context.Context, shards ShardedCluster) error {

	opt, err := redis.ParseURL(string(*ss))
	if err != nil {
		return err
	}

	sentinel := redis.NewSentinelClient(opt)

	// Initialize unmonitored shards
	for name, shard := range shards {

		result := &redistypes.SentinelMasterCmdResult{}
		err = sentinel.Master(ctx, name).Scan(result)
		if err != nil {
			if err.Error() == shardNotInitializedError {
				host, port, err := shard.GetMasterAddr()
				if err != nil {
					return err
				}
				_, err = sentinel.Monitor(ctx, name, host, port, "2").Result()
				if err != nil {
					return util.WrapError("[redis-sentinel/SentinelServer.Monitor]", err)
				}
				_, err = sentinel.Set(ctx, name, "down-after-milliseconds", "5000").Result()
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

func (ss *SentinelServer) Masters(ctx context.Context) ([]redistypes.SentinelMasterCmdResult, error) {
	opt, err := redis.ParseURL(string(*ss))
	if err != nil {
		return nil, err
	}

	sentinel := redis.NewSentinelClient(opt)

	values, err := sentinel.Masters(ctx).Result()
	if err != nil {
		return nil, err
	}

	result := make([]redistypes.SentinelMasterCmdResult, len(values))
	for i, val := range values {
		// shard := val.([]interface{})[1].(string)
		// master := &redistypes.SentinelMasterCmdResult{}
		// err := sentinel.Master(ctx, shard).Scan(master)
		masterResult := &redistypes.SentinelMasterCmdResult{}
		sliceCmdToStruct(val, masterResult)
		if err != nil {
			return nil, err
		}
		result[i] = *masterResult
	}

	return result, nil
}

func (ss *SentinelServer) Slaves(ctx context.Context, shard string) ([]redistypes.SentinelSlaveCmdResult, error) {
	opt, err := redis.ParseURL(string(*ss))
	if err != nil {
		return nil, err
	}

	sentinel := redis.NewSentinelClient(opt)

	values, err := sentinel.Slaves(ctx, shard).Result()
	if err != nil {
		return nil, err
	}

	result := make([]redistypes.SentinelSlaveCmdResult, len(values))
	for i, val := range values {
		slaveResult := &redistypes.SentinelSlaveCmdResult{}
		sliceCmdToStruct(val, slaveResult)
		result[i] = *slaveResult
	}

	return result, nil
}

func (ss *SentinelServer) Subscribe(ctx context.Context, events ...string) (<-chan *redis.Message, func() error, error) {
	opt, err := redis.ParseURL(string(*ss))
	if err != nil {
		return nil, nil, err
	}

	sentinel := redis.NewSentinelClient(opt)
	pubsub := sentinel.PSubscribe(ctx, events...)
	return pubsub.Channel(), pubsub.Close, nil
}

type SentinelPool []string

func NewSentinelPool(ctx context.Context, cl client.Client, key types.NamespacedName, replicas int) (SentinelPool, error) {

	spool := SentinelPool{}
	for i := 0; i < replicas; i++ {
		pod := &corev1.Pod{}
		key := types.NamespacedName{Name: fmt.Sprintf("%s-%d", key.Name, i), Namespace: key.Namespace}
		err := cl.Get(ctx, key, pod)
		if err != nil {
			return nil, err
		}
		spool = append(spool, fmt.Sprintf("redis://%s:%d", pod.Status.PodIP, saasv1alpha1.SentinelPort))
	}
	return spool, nil
}

func (sp SentinelPool) IsMonitoringShards(ctx context.Context, shards []string) (bool, error) {

	for _, connString := range sp {
		sentinel := SentinelServer(connString)
		ok, err := sentinel.IsMonitoringShards(ctx, shards)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

func (sp SentinelPool) Monitor(ctx context.Context, shards ShardedCluster) error {

	for _, connString := range sp {
		sentinel := SentinelServer(connString)
		err := sentinel.Monitor(ctx, shards)
		if err != nil {
			return err
		}
	}
	return nil
}

// This is a horrible function to parse the horrible structs that the redis-go
// client returns for administrative commands. I swear it's not my fault ...
func sliceCmdToStruct(in interface{}, out interface{}) (interface{}, error) {
	m := map[string]string{}
	for i := range in.([]interface{}) {
		if i%2 != 0 {
			continue
		}
		m[in.([]interface{})[i].(string)] = in.([]interface{})[i+1].(string)
	}

	err := redis.NewStringStringMapResult(m, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

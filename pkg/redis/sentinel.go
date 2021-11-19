package redis

import (
	"context"
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-redis/redis/v8"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	shardNotInitializedError = "ERR No such master with that name"
)

type SentinelMasterCmdResult struct {
	Name         string `redis:"name"`
	IP           string `redis:"ip"`
	Port         string `redis:"port"`
	RunID        string `redis:"runid"`
	Flags        string `redis:"flags"`
	RoleReported string `redis:"role-reported"`
	NumSlaves    string `redis:"num-slaves"`
}

type SentinelServer string

func (ss *SentinelServer) Monitor(ctx context.Context, shards ShardedCluster) error {

	opt, err := redis.ParseURL(string(*ss))
	if err != nil {
		return err
	}

	sentinel := redis.NewSentinelClient(opt)

	// Check that all shards are being monitored, initialize if not
	for name, shard := range shards {

		result := &SentinelMasterCmdResult{}
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

				// TODO: change the default failover timeoute.
				// TODO: maybe add a generic mechanism to set/modify parameters

			} else {
				return err
			}
		}
	}

	return nil
}

type SentinelPool []string

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

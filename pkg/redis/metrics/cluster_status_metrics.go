package metrics

import (
	"context"

	"github.com/3scale/saas-operator/pkg/redis/sharded"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	serverInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "server_info",
			Namespace: "saas_redis_cluster_status",
			Help:      `"redis server info"`,
		},
		[]string{"resource", "shard", "redis_server_host", "redis_server_alias", "role", "read_only"})
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(serverInfo)
}

func FromShardedCluster(ctx context.Context, cluster *sharded.Cluster, refresh bool, resource string) error {

	if refresh {
		err := cluster.SentinelDiscover(ctx, sharded.SlaveReadOnlyDiscoveryOpt)
		if err != nil {
			return err
		}
	}

	for _, shard := range cluster.Shards {

		for _, server := range shard.Servers {
			ro, ok := server.Config["slave-read-only"]
			if !ok {
				ro = "no"
			}
			serverInfo.With(prometheus.Labels{"resource": resource, "shard": shard.Name,
				"redis_server_host": server.ID(), "redis_server_alias": server.GetAlias(),
				"role": string(server.Role), "read_only": ro,
			}).Set(float64(1))
		}
	}

	return nil
}

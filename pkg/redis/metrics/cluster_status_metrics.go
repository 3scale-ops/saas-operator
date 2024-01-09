package metrics

import (
	"context"

	"github.com/3scale-ops/saas-operator/pkg/redis/client"
	"github.com/3scale-ops/saas-operator/pkg/redis/sharded"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	serverInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "server_info",
			Namespace: "saas_redis_cluster_status",
			Help:      "redis cluster member info",
		},
		[]string{"resource", "shard", "redis_server_host", "redis_server_alias", "role", "read_only"})
	roSlaveCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "ro_slave_count",
			Namespace: "saas_redis_cluster_status",
			Help:      "read-only slave count",
		},
		[]string{"resource", "shard"},
	)
	rwSlaveCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "rw_slave_count",
			Namespace: "saas_redis_cluster_status",
			Help:      "read-write slave count",
		},
		[]string{"resource", "shard"},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(serverInfo, roSlaveCount, rwSlaveCount)
}

func FromShardedCluster(ctx context.Context, cluster *sharded.Cluster, refresh bool, resource string) error {

	if refresh {
		err := cluster.SentinelDiscover(ctx, sharded.SlaveReadOnlyDiscoveryOpt)
		if err != nil {
			return err
		}
	}

	for _, shard := range cluster.Shards {
		roslave := 0
		rwslave := 0

		for _, server := range shard.Servers {
			ro, ok := server.Config["slave-read-only"]
			if !ok {
				ro = "no"
			}
			serverInfo.With(prometheus.Labels{"resource": resource, "shard": shard.Name,
				"redis_server_host": server.ID(), "redis_server_alias": server.GetAlias(),
				"role": string(server.Role), "read_only": ro,
			}).Set(float64(1))

			if server.Role == client.Slave {
				if ro == "yes" {
					roslave++
				} else {
					rwslave++
				}
			}
		}

		roSlaveCount.With(prometheus.Labels{"resource": resource, "shard": shard.Name}).Set(float64(roslave))
		rwSlaveCount.With(prometheus.Labels{"resource": resource, "shard": shard.Name}).Set(float64(rwslave))
	}

	return nil
}

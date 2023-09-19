package backup

import (
	"math"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// metrics
var (
	backupSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "size",
			Namespace: "saas_redis_backup",
			Help:      `"size of the latest backup in bytes"`,
		},
		[]string{"shard"})
	backupFailureCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "failure_count",
			Namespace: "saas_redis_backup",
			Help:      `"total number of backup failures"`,
		},
		[]string{"shard"})
	backupSuccessCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "success_count",
			Namespace: "saas_redis_backup",
			Help:      `"total number of backup successes"`,
		},
		[]string{"shard"})
	backupDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "duration",
			Namespace: "saas_redis_backup",
			Help:      `"seconds it took to complete the backup"`,
		},
		[]string{"shard"})
)

func init() {
	// Register backup metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		backupSize, backupFailureCount, backupDuration, backupSuccessCount,
	)
}

func (r *Runner) publishMetrics() {
	// ensure counters are initialized
	if err := backupSuccessCount.With(prometheus.Labels{"shard": r.ShardName}).Write(&dto.Metric{}); err != nil {
		backupFailureCount.With(prometheus.Labels{"shard": r.ShardName}).Add(0)
	}
	if err := backupFailureCount.With(prometheus.Labels{"shard": r.ShardName}).Write(&dto.Metric{}); err != nil {
		backupFailureCount.With(prometheus.Labels{"shard": r.ShardName}).Add(0)
	}
	// update metrics
	if r.status.Error != nil {
		backupSize.With(prometheus.Labels{"shard": r.ShardName}).Set(float64(0))
		backupFailureCount.With(prometheus.Labels{"shard": r.ShardName}).Inc()
	} else {
		backupSize.With(prometheus.Labels{"shard": r.ShardName}).Set(float64(r.status.BackupSize))
		backupDuration.With(prometheus.Labels{"shard": r.ShardName}).Set(math.Round(r.status.FinishedAt.Sub(r.Timestamp).Seconds()))
		backupSuccessCount.With(prometheus.Labels{"shard": r.ShardName}).Add(1)
	}
}

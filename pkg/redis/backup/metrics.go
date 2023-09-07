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
	backupFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "failures",
			Namespace: "saas_redis_backup",
			Help:      `"total number of backup failures"`,
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
		backupSize, backupFailures, backupDuration,
	)
}

func (r *Runner) publishMetrics() {
	if r.status.Error != nil {
		backupSize.With(prometheus.Labels{"shard": r.ShardName}).Set(float64(0))
		backupFailures.With(prometheus.Labels{"shard": r.ShardName}).Inc()
	} else {
		backupSize.With(prometheus.Labels{"shard": r.ShardName}).Set(float64(r.status.BackupSize))
		backupDuration.With(prometheus.Labels{"shard": r.ShardName}).Set(math.Round(r.status.FinishedAt.Sub(r.Timestamp).Seconds()))
		// ensure the failure counter is initialized
		if err := backupFailures.With(prometheus.Labels{"shard": r.ShardName}).Write(&dto.Metric{}); err != nil {
			backupFailures.With(prometheus.Labels{"shard": r.ShardName}).Add(0)
		}
	}
}

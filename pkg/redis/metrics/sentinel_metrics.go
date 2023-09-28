package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/3scale/saas-operator/pkg/reconcilers/threads"
	"github.com/3scale/saas-operator/pkg/redis/client"
	redis "github.com/3scale/saas-operator/pkg/redis/server"
	"github.com/3scale/saas-operator/pkg/redis/sharded"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	serverInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "server_info",
			Namespace: "saas_redis_sentinel",
			Help:      `"redis server info"`,
		},
		[]string{"sentinel", "shard", "redis_server", "role"})
	linkPendingCommands = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "link_pending_commands",
			Namespace: "saas_redis_sentinel",
			Help:      `"sentinel master <name> link-pending-commands"`,
		},
		[]string{"sentinel", "shard", "redis_server", "role"},
	)
	lastOkPingReply = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "last_ok_ping_reply",
			Namespace: "saas_redis_sentinel",
			Help:      `"sentinel master <name> last-ok-ping-reply"`,
		},
		[]string{"sentinel", "shard", "redis_server", "role"},
	)
	roleReportedTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "role_reported_time",
			Namespace: "saas_redis_sentinel",
			Help:      `"sentinel master <name> role-reported-time"`,
		},
		[]string{"sentinel", "shard", "redis_server", "role"},
	)
	numSlaves = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "num_slaves",
			Namespace: "saas_redis_sentinel",
			Help:      `"sentinel master <name> num-slaves"`,
		},
		[]string{"sentinel", "shard", "redis_server", "role"},
	)
	numOtherSentinels = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "num_other_sentinels",
			Namespace: "saas_redis_sentinel",
			Help:      `"sentinel master <name> num-other-sentinels"`,
		},
		[]string{"sentinel", "shard", "redis_server", "role"},
	)

	masterLinkDownTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "master_link_down_time",
			Namespace: "saas_redis_sentinel",
			Help:      `"sentinel slaves master-link-down-time"`,
		},
		[]string{"sentinel", "shard", "redis_server", "role"},
	)

	slaveReplOffset = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "slave_repl_offset",
			Namespace: "saas_redis_sentinel",
			Help:      `"sentinel slaves slave-repl-offset"`,
		},
		[]string{"sentinel", "shard", "redis_server", "role"},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		serverInfo, linkPendingCommands, lastOkPingReply, roleReportedTime,
		numSlaves, numOtherSentinels, masterLinkDownTime, slaveReplOffset,
	)
}

// SentinelEventWatcher implements RunnableThread
var _ threads.RunnableThread = &SentinelMetricsGatherer{}

// SentinelMetricsGatherer is used to export sentinel metrics, obtained
// thrugh several admin commands, as prometheus metrics
type SentinelMetricsGatherer struct {
	refreshInterval time.Duration
	sentinelURI     string
	sentinel        *sharded.SentinelServer
	started         bool
	cancel          context.CancelFunc
}

func NewSentinelMetricsGatherer(sentinelURI string, refreshInterval time.Duration, pool *redis.ServerPool) (*SentinelMetricsGatherer, error) {
	sentinel, err := sharded.NewSentinelServerFromPool(sentinelURI, nil, pool)
	if err != nil {
		return nil, err
	}

	return &SentinelMetricsGatherer{
		refreshInterval: refreshInterval,
		sentinelURI:     sentinelURI,
		sentinel:        sentinel,
	}, nil
}

func (fw *SentinelMetricsGatherer) GetID() string {
	return fw.sentinelURI
}

// IsStarted returns whether the metrics gatherer is running or not
func (smg *SentinelMetricsGatherer) IsStarted() bool {
	return smg.started
}

func (smg *SentinelMetricsGatherer) CanBeDeleted() bool {
	return true
}

// SetChannel is required for SentinelMetricsGatherer to implement the RunnableThread
// interface, but it actually does nothing with the channel.
func (fw *SentinelMetricsGatherer) SetChannel(ch chan event.GenericEvent) {}

// Start starts metrics gatherer for sentinel
func (smg *SentinelMetricsGatherer) Start(parentCtx context.Context, l logr.Logger) error {
	log := l.WithValues("sentinel", smg.sentinelURI)
	if smg.started {
		log.Info("the metrics gatherer is already running")
		return nil
	}

	go func() {
		var ctx context.Context
		ctx, smg.cancel = context.WithCancel(parentCtx)

		ticker := time.NewTicker(smg.refreshInterval)

		log.Info("sentinel metrics gatherer running")

		for {
			select {

			case <-ticker.C:
				if err := smg.gatherMetrics(ctx); err != nil {
					log.Error(err, "error gathering sentinel metrics")
				}

			case <-ctx.Done():
				log.Info("shutting down sentinel metrics gatherer")
				smg.started = false
				return
			}
		}
	}()

	smg.started = true
	return nil
}

// Stop stops metrics gatherering for sentinel
func (smg *SentinelMetricsGatherer) Stop() {
	// stop gathering metrics
	smg.cancel()
	// Reset all gauge metrics so the values related to
	// this exporter are deleted from the collection
	serverInfo.Reset()
	linkPendingCommands.Reset()
	lastOkPingReply.Reset()
	roleReportedTime.Reset()
	numSlaves.Reset()
	numOtherSentinels.Reset()
	masterLinkDownTime.Reset()
	slaveReplOffset.Reset()
}

func (smg *SentinelMetricsGatherer) gatherMetrics(ctx context.Context) error {

	mresult, err := smg.sentinel.SentinelMasters(ctx)
	if err != nil {
		return err
	}

	for _, master := range mresult {

		serverInfo.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port), "role": master.RoleReported,
		}).Set(float64(1))

		linkPendingCommands.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port), "role": master.RoleReported,
		}).Set(float64(master.LinkPendingCommands))

		lastOkPingReply.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port), "role": master.RoleReported,
		}).Set(float64(master.LastOkPingReply))

		roleReportedTime.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port), "role": master.RoleReported,
		}).Set(float64(master.RoleReportedTime))

		numSlaves.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port), "role": master.RoleReported,
		}).Set(float64(master.NumSlaves))

		numOtherSentinels.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port), "role": master.RoleReported,
		}).Set(float64(master.NumOtherSentinels))

		sresult, err := smg.sentinel.SentinelSlaves(ctx, master.Name)
		if err != nil {
			return err
		}

		// Cleanup any vector that corresponds to the same server but with a
		// different role to avoid stale metrics after a role switch
		cleanupMetrics(prometheus.Labels{
			"sentinel":     smg.sentinelURI,
			"shard":        master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port),
			"role":         string(client.Slave),
		})

		for _, slave := range sresult {

			serverInfo.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port), "role": slave.RoleReported,
			}).Set(float64(1))

			linkPendingCommands.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port), "role": slave.RoleReported,
			}).Set(float64(slave.LinkPendingCommands))

			lastOkPingReply.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port), "role": slave.RoleReported,
			}).Set(float64(slave.LastOkPingReply))

			roleReportedTime.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port), "role": slave.RoleReported,
			}).Set(float64(slave.RoleReportedTime))

			masterLinkDownTime.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port), "role": slave.RoleReported,
			}).Set(float64(slave.MasterLinkDownTime))

			slaveReplOffset.With(prometheus.Labels{"sentinel": smg.sentinelURI, "shard": master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port), "role": slave.RoleReported,
			}).Set(float64(slave.SlaveReplOffset))

			cleanupMetrics(prometheus.Labels{
				"sentinel":     smg.sentinelURI,
				"shard":        master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port),
				"role":         string(client.Master),
			})
		}
	}

	return nil
}

func cleanupMetrics(labels prometheus.Labels) {
	serverInfo.Delete(labels)
	linkPendingCommands.Delete(labels)
	lastOkPingReply.Delete(labels)
	roleReportedTime.Delete(labels)
	numSlaves.Delete(labels)
	numOtherSentinels.Delete(labels)
	masterLinkDownTime.Delete(labels)
	slaveReplOffset.Delete(labels)
}

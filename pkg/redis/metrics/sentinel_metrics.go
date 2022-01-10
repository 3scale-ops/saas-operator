package metrics

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/3scale/saas-operator/pkg/redis"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
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
		linkPendingCommands, lastOkPingReply, roleReportedTime,
		numSlaves, numOtherSentinels, masterLinkDownTime, slaveReplOffset,
	)
}

// SentinelMetricsGatherer is used to export sentinel metrics, obtained
// thrugh several admin commands, as prometheus metrics
type SentinelMetricsGatherer struct {
	RefreshInterval time.Duration
	SentinelURI     string
	Log             logr.Logger
	started         bool
	cancel          context.CancelFunc
	sentinel        *redis.SentinelServer
}

// IsStarted returns whether the metrics gatherer is running or not
func (smg *SentinelMetricsGatherer) IsStarted() bool {
	return smg.started
}

//Start starts metrics gatherer for sentinel
func (smg *SentinelMetricsGatherer) Start(parentCtx context.Context) {
	log := smg.Log.WithValues("sentinel", smg.SentinelURI)
	if smg.started {
		log.Error(errors.New("already started"), "the metrics gatherer is already running")
		return
	}

	go func() {
		var err error
		var ctx context.Context
		ctx, smg.cancel = context.WithCancel(parentCtx)

		ticker := time.NewTicker(smg.RefreshInterval)

		smg.sentinel, err = redis.NewSentinelServerFromConnectionString(smg.SentinelURI, smg.SentinelURI)
		if err != nil {
			log.Error(err, "cannot create SentinelServer")
		}

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
}

// Stop stops metrics gatherering for sentinel
func (smg *SentinelMetricsGatherer) Stop() {
	// stop gathering metrics
	smg.cancel()
	// Reset all gauge metrics so the values related to
	// this exporter are deleted from the collection
	linkPendingCommands.Reset()
	lastOkPingReply.Reset()
	roleReportedTime.Reset()
	numSlaves.Reset()
	numOtherSentinels.Reset()
	masterLinkDownTime.Reset()
	slaveReplOffset.Reset()
}

func (smg *SentinelMetricsGatherer) gatherMetrics(ctx context.Context) error {

	mresult, err := smg.sentinel.CRUD.SentinelMasters(ctx)
	if err != nil {
		return err
	}

	for _, master := range mresult {
		linkPendingCommands.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port), "role": master.RoleReported,
		}).Set(float64(master.LinkPendingCommands))

		lastOkPingReply.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port), "role": master.RoleReported,
		}).Set(float64(master.LastOkPingReply))

		roleReportedTime.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port), "role": master.RoleReported,
		}).Set(float64(master.RoleReportedTime))

		numSlaves.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port), "role": master.RoleReported,
		}).Set(float64(master.NumSlaves))

		numOtherSentinels.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": master.Name,
			"redis_server": fmt.Sprintf("%s:%d", master.IP, master.Port), "role": master.RoleReported,
		}).Set(float64(master.NumOtherSentinels))

		sresult, err := smg.sentinel.CRUD.SentinelSlaves(ctx, master.Name)
		if err != nil {
			return err
		}

		for _, slave := range sresult {

			linkPendingCommands.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port), "role": slave.RoleReported,
			}).Set(float64(slave.LinkPendingCommands))

			lastOkPingReply.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port), "role": slave.RoleReported,
			}).Set(float64(slave.LastOkPingReply))

			roleReportedTime.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port), "role": slave.RoleReported,
			}).Set(float64(slave.RoleReportedTime))

			masterLinkDownTime.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port), "role": slave.RoleReported,
			}).Set(float64(slave.MasterLinkDownTime))

			slaveReplOffset.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": master.Name,
				"redis_server": fmt.Sprintf("%s:%d", slave.IP, slave.Port), "role": slave.RoleReported,
			}).Set(float64(slave.SlaveReplOffset))
		}
	}

	return nil
}

package events

import (
	"context"
	"fmt"

	"github.com/3scale/saas-operator/pkg/reconcilers/threads"
	redis "github.com/3scale/saas-operator/pkg/redis/server"
	"github.com/3scale/saas-operator/pkg/redis/sharded"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	switchMasterCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "switch_master_count",
			Namespace: "saas_redis_sentinel",
			Help:      "+switch-master (https://redis.io/topics/sentinel#sentinel-api)",
		},
		[]string{"sentinel", "shard"},
	)

	failoverAbortNoGoodSlaveCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "failover_abort_no_good_slave_count",
			Namespace: "saas_redis_sentinel",
			Help:      "-failover-abort-no-good-slave (https://redis.io/topics/sentinel#sentinel-api)",
		},
		[]string{"sentinel", "shard"},
	)
	sdownCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "sdown_count",
			Namespace: "saas_redis_sentinel",
			Help:      "+sdown (https://redis.io/topics/sentinel#sentinel-api)",
		},
		[]string{"sentinel", "shard", "redis_server"},
	)
	sdownSentinelCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "sdown_sentinel_count",
			Namespace: "saas_redis_sentinel",
			Help:      "+sdown (https://redis.io/topics/sentinel#sentinel-api)",
		},
		[]string{"sentinel", "shard", "redis_server"},
	)
	sdownClearedCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "sdown_cleared_count",
			Namespace: "saas_redis_sentinel",
			Help:      "-sdown (https://redis.io/topics/sentinel#sentinel-api)",
		},
		[]string{"sentinel", "shard", "redis_server"},
	)
	sdownClearedSentinelCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "sdown_cleared_sentinel_count",
			Namespace: "saas_redis_sentinel",
			Help:      "-sdown (https://redis.io/topics/sentinel#sentinel-api)",
		},
		[]string{"sentinel", "shard", "redis_server"},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(switchMasterCount, failoverAbortNoGoodSlaveCount,
		sdownCount, sdownSentinelCount, sdownClearedCount, sdownClearedSentinelCount)
}

// SentinelEventWatcher implements RunnableThread
var _ threads.RunnableThread = &SentinelEventWatcher{}

type SentinelEventWatcher struct {
	instance      client.Object
	sentinelURI   string
	exportMetrics bool
	topology      *sharded.Cluster
	eventsCh      chan event.GenericEvent
	started       bool
	cancel        context.CancelFunc
	sentinel      *sharded.SentinelServer
}

func NewSentinelEventWatcher(sentinelURI string, instance client.Object, topology *sharded.Cluster,
	metrics bool, pool *redis.ServerPool) (*SentinelEventWatcher, error) {
	sentinel, err := sharded.NewSentinelServerFromPool(sentinelURI, nil, pool)
	if err != nil {
		return nil, err
	}

	return &SentinelEventWatcher{
		instance:      instance,
		sentinelURI:   sentinelURI,
		exportMetrics: metrics,
		topology:      topology,
		sentinel:      sentinel,
	}, nil
}

func (sew *SentinelEventWatcher) GetID() string {
	return sew.sentinelURI
}

// IsStarted returns whether the metrics gatherer is running or not
func (sew *SentinelEventWatcher) IsStarted() bool {
	return sew.started
}

func (sew *SentinelEventWatcher) CanBeDeleted() bool {
	return true
}

func (sew *SentinelEventWatcher) SetChannel(ch chan event.GenericEvent) {
	sew.eventsCh = ch
}

func (sew *SentinelEventWatcher) Cleanup() error {
	return sew.sentinel.CloseClient()
}

// Start starts metrics gatherer for sentinel
func (sew *SentinelEventWatcher) Start(parentCtx context.Context, l logr.Logger) error {
	log := l.WithValues("sentinel", sew.sentinelURI)
	if sew.started {
		log.Info("the event watcher is already running")
		return nil
	}

	if sew.exportMetrics {
		// Initializes metrics with 0 value
		sew.initCounters()
	}

	go func() {
		var ctx context.Context
		ctx, sew.cancel = context.WithCancel(parentCtx)

		ch, closeWatch := sew.sentinel.SentinelPSubscribe(ctx,
			`+switch-master`,
			`-failover-abort-no-good-slave`,
			`[+\-]sdown`,
		)
		defer closeWatch()

		log.Info("event watcher running")

		for {
			select {

			case msg := <-ch:
				log.V(1).Info("received event from sentinel", "event", msg.String())
				sew.eventsCh <- event.GenericEvent{Object: sew.instance}
				rem, err := NewRedisEventMessage(msg)
				if err == nil {
					log.V(3).Info("redis event message parsed",
						"event", rem.event,
						"target-type", rem.target.role, "target-name", rem.target.name,
						"target-ip", rem.target.ip, "target-port", rem.target.port,
						"master-type", rem.master.role, "master-name", rem.master.name,
						"master-ip", rem.master.ip, "master-port", rem.target.port,
					)
					if sew.exportMetrics {
						sew.metricsFromEvent(rem)
					}
				} else {
					log.Error(err, "invalid event message")
				}

			case <-ctx.Done():
				log.Info("shutting down event watcher")
				sew.started = false
				return
			}
		}
	}()

	sew.started = true
	return nil
}

// Stop stops the sentinel event watcher
func (sew *SentinelEventWatcher) Stop() {
	sew.cancel()
}

func (sew *SentinelEventWatcher) metricsFromEvent(rem RedisEventMessage) {
	switch rem.event {
	case "+switch-master":
		switchMasterCount.With(
			prometheus.Labels{
				"sentinel": sew.sentinelURI, "shard": rem.master.name,
			},
		).Add(1)
	case "-failover-abort-no-good-slave":
		failoverAbortNoGoodSlaveCount.With(
			prometheus.Labels{
				"sentinel": sew.sentinelURI, "shard": rem.target.name,
			},
		).Add(1)
	case "+sdown":
		switch rem.target.role {
		case "sentinel":
			sdownSentinelCount.With(
				prometheus.Labels{
					"sentinel": sew.sentinelURI, "shard": rem.master.name,
					"redis_server": fmt.Sprintf("%s:%s", rem.target.ip, rem.target.port),
				},
			).Add(1)
		default:
			sdownCount.With(
				prometheus.Labels{
					"sentinel": sew.sentinelURI, "shard": rem.master.name,
					"redis_server": fmt.Sprintf("%s:%s", rem.target.ip, rem.target.port),
				},
			).Add(1)
		}
	case "-sdown":
		switch rem.target.role {
		case "sentinel":
			sdownClearedSentinelCount.With(
				prometheus.Labels{
					"sentinel": sew.sentinelURI, "shard": rem.master.name,
					"redis_server": fmt.Sprintf("%s:%s", rem.target.ip, rem.target.port),
				},
			).Add(1)
		default:
			sdownClearedCount.With(
				prometheus.Labels{
					"sentinel": sew.sentinelURI, "shard": rem.master.name,
					"redis_server": fmt.Sprintf("%s:%s", rem.target.ip, rem.target.port),
				},
			).Add(1)
		}
	}
}

func (sew *SentinelEventWatcher) initCounters() {
	if sew.topology != nil {

		for _, shard := range sew.topology.Shards {
			failoverAbortNoGoodSlaveCount.With(
				prometheus.Labels{
					"sentinel": sew.sentinelURI, "shard": shard.Name,
				},
			).Add(0)

			for _, server := range shard.Servers {
				switchMasterCount.With(
					prometheus.Labels{
						"sentinel": sew.sentinelURI, "shard": shard.Name,
					},
				).Add(0)
				sdownSentinelCount.With(
					prometheus.Labels{
						"sentinel": sew.sentinelURI, "shard": shard.Name,
						"redis_server": server.ID(),
					},
				).Add(0)
				sdownCount.With(
					prometheus.Labels{
						"sentinel": sew.sentinelURI, "shard": shard.Name,
						"redis_server": server.ID(),
					},
				).Add(0)
				sdownClearedSentinelCount.With(
					prometheus.Labels{
						"sentinel": sew.sentinelURI, "shard": shard.Name,
						"redis_server": server.ID(),
					},
				).Add(0)
				sdownClearedCount.With(
					prometheus.Labels{
						"sentinel": sew.sentinelURI, "shard": shard.Name,
						"redis_server": server.ID(),
					},
				).Add(0)
			}
		}

	}
}

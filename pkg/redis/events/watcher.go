package events

import (
	"context"
	"fmt"

	"github.com/3scale/saas-operator/pkg/reconcilers/threads"
	"github.com/3scale/saas-operator/pkg/redis"
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
		[]string{"sentinel", "shard", "redis_server"},
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
	Instance      client.Object
	SentinelURI   string
	ExportMetrics bool
	Topology      *redis.ShardedCluster
	eventsCh      chan event.GenericEvent
	started       bool
	cancel        context.CancelFunc
	sentinel      *redis.SentinelServer
}

func (sew *SentinelEventWatcher) GetID() string {
	return sew.SentinelURI
}

// IsStarted returns whether the metrics gatherer is running or not
func (sew *SentinelEventWatcher) IsStarted() bool {
	return sew.started
}

func (sew *SentinelEventWatcher) SetChannel(ch chan event.GenericEvent) {
	sew.eventsCh = ch
}

func (sew *SentinelEventWatcher) Cleanup() error {
	return sew.sentinel.CRUD.CloseClient()
}

// Start starts metrics gatherer for sentinel
func (sew *SentinelEventWatcher) Start(parentCtx context.Context, l logr.Logger) error {
	log := l.WithValues("sentinel", sew.SentinelURI)
	if sew.started {
		log.Info("the event watcher is already running")
		return nil
	}

	if sew.ExportMetrics {
		// Initializes metrics with 0 value
		sew.initCounters()
	}

	var err error
	sew.sentinel, err = redis.NewSentinelServerFromConnectionString(sew.SentinelURI, sew.SentinelURI)
	if err != nil {
		log.Error(err, "cannot create SentinelServer")
		return err
	}

	go func() {
		var ctx context.Context
		ctx, sew.cancel = context.WithCancel(parentCtx)

		ch, closeWatch := sew.sentinel.CRUD.SentinelPSubscribe(ctx,
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
				sew.eventsCh <- event.GenericEvent{Object: sew.Instance}
				rem, err := NewRedisEventMessage(msg)
				if err == nil {
					log.V(3).Info("redis event message parsed",
						"event", rem.event,
						"target-type", rem.target.role, "target-name", rem.target.name,
						"target-ip", rem.target.ip, "target-port", rem.target.port,
						"master-type", rem.master.role, "master-name", rem.master.name,
						"master-ip", rem.master.ip, "master-port", rem.target.port,
					)
					if sew.ExportMetrics {
						sew.metricsFromEvent(rem)
					}
				} else {
					log.Error(err, "invalid event message")
				}

			case <-ctx.Done():
				log.Info("shutting down event watcher")
				sew.sentinel.Cleanup(log)
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
				"sentinel": sew.SentinelURI, "shard": rem.master.name,
				"redis_server": fmt.Sprintf("%s:%s", rem.master.ip, rem.master.port),
			},
		).Add(1)
	case "-failover-abort-no-good-slave":
		failoverAbortNoGoodSlaveCount.With(
			prometheus.Labels{
				"sentinel": sew.SentinelURI, "shard": rem.target.name,
			},
		).Add(1)
	case "+sdown":
		switch rem.target.role {
		case "sentinel":
			sdownSentinelCount.With(
				prometheus.Labels{
					"sentinel": sew.SentinelURI, "shard": rem.master.name,
					"redis_server": fmt.Sprintf("%s:%s", rem.target.ip, rem.target.port),
				},
			).Add(1)
		default:
			sdownCount.With(
				prometheus.Labels{
					"sentinel": sew.SentinelURI, "shard": rem.master.name,
					"redis_server": fmt.Sprintf("%s:%s", rem.target.ip, rem.target.port),
				},
			).Add(1)
		}
	case "-sdown":
		switch rem.target.role {
		case "sentinel":
			sdownClearedSentinelCount.With(
				prometheus.Labels{
					"sentinel": sew.SentinelURI, "shard": rem.master.name,
					"redis_server": fmt.Sprintf("%s:%s", rem.target.ip, rem.target.port),
				},
			).Add(1)
		default:
			sdownClearedCount.With(
				prometheus.Labels{
					"sentinel": sew.SentinelURI, "shard": rem.master.name,
					"redis_server": fmt.Sprintf("%s:%s", rem.target.ip, rem.target.port),
				},
			).Add(1)
		}
	}
}

func (sew *SentinelEventWatcher) initCounters() {
	if sew.Topology != nil {

		for _, shard := range *sew.Topology {
			failoverAbortNoGoodSlaveCount.With(
				prometheus.Labels{
					"sentinel": sew.SentinelURI, "shard": shard.Name,
				},
			).Add(0)

			for _, server := range shard.Servers {
				switchMasterCount.With(
					prometheus.Labels{
						"sentinel": sew.SentinelURI, "shard": shard.Name,
						"redis_server": fmt.Sprintf("%s:%s", server.CRUD.IP, server.CRUD.Port),
					},
				).Add(0)
				sdownSentinelCount.With(
					prometheus.Labels{
						"sentinel": sew.SentinelURI, "shard": shard.Name,
						"redis_server": fmt.Sprintf("%s:%s", server.CRUD.IP, server.CRUD.Port),
					},
				).Add(0)
				sdownCount.With(
					prometheus.Labels{
						"sentinel": sew.SentinelURI, "shard": shard.Name,
						"redis_server": fmt.Sprintf("%s:%s", server.CRUD.IP, server.CRUD.Port),
					},
				).Add(0)
				sdownClearedSentinelCount.With(
					prometheus.Labels{
						"sentinel": sew.SentinelURI, "shard": shard.Name,
						"redis_server": fmt.Sprintf("%s:%s", server.CRUD.IP, server.CRUD.Port),
					},
				).Add(0)
				sdownClearedCount.With(
					prometheus.Labels{
						"sentinel": sew.SentinelURI, "shard": shard.Name,
						"redis_server": fmt.Sprintf("%s:%s", server.CRUD.IP, server.CRUD.Port),
					},
				).Add(0)
			}
		}

	}
}

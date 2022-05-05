package events

import (
	"context"

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
		[]string{"sentinel", "shard", "redis_server", "role"},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(switchMasterCount, failoverAbortNoGoodSlaveCount, sdownCount)
}

// SentinelEventWatcher implements RunnableThread
var _ threads.RunnableThread = &SentinelEventWatcher{}

type SentinelEventWatcher struct {
	Instance      client.Object
	SentinelURI   string
	ExportMetrics bool
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

//Start starts metrics gatherer for sentinel
func (sew *SentinelEventWatcher) Start(parentCtx context.Context, l logr.Logger) error {
	log := l.WithValues("sentinel", sew.SentinelURI)
	if sew.started {
		log.Info("the event watcher is already running")
		return nil
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
func (fw *SentinelEventWatcher) Stop() {
	fw.cancel()
}

func (smg *SentinelEventWatcher) metricsFromEvent(rem RedisEventMessage) {
	switch rem.event {
	case "+switch-master":
		switchMasterCount.With(
			prometheus.Labels{
				"sentinel": smg.SentinelURI, "shard": rem.target.name,
			},
		).Add(1)
	case "-failover-abort-no-good-slave":
		failoverAbortNoGoodSlaveCount.With(
			prometheus.Labels{
				"sentinel": smg.SentinelURI, "shard": rem.target.name,
			},
		).Add(1)
	case "+sdown":
		sdownCount.With(
			prometheus.Labels{
				"sentinel": smg.SentinelURI, "shard": rem.master.name, "role": rem.target.role, "redis_server": rem.target.ip,
			},
		).Add(1)
	}
}

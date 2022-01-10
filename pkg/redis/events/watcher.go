package events

import (
	"context"
	"errors"
	"strings"

	"github.com/3scale/saas-operator/pkg/redis"
	"github.com/go-logr/logr"
	goredis "github.com/go-redis/redis/v8"
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
			Help:      "no-good-slave (https://redis.io/topics/sentinel#sentinel-api)",
		},
		[]string{"sentinel", "shard"},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(switchMasterCount, failoverAbortNoGoodSlaveCount)
}

type SentinelEventWatcher struct {
	Instance      client.Object
	SentinelURI   string
	Log           logr.Logger
	EventsCh      chan event.GenericEvent
	ExportMetrics bool
	started       bool
	cancel        context.CancelFunc
	sentinel      *redis.SentinelServer
}

// IsStarted returns whether the metrics gatherer is running or not
func (fw *SentinelEventWatcher) IsStarted() bool {
	return fw.started
}

//Start starts metrics gatherer for sentinel
func (fw *SentinelEventWatcher) Start(parentCtx context.Context) {
	log := fw.Log.WithValues("sentinel", fw.SentinelURI)
	if fw.started {
		log.Error(errors.New("already started"), "the event watcher is already running")
		return
	}

	go func() {
		var err error
		var ctx context.Context
		ctx, fw.cancel = context.WithCancel(parentCtx)

		fw.sentinel, err = redis.NewSentinelServerFromConnectionString(fw.SentinelURI, fw.SentinelURI)
		if err != nil {
			log.Error(err, "cannot create SentinelServer")
		}

		ch, closeWatch := fw.sentinel.CRUD.SentinelPSubscribe(ctx, "+switch-master", "-failover-abort-no-good-slave")
		defer closeWatch()

		log.Info("event watcher running")

		for {
			select {

			case msg := <-ch:
				log.V(1).Info("received event from sentinel", "event", msg.String())
				fw.EventsCh <- event.GenericEvent{Object: fw.Instance}
				if fw.ExportMetrics {
					fw.metricsFromEvent(msg)
				}

			case <-ctx.Done():
				log.Info("shutting down event watcher")
				fw.started = false
				return
			}
		}
	}()

	fw.started = true
}

// Stop stops the sentinel event watcher
func (fw *SentinelEventWatcher) Stop() {
	fw.cancel()
}

func (smg *SentinelEventWatcher) metricsFromEvent(msg *goredis.Message) {

	switch msg.Channel {
	case "+switch-master":
		shard := strings.Split(msg.Payload, " ")[0]
		switchMasterCount.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": shard}).Add(1)
	case "-failover-abort-no-good-slave":
		shard := strings.Split(msg.Payload, " ")[1]
		failoverAbortNoGoodSlaveCount.With(prometheus.Labels{"sentinel": smg.SentinelURI, "shard": shard}).Add(1)
	}
}

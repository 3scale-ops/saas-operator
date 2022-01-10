package events

import (
	"context"
	"strings"
	"sync"

	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// SentinelEvents is a struct that holds configuration to
// run sentinel event watchers
type SentinelEvents struct {
	mu              sync.Mutex
	sentinelEventCh chan event.GenericEvent
	watchers        map[string]*SentinelEventWatcher
}

// NewSentinelEvents returns a SentinelEvents struct
func NewSentinelEvents() SentinelEvents {
	return SentinelEvents{
		sentinelEventCh: make(chan event.GenericEvent),
		watchers:        map[string]*SentinelEventWatcher{},
	}
}

// RunEventWatcher runs a FailoverWatcher for the given instance.
func (se *SentinelEvents) RunEventWatcher(ctx context.Context, key string, instance client.Object,
	sentinelURI string, ch chan event.GenericEvent, metrics bool, log logr.Logger) {

	// run the exporter for this instance if it is not running, do nothing otherwise
	if w, ok := se.watchers[key]; !ok || !w.IsStarted() {
		se.mu.Lock()
		se.watchers[key] = &SentinelEventWatcher{
			Instance:      instance,
			SentinelURI:   sentinelURI,
			Log:           log,
			EventsCh:      ch,
			ExportMetrics: metrics,
		}
		se.watchers[key].Start(ctx)
		se.mu.Unlock()
	}
}

// StopExporter stops the event watcher for the given instance
func (se *SentinelEvents) StopEventWatcher(key string) {
	se.mu.Lock()
	se.watchers[key].Stop()
	delete(se.watchers, key)
	se.mu.Unlock()
}

//GetStatusChangeChannel returns the channel through which sentinel events can be received
func (se *SentinelEvents) GetSentinelEventsChannel() <-chan event.GenericEvent {
	return se.sentinelEventCh
}

// ReconcileEventWatchers ensures that all Sentinel instances have an event watcher monitoring
// sentinel events
func (se *SentinelEvents) ReconcileEventWatchers(ctx context.Context, instance client.Object, sentinelURIs []string, log logr.Logger) {

	shouldRun := map[string]int{}

	for _, uri := range sentinelURIs {
		key := util.ObjectKey(instance).String() + uri
		shouldRun[key] = 1
		se.RunEventWatcher(ctx, key, instance, uri, se.sentinelEventCh, true, log)
	}

	// Stop event watchers for any sentinel replica that does not exist anymore
	for key := range se.watchers {
		if strings.Contains(key, util.ObjectKey(instance).String()) {
			if _, ok := shouldRun[key]; !ok {
				se.StopEventWatcher(key)
			}
		}
	}
}

// CleanupEventWatchers stops the sentinel event watcher for this client.Object.
// This is used as a cleanup function in the finalize phase of the controller loop.
func (se *SentinelEvents) CleanupEventWatchers(instance client.Object) func() {
	return func() {
		prefix := util.ObjectKey(instance).String()
		for key := range se.watchers {
			if strings.Contains(key, prefix) {
				se.StopEventWatcher(key)
			}
		}
	}
}

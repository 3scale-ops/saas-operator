package threads

import (
	"context"
	"strings"
	"sync"

	"github.com/3scale/saas-operator/pkg/util"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type RunnableThread interface {
	GetID() string
	SetChannel(chan event.GenericEvent)
	Start(context.Context, logr.Logger) error
	Stop()
	IsStarted() bool
}

// SentinelEvents is a struct that holds configuration to
// run sentinel event watchers
type Manager struct {
	mu      sync.Mutex
	channel chan event.GenericEvent
	threads map[string]RunnableThread
}

// NewManager returns a new initialized Manager struct
func NewManager() Manager {
	return Manager{
		channel: make(chan event.GenericEvent),
		threads: map[string]RunnableThread{},
	}
}

// RunThread runs thread and associates it with a given key so it can later be stopped
func (mgr *Manager) RunThread(ctx context.Context, key string, thread RunnableThread, log logr.Logger) error {
	thread.SetChannel(mgr.channel)
	// run the exporter for this instance if it is not running, do nothing otherwise
	if w, ok := mgr.threads[key]; !ok || !w.IsStarted() {
		mgr.mu.Lock()
		mgr.threads[key] = thread
		if err := mgr.threads[key].Start(ctx, log); err != nil {
			return err
		}
		mgr.mu.Unlock()
	}
	return nil
}

// StopExporter stops the thread identified by the given key
func (mgr *Manager) StopThread(key string) {
	mgr.mu.Lock()
	if _, ok := mgr.threads[key]; !ok {
		return
	}
	mgr.threads[key].Stop()
	delete(mgr.threads, key)
	mgr.mu.Unlock()
}

// GetChannel returns the channel through which events can be received
// from the running thread
func (mgr *Manager) GetChannel() <-chan event.GenericEvent {
	return mgr.channel
}

// ReconcileThreads ensures that the threads identified by the provided keys are running. prefix() is used to identify
// which threads belong to each resource.
func (mgr *Manager) ReconcileThreads(ctx context.Context, instance client.Object, threads []RunnableThread, log logr.Logger) error {

	shouldRun := map[string]int{}

	for _, thread := range threads {
		key := prefix(instance) + thread.GetID()
		shouldRun[key] = 1
		if err := mgr.RunThread(ctx, key, thread, log); err != nil {
			return err
		}
	}

	// Stop threads that should not be running anymore
	for key := range mgr.threads {
		if strings.Contains(key, prefix(instance)) {
			if _, ok := shouldRun[key]; !ok {
				mgr.StopThread(key)
			}
		}
	}

	return nil
}

// CleanupThreads returns a function that cleans matching threads when invoked.
// This is intended for use as a cleanup function in the finalize phase of a controller's
// reconcile loop.
func (mgr *Manager) CleanupThreads(instance client.Object) func() {
	return func() {
		for key := range mgr.threads {
			if strings.Contains(key, prefix(instance)) {
				mgr.StopThread(key)
			}
		}
	}
}

func prefix(o client.Object) string {
	return util.ObjectKey(o).String() + "_"
}

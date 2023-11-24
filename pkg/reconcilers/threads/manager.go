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
	CanBeDeleted() bool
}

// Manager is a struct that holds configuration to
// manage concurrent RunnableThreads
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

// runThread runs thread and associates it with a given key so it can later be stopped
func (mgr *Manager) runThread(ctx context.Context, key string, thread RunnableThread, log logr.Logger) error {
	thread.SetChannel(mgr.channel)

	t, ok := mgr.threads[key]
	// do nothing if present and already started
	if ok && t.IsStarted() {
		return nil
	}

	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if !ok {
		mgr.threads[key] = thread
	}
	if err := mgr.threads[key].Start(ctx, log); err != nil {
		return err
	}
	return nil
}

// stopThread stops the thread identified by the given key
func (mgr *Manager) stopThread(key string) {
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
func (mgr *Manager) ReconcileThreads(ctx context.Context, owner client.Object, threads []RunnableThread, log logr.Logger) error {

	shouldRun := map[string]int{}

	for _, thread := range threads {
		key := prefix(owner) + thread.GetID()
		shouldRun[key] = 1
		if err := mgr.runThread(ctx, key, thread, log); err != nil {
			return err
		}
	}

	// Stop threads that should not be running anymore
	for key := range mgr.threads {
		if strings.Contains(key, prefix(owner)) {
			if _, ok := shouldRun[key]; !ok && mgr.threads[key].CanBeDeleted() {
				mgr.stopThread(key)
			}
		}
	}

	return nil
}

// CleanupThreads returns a function that cleans matching threads when invoked.
// This is intended for use as a cleanup function in the finalize phase of a controller's
// reconcile loop.
func (mgr *Manager) CleanupThreads(owner client.Object) func(context.Context, client.Client) error {
	return func(context.Context, client.Client) error {
		for key := range mgr.threads {
			if strings.Contains(key, prefix(owner)) {
				mgr.stopThread(key)
			}
		}
		return nil
	}
}

// Returns a thread, typically for inspection by the caller (ie get status/errors)
func (mgr *Manager) GetThread(id string, owner client.Object, log logr.Logger) RunnableThread {
	key := prefix(owner) + id
	if _, ok := mgr.threads[key]; ok {
		return mgr.threads[key]
	} else {
		return nil
	}
}

func (mgr *Manager) GetKeys() []string {
	keys := make([]string, 0, len(mgr.threads))
	for k, _ := range mgr.threads {
		keys = append(keys, k)
	}
	return keys
}

func prefix(o client.Object) string {
	return util.ObjectKey(o).String() + "_"
}

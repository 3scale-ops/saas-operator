package backup

import (
	"context"
	"fmt"
	"time"

	"github.com/3scale/saas-operator/pkg/redis/sharded"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Runner struct {
	shardName string
	server    *sharded.RedisServer
	timeout   time.Duration
	timestamp time.Time
	eventsCh  chan event.GenericEvent
	cancel    context.CancelFunc
	status    RunnerStatus
	instance  client.Object
}

type RunnerStatus struct {
	Started  bool
	Finished bool
	Error    error
}

func NewBackupRunner(shardName string, server *sharded.RedisServer, timestamp time.Time, timeout time.Duration, instance client.Object) *Runner {
	return &Runner{
		shardName: shardName,
		server:    server,
		timestamp: timestamp,
		timeout:   timeout,
		instance:  instance,
	}
}

func ID(shard, alias string, ts time.Time) string {
	return fmt.Sprintf("%s-%s-%d", shard, alias, ts.UTC().UnixMilli())
}

func (br *Runner) GetID() string {
	return ID(br.shardName, br.server.GetAlias(), br.timestamp)
}

// IsStarted returns whether the backup runner is started or not
func (br *Runner) IsStarted() bool {
	return br.status.Started
}

func (br *Runner) CanBeDeleted() bool {
	return time.Since(br.timestamp) > 1*time.Hour
}

func (br *Runner) SetChannel(ch chan event.GenericEvent) {
	br.eventsCh = ch
}

// Start starts the backup runner
func (br *Runner) Start(parentCtx context.Context, l logr.Logger) error {
	logger := l.WithValues("server", br.server.GetAlias(), "shard", br.shardName)

	var ctx context.Context
	ctx, br.cancel = context.WithCancel(parentCtx)
	ctx = log.IntoContext(ctx, logger)

	done := make(chan bool)
	errCh := make(chan error)

	// this go routine runs the backup
	go func() {
		if err := br.BackgroundSave(ctx); err != nil {
			errCh <- err
		}
		close(done)
	}()

	br.status = RunnerStatus{Started: true, Finished: false, Error: nil}
	logger.Info("backup running")

	// this goroutine listens controls the max time execution of the backup
	// and listens for status updates
	go func() {
		// apply a time boundary to the backup and listen for errors
		timer := time.NewTimer(br.timeout)
		for {
			select {

			case <-timer.C:
				err := fmt.Errorf("timeout reached (%v)", br.timeout)
				br.cancel()
				logger.Error(err, "backup failed")
				br.status.Finished = true
				br.status.Error = err
				br.eventsCh <- event.GenericEvent{Object: br.instance}
				return

			case err := <-errCh:
				logger.Error(err, "backup failed")
				br.status.Finished = true
				br.status.Error = err
				br.eventsCh <- event.GenericEvent{Object: br.instance}
				return

			case <-done:
				logger.V(1).Info("backup completed")
				br.status.Finished = true
				br.eventsCh <- event.GenericEvent{Object: br.instance}
				return
			}
		}
	}()

	return nil
}

// Stop stops the sentinel event watcher
func (br *Runner) Stop() {
	br.cancel()
}

func (br *Runner) Status() RunnerStatus {
	return br.status
}

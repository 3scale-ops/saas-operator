package backup

import (
	"context"
	"fmt"
	"time"

	"github.com/3scale/saas-operator/pkg/redis/sharded"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type BackupRunner struct {
	shardName string
	server    *sharded.RedisServer
	timeout   time.Duration
	// timestamp time.Time
	eventsCh chan event.GenericEvent
	started  bool
	cancel   context.CancelFunc
}

func (br *BackupRunner) GetID() string {
	return br.shardName
}

// IsStarted returns whether the backup runner is started or not
func (br *BackupRunner) IsStarted() bool {
	return br.started
}

func (br *BackupRunner) SetChannel(ch chan event.GenericEvent) {
	br.eventsCh = ch
}

// Start starts the backup runner
func (br *BackupRunner) Start(parentCtx context.Context, l logr.Logger) error {
	log := l.WithValues("server", br.server.GetAlias())

	var ctx context.Context
	ctx, br.cancel = context.WithCancel(parentCtx)
	defer br.cancel()

	var done chan bool
	var errCh chan error
	done, errCh = br.RunBackup(ctx)
	log.Info("backup running")

	// apply a time boundary to the backup and listen for errors
	timer := time.NewTimer(br.timeout)
	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timeout reached (%v)", br.timeout)

		case <-done:
			return nil

		case err := <-errCh:
			return err

		}
	}

}

// Stop stops the sentinel event watcher
func (br *BackupRunner) Stop() {
	br.cancel()
}

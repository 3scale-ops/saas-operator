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
	Instance           client.Object
	ShardName          string
	Server             *sharded.RedisServer
	Timeout            time.Duration
	PollInterval       time.Duration
	Timestamp          time.Time
	RedisDBFile        string
	SSHUser            string
	SSHKey             string
	SSHPort            uint32
	S3Bucket           string
	S3Path             string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSS3Endpoint      *string
	eventsCh           chan event.GenericEvent
	cancel             context.CancelFunc
	status             RunnerStatus
}

type RunnerStatus struct {
	Started    bool
	Finished   bool
	Error      error
	BackupFile string
	BackupSize int64
	FinishedAt time.Time
}

func ID(shard, alias string, ts time.Time) string {
	return fmt.Sprintf("%s-%s-%d", shard, alias, ts.UTC().UnixMilli())
}

func (br *Runner) GetID() string {
	return ID(br.ShardName, br.Server.GetAlias(), br.Timestamp)
}

// IsStarted returns whether the backup runner is started or not
func (br *Runner) IsStarted() bool {
	return br.status.Started
}

func (br *Runner) CanBeDeleted() bool {
	return time.Since(br.Timestamp) > 1*time.Hour
}

func (br *Runner) SetChannel(ch chan event.GenericEvent) {
	br.eventsCh = ch
}

// Start starts the backup runner
func (br *Runner) Start(parentCtx context.Context, l logr.Logger) error {
	logger := l.WithValues("server", br.Server.GetAlias(), "shard", br.ShardName)

	var ctx context.Context
	ctx, br.cancel = context.WithCancel(parentCtx)
	ctx = log.IntoContext(ctx, logger)

	done := make(chan bool)
	errCh := make(chan error)

	// this go routine runs the backup
	go func() {
		if err := br.BackgroundSave(ctx); err != nil {
			errCh <- err
			return
		}
		if err := br.UploadBackup(ctx); err != nil {
			errCh <- err
			return
		}
		if err := br.TagBackup(ctx); err != nil {
			errCh <- err
			return
		}
		close(done)
	}()

	br.status = RunnerStatus{Started: true, Finished: false, Error: nil}
	logger.Info("backup running")

	// this goroutine controls the max time execution of the backup
	// and listens for status updates
	go func() {
		// apply a time boundary to the backup and listen for errors
		timer := time.NewTimer(br.Timeout)
		for {
			select {

			case <-timer.C:
				err := fmt.Errorf("timeout reached (%v)", br.Timeout)
				br.cancel()
				logger.Error(err, "backup failed")
				br.status.Finished = true
				br.status.Error = err
				br.eventsCh <- event.GenericEvent{Object: br.Instance}
				return

			case err := <-errCh:
				logger.Error(err, "backup failed")
				br.status.Finished = true
				br.status.Error = err
				br.eventsCh <- event.GenericEvent{Object: br.Instance}
				return

			case <-done:
				logger.V(1).Info("backup completed")
				br.status.Finished = true
				br.status.BackupFile = fmt.Sprintf("s3://%s/%s/%s", br.S3Bucket, br.S3Path, br.BackupFileCompressed())
				br.status.FinishedAt = time.Now()
				br.eventsCh <- event.GenericEvent{Object: br.Instance}
				br.publishMetrics()
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

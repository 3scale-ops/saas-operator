package backup

import (
	"context"
	"fmt"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (br *Runner) BackgroundSave(ctx context.Context) error {
	logger := log.FromContext(ctx, "function", "(br *Runner) BackgroundSave")
	prevsave, err := br.Server.RedisLastSave(ctx)
	if err != nil {
		logger.Error(errLastSave(err), "backup error")
		return errLastSave(err)
	}

	err = br.Server.RedisBGSave(ctx)
	if err != nil {
		// TODO: need to hanlde the case that a save is already running
		logger.Error(errBGSave(err), "backup error")
		return errBGSave(err)
	}
	logger.V(1).Info("BGSave running")

	ticker := time.NewTicker(br.PollInterval)

	// wait until BGSAVE completes
	for {
		select {
		case <-ticker.C:
			lastsave, err := br.Server.RedisLastSave(ctx)
			if err != nil {
				// retry at next tick
				logger.Error(errLastSave(err), "transient backup error")
				continue
			}
			if lastsave > prevsave {
				// BGSAVE completed
				logger.Info("BGSave finished")
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("context cancelled")
		}
	}
}

func errBGSave(err error) error {
	return fmt.Errorf("redis cmd (BGSAVE) error: %w", err)
}

func errLastSave(err error) error {
	return fmt.Errorf("redis cmd (LASTSAVE) error: %w", err)
}

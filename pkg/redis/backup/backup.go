package backup

import (
	"context"
	"time"
)

func (br BackupRunner) RunBackup(ctx context.Context) (chan bool, chan error) {
	done := make(chan bool)
	errCh := make(chan error)

	go func() {
		if err := br.BackgroundSave(ctx); err != nil {
			errCh <- err
		}
		close(done)

	}()

	return done, errCh
}

func (br BackupRunner) BackgroundSave(ctx context.Context) error {
	lastsave, err := br.server.RedisLastSave(ctx)
	if err != nil {
		return err
	}

	err = br.server.RedisBGSave(ctx)
	if err != nil {
		// TODO: need to hanlde the case that a save is already running
		return err
	}

	ticker := time.NewTicker(60 * time.Second)
	done := make(chan bool)

	// wait until BGSAVE completes
	for {
		select {
		case <-ticker.C:
			seconds, err := br.server.RedisLastSave(ctx)
			if err != nil {
				// retry in next ticker
				continue
			}
			if seconds < lastsave {
				// BGSAVE completed
				close(done)
			}
		case <-done:
			return nil
		}
	}

}

func (br BackupRunner) UploadBackup(ctx context.Context) error {
	return nil
}

func (br BackupRunner) TagBackup(ctx context.Context) error {
	return nil
}

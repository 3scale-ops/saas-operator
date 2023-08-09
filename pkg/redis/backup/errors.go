package backup

import "fmt"

func errBGSave(err error) error {
	return fmt.Errorf("redis cmd (BGSAVE) error: %w", err)
}

func errLastSave(err error) error {
	return fmt.Errorf("redis cmd (LASTSAVE) error: %w", err)
}

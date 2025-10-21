package pyupdater

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path"
)

func removeLockFile(directory string) error {
	log.Info().Msg("Removing uv.lock file")
	if err := os.Remove(path.Join(directory, "uv.lock")); err != nil {
		if os.IsNotExist(err) {
			log.Info().
				Err(err).
				Msg("uv.lock file does not exist, skipping removal")
		} else {
			return fmt.Errorf("failed to remove uv.lock: %w", err)
		}
	}
	return nil
}

// Update or generate the lock file
func updateLockFile(directory string) error {
	return executeCommand([]string{"uv", "sync", "--directory", directory})
}

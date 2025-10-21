package pyupdater

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
	"strings"
)

func executeCommand(command []string) error {
	log.Debug().
		Strs("command", command).
		Msg("Executing shell command")
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command '%s' failed: %w\nOutput: %s", strings.Join(command, " "), err, string(output))
	}
	log.Debug().
		Str("output", string(output)).
		Msg("Command executed successfully")
	return nil
}

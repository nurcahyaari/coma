package file

import (
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func NewDir(fd string) error {
	cmd := exec.Command("mkdir", "-m", "0777", "-p", fd)
	cmd.Env = append(os.Environ(), "SUDO_COMMAND=true")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Error().Err(err).
			Str("path", fd).
			Msg("creating file directory")
		return err
	}

	cmd = exec.Command("chmod", "755", fd)
	cmd.Env = append(os.Environ(), "SUDO_COMMAND=true")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Error().Err(err).
			Str("path", fd).
			Msg("creating access file directory")
		return err
	}
	return nil
}

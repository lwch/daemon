package background

import (
	"os"
	"os/exec"
)

func makeCommand(username string, arg ...string) *exec.Cmd {
	cmd = exec.Command(os.Args[0], arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}

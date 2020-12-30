package background

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// Start start daemon
func Start(pid string, arg ...string) {
	cmd := exec.Command(os.Args[0], arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err == nil {
		if len(pid) > 0 {
			writePidFile(pid, cmd.Process.Pid)
		}
		cmd.Process.Release()
		os.Exit(0)
	}
}

func writePidFile(dir string, pid int) {
	ioutil.WriteFile(dir, []byte(fmt.Sprintf("%d", pid)), 0644)
}

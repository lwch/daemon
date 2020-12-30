package background

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var cmd *exec.Cmd

// Start start daemon
func Start(pid string, arg ...string) {
	if len(pid) > 0 {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			<-c
			os.Remove(pid)
		}()
	}
	for {
		run(pid, arg...)
	}
}

func run(pid string, arg ...string) {
	cmd := exec.Command(os.Args[0], arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err == nil {
		if len(pid) > 0 {
			writePidFile(pid, cmd.Process.Pid)
		}
		cmd.Process.Release()
		cmd.Wait()
	}
}

func writePidFile(dir string, pid int) {
	ioutil.WriteFile(dir, []byte(fmt.Sprintf("%d", pid)), 0644)
}

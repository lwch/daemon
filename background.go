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
	chExit := make(chan struct{})
	if len(pid) > 0 {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			sig := <-c
			cmd.Process.Signal(sig)
			<-chExit
			os.Remove(pid)
		}()
	}
	for {
		run(chExit, pid, arg...)
	}
}

func run(ch chan struct{}, pid string, arg ...string) {
	cmd = exec.Command(os.Args[0], arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err == nil {
		if len(pid) > 0 {
			writePidFile(pid, cmd.Process.Pid)
		}
		cmd.Wait()
		ch <- struct{}{}
	} else {
		os.Exit(1)
	}
}

func writePidFile(dir string, pid int) {
	ioutil.WriteFile(dir, []byte(fmt.Sprintf("%d", pid)), 0644)
}

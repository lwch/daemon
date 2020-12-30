package daemon

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

var cmd *exec.Cmd

// Start start daemon
func Start(pid, username string, arg ...string) {
	chExit := make(chan struct{})
	onExit := false
	var wg sync.WaitGroup
	if len(pid) > 0 {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		wg.Add(1)
		go func() {
			defer wg.Done()
			sig := <-c
			onExit = true
			cmd.Process.Signal(sig)
			<-chExit
			os.Remove(pid)
			os.Exit(0)
		}()
	}
	for {
		run(chExit, pid, username, arg...)
		if onExit {
			break
		}
	}
	wg.Wait()
}

func run(ch chan struct{}, pid, username string, arg ...string) {
	cmd = makeCommand(username, arg...)
	if err := cmd.Start(); err == nil {
		writePidFile(pid, os.Getpid())
		cmd.Wait()
		ch <- struct{}{}
	} else {
		fmt.Println("create child process failed")
		os.Exit(1)
	}
}

func writePidFile(dir string, pid int) {
	ioutil.WriteFile(dir, []byte(fmt.Sprintf("%d", pid)), 0644)
}

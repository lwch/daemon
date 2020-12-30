package background

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"syscall"
)

var cmd *exec.Cmd

// Start start daemon
func Start(pid, username string, arg ...string) {
	chExit := make(chan struct{})
	onExit := false
	if len(pid) > 0 {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
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
			return
		}
	}
}

func run(ch chan struct{}, pid, username string, arg ...string) {
	cmd = exec.Command(os.Args[0], arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if len(user) > 0 {
		u, err := user.Lookup(username)
		if err != nil {
			fmt.Println("user not found")
			os.Exit(1)
		}
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uid,
			},
		}
	}
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

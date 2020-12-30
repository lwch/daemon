// +build !windows

package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func makeCommand(username string, arg ...string) *exec.Cmd {
	cmd = exec.Command(os.Args[0], arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if len(username) > 0 {
		u, err := user.Lookup(username)
		if err != nil {
			fmt.Println("user not found")
			os.Exit(1)
		}
		uid, _ := strconv.ParseUint(u.Uid, 10, 32)
		gid, _ := strconv.ParseUint(u.Gid, 10, 32)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(uid),
				Gid: uint32(gid),
			},
		}
	}
	return cmd
}

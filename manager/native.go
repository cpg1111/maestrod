package manager

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Native runs maestro in it's own native vanilla container
type Native struct {
	Driver
}

// Run runs the native driver
func (n Native) Run() error {
	maestro, lookErr := exec.LookPath("maestro")
	if lookErr != nil {
		return lookErr
	}
	cmd := exec.Command("./container", maestro)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("ERROR", err)
		return err
	}
}

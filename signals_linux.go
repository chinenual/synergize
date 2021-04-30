// +build !darwin, linux, !windows

package main

import (
	"os"
	"os/exec"
	"syscall"
)

var ignoredSignals = []os.Signal{syscall.SIGURG}

func SetCmdNoConsoleWindow(cmd *exec.Cmd) {
	// only necessary on windows
}

// +build windows, !darwin, !linux

package main

import (
	"os"
	"os/exec"
	"syscall"
)

var ignoredSignals = []os.Signal{}

func SetCmdNoConsoleWindow(cmd *exec.Cmd) {
	(*cmd).SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

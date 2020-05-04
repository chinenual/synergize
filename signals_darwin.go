// +build darwin, !linux, !windows

package main

import (
	"os"
	"syscall"
)

var ignoredSignals = []os.Signal{syscall.SIGURG}

// +build windows, !darwin, !linux

package main

import (
	"os"
)

var ignoredSignals = []os.Signal{}

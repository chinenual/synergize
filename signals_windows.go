// +build windows, !darwin

package main

import (
	"os"
)

var ignoredSignals = []os.Signal{}

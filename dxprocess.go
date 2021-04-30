package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chinenual/synergize/logger"
)

func getExeDirectory() (path string) {
	var err error
	if path, err = os.Executable(); err != nil {
		logger.Errorf("Could not determine the location of the executable: %v\n", err)
		path = ""
	}
	path = filepath.Dir(path)
	logger.Infof("ExeDirectory: %s\n", path)
	return
}

var dxcmd *exec.Cmd

func dx2synProcessStart(path string) (err error) {
	exeName := "dx2syn"
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}
	exePath := filepath.Join(getExeDirectory(), exeName)

	dxcmd = exec.Command(exePath, "-makecrt", path)
	// suppress the console window on windows
	SetCmdNoConsoleWindow(dxcmd)

	var stderr io.ReadCloser
	if stderr, err = dxcmd.StderrPipe(); err != nil {
		return
	}
	if err = dxcmd.Start(); err != nil {
		return
	}

	go _slurpLog(stderr)

	return
}

func dx2SynProcessCancel() (err error) {
	if dxcmd != nil && dxcmd.Process != nil {
		if err = dxcmd.Process.Kill(); err != nil {
			return
		}
		_finishToUI("Cancelled!")
	}
	return
}

func _slurpLog(pipe io.ReadCloser) {
	for {
		r := bufio.NewReader(pipe)
		var line []byte
		var err error
		line, err = r.ReadBytes('\n')
		_logToUI(string(line))
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				logger.Errorf("error reading subprocess pipe: %v\n", err)
			}
			break
		}
	}
	dxcmd.Wait()
	msg := "Success"
	if !dxcmd.ProcessState.Success() {
		msg = fmt.Sprintf("dx2syn returned error status: %d", dxcmd.ProcessState.ExitCode())
	}
	_finishToUI(msg)
}

var astilectronWindow *astilectron.Window

func dx2synRegisterBridge(w *astilectron.Window) (err error) {
	astilectronWindow = w
	return
}

func _finishToUI(msg string) (err error) {
	logger.Infof("dx2syn finished: %s\n", strings.TrimSpace(msg))
	if astilectronWindow == nil {
		return
	}
	var strval string
	if err = bootstrap.SendMessage(astilectronWindow, "dx2synFinish", msg,
		func(m *bootstrap.MessageIn) {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &strval); err != nil {
				logger.Errorf(" _finishToUI failed to decode json response : %v\n", err)
				return
			}
		}); err != nil {
		return
	}
	return
}

func _logToUI(line string) (err error) {
	logger.Infof("dx2syn: %s\n", strings.TrimSpace(line))
	if astilectronWindow == nil {
		return
	}
	var strval string
	if err = bootstrap.SendMessage(astilectronWindow, "dx2synAddProcessLog", line,
		func(m *bootstrap.MessageIn) {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &strval); err != nil {
				logger.Errorf(" _logToUI failed to decode json response : %v\n", err)
				return
			}
		}); err != nil {
		return
	}
	return
}

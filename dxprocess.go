package main

import (
	"bufio"
	"fmt"
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

func _finishToUI(msg string) (err error) {
	logger.Infof("dx2syn finished: %s\n", strings.TrimSpace(msg))
	if err = dx2synFinish(msg); err != nil {
		logger.Error("error finishing log to frontend", err)
	}
	return
}

func _logToUI(line string) (err error) {
	logger.Infof("dx2syn: %s\n", strings.TrimSpace(line))
	if err = dx2synAddProcessLog(line); err != nil {
		logger.Error("error sending log to frontend", err)
	}
	return
}

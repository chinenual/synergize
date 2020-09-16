package logger

import (
	"io"
	"log"
	"os"

	"github.com/asticode/go-astikit"

	"gopkg.in/natefinch/lumberjack.v2"
)

// a global logger implementing the astikit.StdLogger and astikit.SeverityLogger interfaces
type Logger interface {
	astikit.StdLogger
	astikit.SeverityLogger
}

// Levels
type Level int

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

type synLogger struct {
	level Level
	l     *log.Logger
}

// initialize with safe defaults so that its safe to call the log functions
// before a call to Init() (e.g. when the preferences are loaded in order to initialize
// command line flag defaults)
var l = synLogger{
	level: LevelInfo,
	l:     log.New(os.Stdout, "", 0),
}

func GetLogger() Logger {
	return &l
}

func Init(logPath string, level Level) {
	multi := io.MultiWriter(
		&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    5, // megabytes
			MaxBackups: 2,
			Compress:   false,
		},
		os.Stderr)
	log.SetOutput(multi)

	l.level = level
	l.l = log.New(log.Writer(), log.Prefix(), log.Flags()) //|log.Lshortfile)
}

func (l *synLogger) SetRaw() {
	l.l.SetFlags(0)
}

func (l *synLogger) Print(v ...interface{}) {
	l.Info(v...)
}

func (l *synLogger) Printf(format string, v ...interface{}) {
	l.Infof(format, v...)
}

func (l *synLogger) Debug(v ...interface{}) {
	if l.level <= LevelDebug {
		newv := append([]interface{}{"DEBUG: "}, v...)
		l.l.Print(newv...)
	}
}

func (l *synLogger) Debugf(format string, v ...interface{}) {
	if l.level <= LevelDebug {
		l.l.Printf("DEBUG: "+format, v...)
	}
}

func (l *synLogger) Info(v ...interface{}) {
	if l.level <= LevelInfo {
		newv := append([]interface{}{"INFO: "}, v...)
		l.l.Print(newv...)
	}
}

func (l *synLogger) Infof(format string, v ...interface{}) {
	if l.level <= LevelInfo {
		l.l.Printf("INFO: "+format, v...)
	}
}

func (l *synLogger) Warn(v ...interface{}) {
	if l.level <= LevelWarn {
		newv := append([]interface{}{"WARN: "}, v...)
		l.l.Print(newv...)
	}
}

func (l *synLogger) Warnf(format string, v ...interface{}) {
	if l.level <= LevelWarn {
		l.l.Printf("WARN: "+format, v...)
	}
}

func (l *synLogger) Error(v ...interface{}) {
	if l.level <= LevelError {
		newv := append([]interface{}{"ERROR: "}, v...)
		l.l.Print(newv...)
	}
}

func (l *synLogger) Errorf(format string, v ...interface{}) {
	if l.level <= LevelError {
		l.l.Printf("ERROR: "+format, v...)
	}
}

/// --- convenience variants that call the global instance:

func SetRaw() {
	l.SetRaw()
}

func Print(v ...interface{}) {
	l.Info(v...)
}

func Printf(format string, v ...interface{}) {
	l.Infof(format, v...)
}

func Debug(v ...interface{}) {
	l.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	l.Debugf(format, v...)
}

func Info(v ...interface{}) {
	l.Info(v...)
}

func Infof(format string, v ...interface{}) {
	l.Infof(format, v...)
}

func Warn(v ...interface{}) {
	l.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	l.Warnf(format, v...)
}

func Error(v ...interface{}) {
	l.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	l.Errorf(format, v...)
}

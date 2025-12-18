package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

//Pre-wails3 rewrite this used astilogger interfaces; Wails uses new go standard log/slog  - so we do too.

func InitViaString(logPath string, wailsLogPath string, logLevelString string) {
	var level slog.Level
	var levelMsg = ""
	switch logLevelString {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
		levelMsg = fmt.Sprintf("Invalid value for -loglevel (%s).  Defaulting to INFO\n", logLevelString)
	}
	Init(logPath, wailsLogPath, level)
	if levelMsg != "" {
		slog.Warn(levelMsg)
	}
}

var SynergizeLogger *slog.Logger
var SynergizeWailsLogger *slog.Logger

func Init(logPath string, wailsLogPath string, level slog.Level) {
	var destSynergize io.Writer = os.Stdout
	var destWails io.Writer = os.Stdout
	if logPath != "" {
		multiSynergize := io.MultiWriter(
			&lumberjack.Logger{
				Filename:   logPath,
				MaxSize:    5, // megabytes
				MaxBackups: 2,
				Compress:   false,
			},
			os.Stderr)
		multiWails := io.MultiWriter(
			&lumberjack.Logger{
				Filename:   wailsLogPath,
				MaxSize:    5, // megabytes
				MaxBackups: 2,
				Compress:   false,
			},
			os.Stderr)
		destSynergize = multiSynergize
		destWails = multiWails
	}
	SynergizeWailsLogger = slog.New(NewSimpleHandler(destWails, nil))
	SynergizeLogger = slog.New(NewSimpleHandler(destSynergize, nil))
	slog.SetDefault(SynergizeLogger)
	// odd slog API - logger level effects "all" slog loggers?
	slog.SetLogLoggerLevel(level)
}

func Printf(format string, v ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	msg := ""
	for _, ele := range v {
		msg += fmt.Sprintf("%v ", ele)
	}
	slog.Debug(msg)
}
func Debugf(format string, v ...interface{}) {
	if strings.HasSuffix(format, "\n") {
		format = format[:len(format)-1]
	}
	slog.Debug(fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	msg := ""
	for _, ele := range v {
		msg += fmt.Sprintf("%v ", ele)
	}
	slog.Info(msg)
}
func Infof(format string, v ...interface{}) {
	if strings.HasSuffix(format, "\n") {
		format = format[:len(format)-1]
	}
	slog.Info(fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	msg := ""
	for _, ele := range v {
		msg += fmt.Sprintf("%v ", ele)
	}
	slog.Warn(msg)
}
func Warnf(format string, v ...interface{}) {
	if strings.HasSuffix(format, "\n") {
		format = format[:len(format)-1]
	}
	slog.Warn(fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	msg := ""
	for _, ele := range v {
		msg += fmt.Sprintf("%v ", ele)
	}
	slog.Error(msg)
}
func Errorf(format string, v ...interface{}) {
	if strings.HasSuffix(format, "\n") {
		format = format[:len(format)-1]
	}
	slog.Error(fmt.Sprintf(format, v...))
}

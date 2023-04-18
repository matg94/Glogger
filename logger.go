package glogger

import (
	"time"
)

var colors = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"reset":   "\033[0m",
}

type Logger interface {
	Error(Loggable)
	Info(Loggable)
	Debug(Loggable)
}

type Loggable interface {
	Error() string
}

type ConsoleLogger struct {
	Config  *LogFormatConfig
	LogFunc func(...any)
	Now     func() time.Time
}

func CreateConsoleLogger(dateFormat, logFormat string, levelColors map[string]string, logFunc func(...any)) *ConsoleLogger {
	return &ConsoleLogger{
		Config: &LogFormatConfig{
			DateFormat:  dateFormat,
			Format:      logFormat,
			LevelColor:  true,
			LevelColors: levelColors,
		},
		LogFunc: logFunc,
		Now:     time.Now().Local,
	}
}

func CreateFormattedLogger(logFormat string, logFunc func(...any)) *ConsoleLogger {
	return &ConsoleLogger{
		Config: &LogFormatConfig{
			DateFormat: "[2006-01-02 15:04:05]",
			Format:     logFormat,
			LevelColor: true,
			LevelColors: map[string]string{
				"error": "red",
				"info":  "blue",
				"debug": "yellow",
			},
		},
		LogFunc: logFunc,
		Now:     time.Now,
	}
}

func CreateSimpleConsoleLogger(logFunc func(...any)) *ConsoleLogger {
	return &ConsoleLogger{
		Config: &LogFormatConfig{
			DateFormat: "[2006-01-02 15:04:05]",
			Format:     "[green][date][reset] [yellow]->[reset] [level]: [log]",
			LevelColor: true,
			LevelColors: map[string]string{
				"error": "red",
				"info":  "blue",
				"debug": "yellow",
			},
		},
		LogFunc: logFunc,
		Now:     time.Now,
	}
}

func (l *ConsoleLogger) Error(log Loggable) {
	l.log(log, "error")
}

func (l *ConsoleLogger) Info(log Loggable) {
	l.log(log, "info")
}

func (l *ConsoleLogger) Debug(log Loggable) {
	l.log(log, "debug")
}

func (l *ConsoleLogger) log(log Loggable, level string) {
	l.LogFunc(FormatLog(level, log, l.Config, l.Now()))
}

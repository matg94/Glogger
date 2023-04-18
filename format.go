package glogger

import (
	"regexp"
	"strings"
	"time"
)

type LogFormatConfig struct {
	DateFormat  string
	Format      string
	LevelColor  bool
	LevelColors map[string]string
}

func FormatLog(level string, message Loggable, config *LogFormatConfig, time time.Time) string {
	re := regexp.MustCompile(`\[(.*?)\]`)
	modifiedStr := re.ReplaceAllStringFunc(config.Format, func(cmd string) string {
		return formatParser(cmd, level, message.Error(), config, time)
	})
	return modifiedStr + colors["reset"]
}

func Colorize(message, color string) string {
	return colors[color] + message + colors["reset"]
}

func formatParser(cmd, level, log string, config *LogFormatConfig, time time.Time) string {
	switch cmd {
	case "[level]":
		if config.LevelColor {
			return Colorize(strings.ToUpper(level), config.LevelColors[level])
		}
		return strings.ToUpper(level)
	case "[date]":
		return time.Format(config.DateFormat)
	case "[log]":
		return log
	case "[red]":
		return colors["red"]
	case "[green]":
		return colors["green"]
	case "[magenta]":
		return colors["magenta"]
	case "[blue]":
		return colors["blue"]
	case "[yellow]":
		return colors["yellow"]
	case "[reset]":
		return colors["reset"]
	default:
		return cmd
	}
}

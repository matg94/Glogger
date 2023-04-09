package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}

func fixedNow() time.Time {
	return time.Date(2023, 2, 2, 15, 15, 15, 0, time.UTC)
}

func TestDefaultLogger(t *testing.T) {
	buffer := new(bytes.Buffer)

	logger := CreateConsoleLogger(
		"2006-01-02",
		"[date]: [log]",
		map[string]string{
			"error": "red",
		},
		func(args ...interface{}) {
			fmt.Fprintln(buffer, args...)
		},
	)

	log := Str("test message")

	logger.Error(log)
	expected := fmt.Sprintf("%s: test message%s\n", time.Now().Format("2006-01-02"), colors["reset"])
	assertEqual(t, expected, buffer.String())
	buffer.Reset()
}

func TestSimpleLogger(t *testing.T) {
	buffer := new(bytes.Buffer)

	logger := CreateSimpleConsoleLogger(func(args ...interface{}) {
		fmt.Fprintln(buffer, args...)
	})

	log := Str("test message")

	logger.Error(log)
	expected := fmt.Sprintf("%s%s%s %s->%s %sERROR%s: test message%s\n",
		colors["green"],
		time.Now().Format("[2006-01-02 15:04:05]"),
		colors["reset"],
		colors["yellow"],
		colors["reset"],
		colors["red"],
		colors["reset"],
		colors["reset"],
	)
	assertEqual(t, expected, buffer.String())
	buffer.Reset()
}
func TestLoggerFormatter(t *testing.T) {

	config := &LogFormatConfig{
		DateFormat: "2023-02-02 15:15:15",
		Format:     "[date] [level]: [log] [default]",
		LevelColors: map[string]string{
			"error": "red",
			"info":  "blue",
			"debug": "yellow",
		},
	}

	buffer := new(bytes.Buffer)
	logger := &ConsoleLogger{
		Config: config,
		LogFunc: func(args ...interface{}) {
			fmt.Fprintln(buffer, args...)
		},
		Now: fixedNow,
	}

	log := Str("test message")

	// Test error log
	logger.Error(log)
	expected := fmt.Sprintf("2023-02-02 15:15:15 %sERROR%s: test message [default]%s\n", colors["red"], colors["reset"], colors["reset"])
	assertEqual(t, expected, buffer.String())
	buffer.Reset()

	// Test info log
	logger.Info(log)
	expected = fmt.Sprintf("2023-02-02 15:15:15 %sINFO%s: test message [default]%s\n", colors["blue"], colors["reset"], colors["reset"])
	assertEqual(t, expected, buffer.String())
	buffer.Reset()

	// Test debug log
	logger.Debug(log)
	expected = fmt.Sprintf("2023-02-02 15:15:15 %sDEBUG%s: test message [default]%s\n", colors["yellow"], colors["reset"], colors["reset"])
	assertEqual(t, expected, buffer.String())
	buffer.Reset()

}

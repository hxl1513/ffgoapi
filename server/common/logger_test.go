package common

import (
	"bytes"
	"strings"
	"testing"
)

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestNewLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, "info")
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}
}

func TestLoggerLevels(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, "debug")

	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
	logger.Debug("debug message")

	output := buf.String()
	if !contains(output, "info message") {
		t.Errorf("Expected log to contain 'info message', got: %s", output)
	}
	if !contains(output, "warn message") {
		t.Errorf("Expected log to contain 'warn message', got: %s", output)
	}
	if !contains(output, "error message") {
		t.Errorf("Expected log to contain 'error message', got: %s", output)
	}
	if !contains(output, "debug message") {
		t.Errorf("Expected log to contain 'debug message', got: %s", output)
	}
}

func TestLoggerInfoLevelFiltersDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, "info")

	logger.Debug("should not appear")

	if contains(buf.String(), "should not appear") {
		t.Error("Debug message should not appear at info level")
	}
}

func TestLoggerErrorLevelFiltersInfoAndDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, "error")

	logger.Info("info should not appear")
	logger.Debug("debug should not appear")
	logger.Error("error should appear")

	output := buf.String()
	if contains(output, "info should not appear") {
		t.Error("Info message should not appear at error level")
	}
	if contains(output, "debug should not appear") {
		t.Error("Debug message should not appear at error level")
	}
	if !contains(output, "error should appear") {
		t.Error("Error message should appear at error level")
	}
}

func TestLoggerWarnLevelFiltersDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, "warn")

	logger.Debug("debug should not appear")
	logger.Warn("warn should appear")

	output := buf.String()
	if contains(output, "debug should not appear") {
		t.Error("Debug message should not appear at warn level")
	}
	if !contains(output, "warn should appear") {
		t.Error("Warn message should appear at warn level")
	}
}

func TestLoggerFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, "info")

	logger.Info("hello world")

	output := buf.String()
	if !strings.HasPrefix(output, "[INFO]") {
		t.Errorf("Expected log to start with [INFO], got: %s", output)
	}
	if !strings.Contains(output, "hello world") {
		t.Errorf("Expected log to contain 'hello world', got: %s", output)
	}
}

func TestLoggerThreadSafety(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, "debug")

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			logger.Info("concurrent message")
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
	// If no panic, thread safety is working
}

package common

import (
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	level       string
	mu          sync.Mutex
}

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

var levelOrder = map[string]int{
	LogLevelDebug: 0,
	LogLevelInfo:  1,
	LogLevelWarn:  2,
	LogLevelError: 3,
}

// NewLogger creates a new Logger writing to the given io.Writer at the specified level.
func NewLogger(output io.Writer, level string) *Logger {
	prefix := func(lvl string) string {
		return "[" + lvl + "]"
	}

	return &Logger{
		infoLogger:  log.New(output, prefix("INFO"), 0),
		warnLogger:  log.New(output, prefix("WARN"), 0),
		errorLogger: log.New(output, prefix("ERROR"), 0),
		debugLogger: log.New(output, prefix("DEBUG"), 0),
		level:       level,
	}
}

func (l *Logger) shouldLog(level string) bool {
	lvl, ok := levelOrder[level]
	if !ok {
		return false
	}
	configuredLvl, ok := levelOrder[l.level]
	if !ok {
		return true
	}
	return lvl >= configuredLvl
}

func (l *Logger) formatMessage(msg string, args ...any) string {
	if len(args) > 0 {
		return msg + " " + formatArgs(args...)
	}
	return msg
}

func formatArgs(args ...any) string {
	var sb []byte
	for i, arg := range args {
		if i > 0 {
			sb = append(sb, ' ')
		}
		sb = append(sb, []byte(formatValue(arg))...)
	}
	return string(sb)
}

func formatValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return ""
	}
}

// Info logs at info level.
func (l *Logger) Info(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.shouldLog(LogLevelInfo) {
		l.infoLogger.Printf("%s | %s", timestamp(), l.formatMessage(msg, args...))
	}
}

// Warn logs at warn level.
func (l *Logger) Warn(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.shouldLog(LogLevelWarn) {
		l.warnLogger.Printf("%s | %s", timestamp(), l.formatMessage(msg, args...))
	}
}

// Error logs at error level.
func (l *Logger) Error(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.shouldLog(LogLevelError) {
		l.errorLogger.Printf("%s | %s", timestamp(), l.formatMessage(msg, args...))
	}
}

// Debug logs at debug level.
func (l *Logger) Debug(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.shouldLog(LogLevelDebug) {
		l.debugLogger.Printf("%s | %s", timestamp(), l.formatMessage(msg, args...))
	}
}

// NewFileLogger creates a logger writing to both a file and the console.
func NewFileLogger(filePath string, level string) (*Logger, *os.File, error) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}

	multiWriter := io.MultiWriter(os.Stdout, f)
	logger := NewLogger(multiWriter, level)
	return logger, f, nil
}

func timestamp() string {
	return time.Now().Format(time.RFC3339)
}

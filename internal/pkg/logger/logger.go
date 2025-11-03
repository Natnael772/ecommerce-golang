package logger

import (
	"fmt"
	stdlog "log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	std *stdlog.Logger
}

var logger = New()

func New() *Logger {
	l := stdlog.New(os.Stdout, "", 0)
	return &Logger{std: l}
}

func (l *Logger) log(level, msg string, args ...interface{}) {
	now := time.Now().Format(time.RFC3339)
	file, line := callerInfo()
	message := fmt.Sprintf(msg, args...)
	l.std.Printf("[%s] %-5s %s:%d → %s", now, strings.ToUpper(level), file, line, message)
}

// Info logs general information.
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log("info", msg, args...)
}

// Debug logs debug information.
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log("debug", msg, args...)
}

// Warn logs a warning.
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log("warn", msg, args...)
}

// Error logs an error.
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log("error", msg, args...)
}

// Fatal logs a fatal error and exits.
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log("fatal", msg, args...)
	os.Exit(1)
}

// helper for file + line info
func callerInfo() (string, int) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown", 0
	}
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		file = fmt.Sprintf("%s (%s)", shortFile(file), shortFunc(fn.Name()))
	} else {
		file = shortFile(file)
	}
	return file, line
}

func shortFile(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return path
}

func shortFunc(name string) string {
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}

// Global log helpers for convenience
func Info(msg string, args ...interface{})   { logger.Info(msg, args...) }
func Debug(msg string, args ...interface{})  { logger.Debug(msg, args...) }
func Warn(msg string, args ...interface{})   { logger.Warn(msg, args...) }
func Error(msg string, args ...interface{})  { logger.Error(msg, args...) }
func Fatal(msg string, args ...interface{})  { logger.Fatal(msg, args...) }

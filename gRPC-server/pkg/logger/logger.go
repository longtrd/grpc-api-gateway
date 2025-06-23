package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
)

// Logger represents a structured logger
type Logger struct {
	level  LogLevel
	format LogFormat
	logger *log.Logger
}

// LogLevel represents the logging level
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// LogFormat represents the logging format
type LogFormat int

const (
	TextFormat LogFormat = iota
	JSONFormat
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Time    time.Time              `json:"time"`
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
}

// New creates a new logger instance
func New(levelStr, formatStr string) domain.Logger {
	level := parseLogLevel(levelStr)
	format := parseLogFormat(formatStr)

	return &Logger{
		level:  level,
		format: format,
		logger: log.New(os.Stdout, "", 0),
	}
}

// NewNullLogger creates a logger that discards all output (for testing)
func NewNullLogger() domain.Logger {
	return &NullLogger{}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.log(InfoLevel, msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.log(ErrorLevel, msg, fields...)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.log(DebugLevel, msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.log(WarnLevel, msg, fields...)
}

// log performs the actual logging
func (l *Logger) log(level LogLevel, msg string, fields ...interface{}) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Time:    time.Now().UTC(),
		Level:   l.levelToString(level),
		Message: msg,
		Fields:  l.parseFields(fields...),
	}

	var output string
	switch l.format {
	case JSONFormat:
		output = l.formatJSON(entry)
	default:
		output = l.formatText(entry)
	}

	l.logger.Print(output)
}

// formatJSON formats the log entry as JSON
func (l *Logger) formatJSON(entry LogEntry) string {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Sprintf(`{"time":"%s","level":"error","message":"failed to marshal log entry: %v"}`,
			entry.Time.Format(time.RFC3339), err)
	}
	return string(data)
}

// formatText formats the log entry as plain text
func (l *Logger) formatText(entry LogEntry) string {
	output := fmt.Sprintf("[%s] %s: %s",
		entry.Time.Format("2006-01-02 15:04:05"),
		strings.ToUpper(entry.Level),
		entry.Message)

	if len(entry.Fields) > 0 {
		for key, value := range entry.Fields {
			output += fmt.Sprintf(" %s=%v", key, value)
		}
	}

	return output
}

// parseFields converts variadic interface{} to a map
func (l *Logger) parseFields(fields ...interface{}) map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	result := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			if key, ok := fields[i].(string); ok {
				result[key] = fields[i+1]
			}
		}
	}

	return result
}

// levelToString converts LogLevel to string
func (l *Logger) levelToString(level LogLevel) string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return "info"
	}
}

// parseLogLevel parses string to LogLevel
func parseLogLevel(levelStr string) LogLevel {
	switch strings.ToLower(levelStr) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

// parseLogFormat parses string to LogFormat
func parseLogFormat(formatStr string) LogFormat {
	switch strings.ToLower(formatStr) {
	case "json":
		return JSONFormat
	case "text":
		return TextFormat
	default:
		return TextFormat
	}
}

// NullLogger is a logger that discards all output
type NullLogger struct{}

// NewNullLoggerInstance creates a null logger that implements domain.Logger
func NewNullLoggerInstance() domain.Logger {
	return &NullLogger{}
}

func (n *NullLogger) Info(msg string, fields ...interface{})  {}
func (n *NullLogger) Error(msg string, fields ...interface{}) {}
func (n *NullLogger) Debug(msg string, fields ...interface{}) {}
func (n *NullLogger) Warn(msg string, fields ...interface{})  {}

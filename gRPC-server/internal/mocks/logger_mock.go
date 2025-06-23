package mocks

import "fmt"

// MockLogger is a mock implementation of Logger for testing
type MockLogger struct {
	// Call tracking
	InfoCalls  []LogCall
	ErrorCalls []LogCall
	DebugCalls []LogCall
	WarnCalls  []LogCall
}

type LogCall struct {
	Message string
	Fields  []interface{}
}

// NewMockLogger creates a new mock logger
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// Info logs an info message
func (m *MockLogger) Info(msg string, fields ...interface{}) {
	m.InfoCalls = append(m.InfoCalls, LogCall{Message: msg, Fields: fields})
}

// Error logs an error message
func (m *MockLogger) Error(msg string, fields ...interface{}) {
	m.ErrorCalls = append(m.ErrorCalls, LogCall{Message: msg, Fields: fields})
}

// Debug logs a debug message
func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	m.DebugCalls = append(m.DebugCalls, LogCall{Message: msg, Fields: fields})
}

// Warn logs a warning message
func (m *MockLogger) Warn(msg string, fields ...interface{}) {
	m.WarnCalls = append(m.WarnCalls, LogCall{Message: msg, Fields: fields})
}

// Reset resets the mock state
func (m *MockLogger) Reset() {
	m.InfoCalls = nil
	m.ErrorCalls = nil
	m.DebugCalls = nil
	m.WarnCalls = nil
}

// GetLastInfo returns the last info log call
func (m *MockLogger) GetLastInfo() *LogCall {
	if len(m.InfoCalls) == 0 {
		return nil
	}
	return &m.InfoCalls[len(m.InfoCalls)-1]
}

// GetLastError returns the last error log call
func (m *MockLogger) GetLastError() *LogCall {
	if len(m.ErrorCalls) == 0 {
		return nil
	}
	return &m.ErrorCalls[len(m.ErrorCalls)-1]
}

// String returns a string representation of all log calls (useful for debugging)
func (m *MockLogger) String() string {
	result := "MockLogger calls:\n"
	for _, call := range m.InfoCalls {
		result += fmt.Sprintf("INFO: %s %v\n", call.Message, call.Fields)
	}
	for _, call := range m.ErrorCalls {
		result += fmt.Sprintf("ERROR: %s %v\n", call.Message, call.Fields)
	}
	for _, call := range m.DebugCalls {
		result += fmt.Sprintf("DEBUG: %s %v\n", call.Message, call.Fields)
	}
	for _, call := range m.WarnCalls {
		result += fmt.Sprintf("WARN: %s %v\n", call.Message, call.Fields)
	}
	return result
}

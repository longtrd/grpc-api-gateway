package mocks

import (
	"context"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
)

// MockEventPublisher is a mock implementation of EventPublisher for testing
type MockEventPublisher struct {
	PublishError error

	// Call tracking
	PublishCalls []PublishCall
}

type PublishCall struct {
	Event domain.Event
}

// NewMockEventPublisher creates a new mock event publisher
func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{}
}

// Publish publishes an event
func (m *MockEventPublisher) Publish(ctx context.Context, event domain.Event) error {
	m.PublishCalls = append(m.PublishCalls, PublishCall{Event: event})

	if m.PublishError != nil {
		return m.PublishError
	}

	return nil
}

// Reset resets the mock state
func (m *MockEventPublisher) Reset() {
	m.PublishError = nil
	m.PublishCalls = nil
}

// GetLastEvent returns the last published event
func (m *MockEventPublisher) GetLastEvent() domain.Event {
	if len(m.PublishCalls) == 0 {
		return nil
	}
	return m.PublishCalls[len(m.PublishCalls)-1].Event
}

// GetEventCount returns the number of published events
func (m *MockEventPublisher) GetEventCount() int {
	return len(m.PublishCalls)
}

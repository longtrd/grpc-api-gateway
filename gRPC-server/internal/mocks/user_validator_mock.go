package mocks

import "github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"

// MockUserValidator is a mock implementation of UserValidator for testing
type MockUserValidator struct {
	ValidateError      error
	ValidateEmailError error

	// Call tracking
	ValidateCalls      []ValidateCall
	ValidateEmailCalls []ValidateEmailCall
}

type ValidateCall struct {
	User *domain.User
}

type ValidateEmailCall struct {
	Email string
}

// NewMockUserValidator creates a new mock user validator
func NewMockUserValidator() *MockUserValidator {
	return &MockUserValidator{}
}

// Validate validates a user
func (m *MockUserValidator) Validate(user *domain.User) error {
	m.ValidateCalls = append(m.ValidateCalls, ValidateCall{User: user})

	if m.ValidateError != nil {
		return m.ValidateError
	}

	// Default validation logic
	return user.Validate()
}

// ValidateEmail validates an email address
func (m *MockUserValidator) ValidateEmail(email string) error {
	m.ValidateEmailCalls = append(m.ValidateEmailCalls, ValidateEmailCall{Email: email})

	if m.ValidateEmailError != nil {
		return m.ValidateEmailError
	}

	// Default email validation logic
	user := domain.User{Email: email}
	return user.Validate()
}

// Reset resets the mock state
func (m *MockUserValidator) Reset() {
	m.ValidateError = nil
	m.ValidateEmailError = nil
	m.ValidateCalls = nil
	m.ValidateEmailCalls = nil
}

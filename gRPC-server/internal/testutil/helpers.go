package testutil

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/mocks"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/usecase"
)

// CreateTestUser creates a test user with default values
func CreateTestUser() *domain.User {
	return &domain.User{
		ID:    "test-user-123",
		Name:  "Test User",
		Email: "test.user@example.com",
	}
}

// CreateTestUserWithID creates a test user with a specific ID
func CreateTestUserWithID(id string) *domain.User {
	return &domain.User{
		ID:    id,
		Name:  "Test User " + id,
		Email: "user" + id + "@example.com",
	}
}

// CreateTestUsers creates multiple test users
func CreateTestUsers(count int) []*domain.User {
	users := make([]*domain.User, count)
	for i := 0; i < count; i++ {
		users[i] = CreateTestUserWithID(fmt.Sprintf("user-%d", i+1))
	}
	return users
}

// AssertUserEqual asserts that two users are equal
func AssertUserEqual(t *testing.T, expected, actual *domain.User) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	if expected == nil || actual == nil {
		t.Errorf("Expected user %v, got %v", expected, actual)
		return
	}

	if expected.ID != actual.ID {
		t.Errorf("Expected user ID %s, got %s", expected.ID, actual.ID)
	}

	if expected.Name != actual.Name {
		t.Errorf("Expected user name %s, got %s", expected.Name, actual.Name)
	}

	if expected.Email != actual.Email {
		t.Errorf("Expected user email %s, got %s", expected.Email, actual.Email)
	}
}

// AssertUsersEqual asserts that two user slices are equal
func AssertUsersEqual(t *testing.T, expected, actual []*domain.User) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("Expected %d users, got %d", len(expected), len(actual))
		return
	}

	for i, expectedUser := range expected {
		AssertUserEqual(t, expectedUser, actual[i])
	}
}

// SetupTestUserUseCase creates a UserUseCase with mocked dependencies for testing
func SetupTestUserUseCase() (domain.UserUseCase, *TestMocks) {
	testMocks := &TestMocks{
		Repo:      mocks.NewMockUserRepository(),
		Validator: mocks.NewMockUserValidator(),
		Logger:    mocks.NewMockLogger(),
		Publisher: mocks.NewMockEventPublisher(),
	}

	uc := usecase.NewUserUseCase(testMocks.Repo, testMocks.Validator, testMocks.Logger, testMocks.Publisher)
	return uc, testMocks
}

// TestMocks holds all mock dependencies for testing
type TestMocks struct {
	Repo      *mocks.MockUserRepository
	Validator *mocks.MockUserValidator
	Logger    *mocks.MockLogger
	Publisher *mocks.MockEventPublisher
}

// Reset resets all mocks to their initial state
func (m *TestMocks) Reset() {
	m.Repo.Reset()
	m.Validator.Reset()
	m.Logger.Reset()
	m.Publisher.Reset()
}

// SeedUsers adds test users to the mock repository
func (m *TestMocks) SeedUsers(users ...*domain.User) {
	for _, user := range users {
		m.Repo.Create(context.Background(), user)
	}
}

// AssertErrorIs asserts that an error is of a specific type
func AssertErrorIs(t *testing.T, err, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Errorf("Expected error %v, got %v", target, err)
	}
}

// AssertNoError asserts that there is no error
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// AssertError asserts that there is an error
func AssertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

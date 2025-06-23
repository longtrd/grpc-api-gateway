package mocks

import (
	"context"
	"sync"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User

	// Mock behavior controls
	CreateError error
	GetError    error
	UpdateError error
	DeleteError error
	ListError   error

	// Call tracking
	CreateCalls []CreateCall
	GetCalls    []GetCall
	UpdateCalls []UpdateCall
	DeleteCalls []DeleteCall
	ListCalls   []ListCall
}

type CreateCall struct {
	User *domain.User
}

type GetCall struct {
	ID string
}

type UpdateCall struct {
	User *domain.User
}

type DeleteCall struct {
	ID string
}

type ListCall struct {
	Limit  int
	Offset int
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

// Create creates a new user
func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CreateCalls = append(m.CreateCalls, CreateCall{User: user})

	if m.CreateError != nil {
		return m.CreateError
	}

	cloned := user.Clone()
	m.users[user.ID] = &cloned
	return nil
}

// GetByID retrieves a user by ID
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.GetCalls = append(m.GetCalls, GetCall{ID: id})

	if m.GetError != nil {
		return nil, m.GetError
	}

	user, exists := m.users[id]
	if !exists {
		return nil, nil
	}

	cloned := user.Clone()
	return &cloned, nil
}

// Update updates an existing user
func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.UpdateCalls = append(m.UpdateCalls, UpdateCall{User: user})

	if m.UpdateError != nil {
		return m.UpdateError
	}

	cloned := user.Clone()
	m.users[user.ID] = &cloned
	return nil
}

// Delete deletes a user by ID
func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DeleteCalls = append(m.DeleteCalls, DeleteCall{ID: id})

	if m.DeleteError != nil {
		return m.DeleteError
	}

	delete(m.users, id)
	return nil
}

// List retrieves users with pagination
func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.ListCalls = append(m.ListCalls, ListCall{Limit: limit, Offset: offset})

	if m.ListError != nil {
		return nil, m.ListError
	}

	var users []*domain.User
	count := 0
	for _, user := range m.users {
		if count >= offset && len(users) < limit {
			cloned := user.Clone()
			users = append(users, &cloned)
		}
		count++
	}

	// Return empty slice instead of nil when no users found
	if users == nil {
		users = []*domain.User{}
	}

	return users, nil
}

// Count returns the total number of users
func (m *MockUserRepository) Count(ctx context.Context) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.users), nil
}

// Clear removes all users (useful for testing)
func (m *MockUserRepository) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.users = make(map[string]*domain.User)
	return nil
}

// Reset resets the mock state
func (m *MockUserRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.users = make(map[string]*domain.User)
	m.CreateError = nil
	m.GetError = nil
	m.UpdateError = nil
	m.DeleteError = nil
	m.ListError = nil
	m.CreateCalls = nil
	m.GetCalls = nil
	m.UpdateCalls = nil
	m.DeleteCalls = nil
	m.ListCalls = nil
}

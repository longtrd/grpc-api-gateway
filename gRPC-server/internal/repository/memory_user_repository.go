package repository

import (
	"context"
	"sync"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
)

// MemoryUserRepository is an in-memory implementation of UserRepository
type MemoryUserRepository struct {
	users map[string]*domain.User
	mutex sync.RWMutex
}

// NewMemoryUserRepository creates a new in-memory user repository
func NewMemoryUserRepository() domain.UserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

// Create creates a new user
func (r *MemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if user == nil {
		return domain.ErrInvalidInput
	}

	// Check if user already exists
	if _, exists := r.users[user.ID]; exists {
		return domain.ErrUserAlreadyExists
	}

	// Clone the user to avoid external mutations
	cloned := user.Clone()
	r.users[user.ID] = &cloned

	return nil
}

// GetByID retrieves a user by ID
func (r *MemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if id == "" {
		return nil, domain.ErrInvalidInput
	}

	user, exists := r.users[id]
	if !exists {
		return nil, nil // Return nil for not found (not an error)
	}

	// Clone the user to avoid external mutations
	cloned := user.Clone()
	return &cloned, nil
}

// Update updates an existing user
func (r *MemoryUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if user == nil {
		return domain.ErrInvalidInput
	}

	// Check if user exists
	if _, exists := r.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}

	// Clone the user to avoid external mutations
	cloned := user.Clone()
	r.users[user.ID] = &cloned

	return nil
}

// Delete deletes a user by ID
func (r *MemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if id == "" {
		return domain.ErrInvalidInput
	}

	// Check if user exists
	if _, exists := r.users[id]; !exists {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

// List retrieves users with pagination
func (r *MemoryUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	// Convert map to slice
	allUsers := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		cloned := user.Clone()
		allUsers = append(allUsers, &cloned)
	}

	// Apply pagination
	start := offset
	if start >= len(allUsers) {
		return []*domain.User{}, nil
	}

	end := start + limit
	if end > len(allUsers) {
		end = len(allUsers)
	}

	return allUsers[start:end], nil
}

// Count returns the total number of users
func (r *MemoryUserRepository) Count(ctx context.Context) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.users), nil
}

// Clear removes all users (useful for testing)
func (r *MemoryUserRepository) Clear(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.users = make(map[string]*domain.User)
	return nil
}

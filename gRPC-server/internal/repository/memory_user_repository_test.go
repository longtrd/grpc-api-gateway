package repository

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
)

func TestMemoryUserRepository_Create(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	tests := []struct {
		name        string
		user        *domain.User
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful creation",
			user: &domain.User{
				ID:    "123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name:        "nil user",
			user:        nil,
			wantErr:     true,
			expectedErr: domain.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.expectedErr != nil && err != tt.expectedErr {
				t.Errorf("Create() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestMemoryUserRepository_Create_DuplicateUser(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := &domain.User{
		ID:    "123",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Create user first time - should succeed
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("First Create() failed: %v", err)
	}

	// Create same user again - should fail
	err = repo.Create(ctx, user)
	if err != domain.ErrUserAlreadyExists {
		t.Errorf("Second Create() error = %v, expected %v", err, domain.ErrUserAlreadyExists)
	}
}

func TestMemoryUserRepository_GetByID(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	// Create test user
	user := &domain.User{
		ID:    "123",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	repo.Create(ctx, user)

	tests := []struct {
		name        string
		id          string
		expectUser  bool
		wantErr     bool
		expectedErr error
	}{
		{
			name:       "existing user",
			id:         "123",
			expectUser: true,
			wantErr:    false,
		},
		{
			name:       "non-existent user",
			id:         "456",
			expectUser: false,
			wantErr:    false,
		},
		{
			name:        "empty ID",
			id:          "",
			expectUser:  false,
			wantErr:     true,
			expectedErr: domain.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByID(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.expectedErr != nil && err != tt.expectedErr {
				t.Errorf("GetByID() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if tt.expectUser && result == nil {
				t.Error("GetByID() returned nil user when user was expected")
			}

			if !tt.expectUser && result != nil {
				t.Errorf("GetByID() returned user %v when nil was expected", result)
			}

			if tt.expectUser && result != nil {
				if result.ID != user.ID || result.Name != user.Name || result.Email != user.Email {
					t.Errorf("GetByID() returned %v, expected %v", result, user)
				}
			}
		})
	}
}

func TestMemoryUserRepository_Update(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	// Create test user
	originalUser := &domain.User{
		ID:    "123",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	repo.Create(ctx, originalUser)

	tests := []struct {
		name        string
		user        *domain.User
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful update",
			user: &domain.User{
				ID:    "123",
				Name:  "John Updated",
				Email: "john.updated@example.com",
			},
			wantErr: false,
		},
		{
			name:        "nil user",
			user:        nil,
			wantErr:     true,
			expectedErr: domain.ErrInvalidInput,
		},
		{
			name: "non-existent user",
			user: &domain.User{
				ID:    "456",
				Name:  "Jane Doe",
				Email: "jane@example.com",
			},
			wantErr:     true,
			expectedErr: domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.expectedErr != nil && err != tt.expectedErr {
				t.Errorf("Update() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestMemoryUserRepository_Delete(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	// Create test user
	user := &domain.User{
		ID:    "123",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	repo.Create(ctx, user)

	tests := []struct {
		name        string
		id          string
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "successful deletion",
			id:      "123",
			wantErr: false,
		},
		{
			name:        "empty ID",
			id:          "",
			wantErr:     true,
			expectedErr: domain.ErrInvalidInput,
		},
		{
			name:        "non-existent user",
			id:          "456",
			wantErr:     true,
			expectedErr: domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.expectedErr != nil && err != tt.expectedErr {
				t.Errorf("Delete() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestMemoryUserRepository_List(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	// Create test users
	users := []*domain.User{
		{ID: "1", Name: "User 1", Email: "user1@example.com"},
		{ID: "2", Name: "User 2", Email: "user2@example.com"},
		{ID: "3", Name: "User 3", Email: "user3@example.com"},
	}

	for _, user := range users {
		repo.Create(ctx, user)
	}

	tests := []struct {
		name        string
		limit       int
		offset      int
		expectedLen int
		wantErr     bool
	}{
		{
			name:        "get all users",
			limit:       10,
			offset:      0,
			expectedLen: 3,
			wantErr:     false,
		},
		{
			name:        "pagination",
			limit:       2,
			offset:      1,
			expectedLen: 2,
			wantErr:     false,
		},
		{
			name:        "offset beyond range",
			limit:       10,
			offset:      10,
			expectedLen: 0,
			wantErr:     false,
		},
		{
			name:        "negative offset",
			limit:       10,
			offset:      -1,
			expectedLen: 3,
			wantErr:     false,
		},
		{
			name:        "zero limit uses default",
			limit:       0,
			offset:      0,
			expectedLen: 3,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.List(ctx, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(result) != tt.expectedLen {
				t.Errorf("List() returned %d users, expected %d", len(result), tt.expectedLen)
			}
		})
	}
}

func TestMemoryUserRepository_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	const numGoroutines = 100
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Test concurrent creates
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			user := &domain.User{
				ID:    fmt.Sprintf("user-%d", id),
				Name:  fmt.Sprintf("User %d", id),
				Email: fmt.Sprintf("user%d@example.com", id),
			}
			if err := repo.Create(ctx, user); err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent Create() failed: %v", err)
	}

	// Verify all users were created
	users, err := repo.List(ctx, 200, 0)
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(users) != numGoroutines {
		t.Errorf("Expected %d users, got %d", numGoroutines, len(users))
	}
}

func TestMemoryUserRepository_Count(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	// Initially empty
	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count() failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	// Add users
	users := []*domain.User{
		{ID: "1", Name: "User 1", Email: "user1@example.com"},
		{ID: "2", Name: "User 2", Email: "user2@example.com"},
	}

	for _, user := range users {
		repo.Create(ctx, user)
	}

	count, err = repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count() failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestMemoryUserRepository_Clear(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	// Add a user
	user := &domain.User{
		ID:    "123",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	repo.Create(ctx, user)

	// Clear repository
	err := repo.Clear(ctx)
	if err != nil {
		t.Fatalf("Clear() failed: %v", err)
	}

	// Verify empty
	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count() failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0 after clear, got %d", count)
	}
}

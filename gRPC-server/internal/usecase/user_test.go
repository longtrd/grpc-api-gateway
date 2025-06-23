package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/mocks"
)

func TestUserUseCase_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		user           *domain.User
		repoError      error
		validatorError error
		existingUser   *domain.User
		wantErr        bool
		expectedErr    error
	}{
		{
			name: "successful user creation",
			user: &domain.User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john.doe@example.com",
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
			name: "validation error",
			user: &domain.User{
				ID:    "",
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			validatorError: errors.New("user ID cannot be empty"),
			wantErr:        true,
		},
		{
			name: "user already exists",
			user: &domain.User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			existingUser: &domain.User{
				ID:    "user-123",
				Name:  "Existing User",
				Email: "existing@example.com",
			},
			wantErr:     true,
			expectedErr: domain.ErrUserAlreadyExists,
		},
		{
			name: "repository error",
			user: &domain.User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			repoError: errors.New("database error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := mocks.NewMockUserRepository()
			mockValidator := mocks.NewMockUserValidator()
			mockLogger := mocks.NewMockLogger()
			mockPublisher := mocks.NewMockEventPublisher()

			// Configure mock behaviors
			if tt.repoError != nil {
				mockRepo.CreateError = tt.repoError
			}
			if tt.validatorError != nil {
				mockValidator.ValidateError = tt.validatorError
			}
			if tt.existingUser != nil {
				mockRepo.Create(context.Background(), tt.existingUser)
			}

			// Create use case
			uc := NewUserUseCase(mockRepo, mockValidator, mockLogger, mockPublisher)

			// Execute
			err := uc.CreateUser(context.Background(), tt.user)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Errorf("CreateUser() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			// Verify calls
			if tt.user != nil && tt.validatorError == nil {
				if len(mockValidator.ValidateCalls) != 1 {
					t.Errorf("Expected 1 validator call, got %d", len(mockValidator.ValidateCalls))
				}
			}

			if tt.user != nil && tt.validatorError == nil && tt.existingUser == nil && tt.repoError == nil {
				if len(mockRepo.CreateCalls) != 1 {
					t.Errorf("Expected 1 repository create call, got %d", len(mockRepo.CreateCalls))
				}
			}
		})
	}
}

func TestUserUseCase_GetUser(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		repoUser    *domain.User
		repoError   error
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "successful user retrieval",
			userID: "user-123",
			repoUser: &domain.User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name:        "empty user ID",
			userID:      "",
			wantErr:     true,
			expectedErr: domain.ErrInvalidInput,
		},
		{
			name:        "user not found",
			userID:      "user-123",
			repoUser:    nil,
			wantErr:     true,
			expectedErr: domain.ErrUserNotFound,
		},
		{
			name:      "repository error",
			userID:    "user-123",
			repoError: errors.New("database error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := mocks.NewMockUserRepository()
			mockValidator := mocks.NewMockUserValidator()
			mockLogger := mocks.NewMockLogger()
			mockPublisher := mocks.NewMockEventPublisher()

			// Configure mock behaviors
			if tt.repoError != nil {
				mockRepo.GetError = tt.repoError
			} else if tt.repoUser != nil {
				mockRepo.Create(context.Background(), tt.repoUser)
			}

			// Create use case
			uc := NewUserUseCase(mockRepo, mockValidator, mockLogger, mockPublisher)

			// Execute
			user, err := uc.GetUser(context.Background(), tt.userID)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Errorf("GetUser() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			if !tt.wantErr && user == nil {
				t.Error("GetUser() returned nil user when success expected")
			}

			if !tt.wantErr && user != nil && user.ID != tt.userID {
				t.Errorf("GetUser() returned user with ID %v, expected %v", user.ID, tt.userID)
			}
		})
	}
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		user           *domain.User
		existingUser   *domain.User
		repoError      error
		validatorError error
		wantErr        bool
		expectedErr    error
	}{
		{
			name: "successful user update",
			user: &domain.User{
				ID:    "user-123",
				Name:  "John Updated",
				Email: "john.updated@example.com",
			},
			existingUser: &domain.User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john.doe@example.com",
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
			name: "validation error",
			user: &domain.User{
				ID:    "user-123",
				Name:  "",
				Email: "john.doe@example.com",
			},
			validatorError: errors.New("user name cannot be empty"),
			wantErr:        true,
		},
		{
			name: "user not found",
			user: &domain.User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			existingUser: nil,
			wantErr:      true,
			expectedErr:  domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := mocks.NewMockUserRepository()
			mockValidator := mocks.NewMockUserValidator()
			mockLogger := mocks.NewMockLogger()
			mockPublisher := mocks.NewMockEventPublisher()

			// Configure mock behaviors
			if tt.repoError != nil {
				mockRepo.UpdateError = tt.repoError
			}
			if tt.validatorError != nil {
				mockValidator.ValidateError = tt.validatorError
			}
			if tt.existingUser != nil {
				mockRepo.Create(context.Background(), tt.existingUser)
			}

			// Create use case
			uc := NewUserUseCase(mockRepo, mockValidator, mockLogger, mockPublisher)

			// Execute
			err := uc.UpdateUser(context.Background(), tt.user)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Errorf("UpdateUser() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		existingUser *domain.User
		repoError    error
		wantErr      bool
		expectedErr  error
	}{
		{
			name:   "successful user deletion",
			userID: "user-123",
			existingUser: &domain.User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name:        "empty user ID",
			userID:      "",
			wantErr:     true,
			expectedErr: domain.ErrInvalidInput,
		},
		{
			name:         "user not found",
			userID:       "user-123",
			existingUser: nil,
			wantErr:      true,
			expectedErr:  domain.ErrUserNotFound,
		},
		{
			name:   "repository error",
			userID: "user-123",
			existingUser: &domain.User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			repoError: errors.New("database error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := mocks.NewMockUserRepository()
			mockValidator := mocks.NewMockUserValidator()
			mockLogger := mocks.NewMockLogger()
			mockPublisher := mocks.NewMockEventPublisher()

			// Configure mock behaviors
			if tt.repoError != nil {
				mockRepo.DeleteError = tt.repoError
			}
			if tt.existingUser != nil {
				mockRepo.Create(context.Background(), tt.existingUser)
			}

			// Create use case
			uc := NewUserUseCase(mockRepo, mockValidator, mockLogger, mockPublisher)

			// Execute
			err := uc.DeleteUser(context.Background(), tt.userID)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Errorf("DeleteUser() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestUserUseCase_ListUsers(t *testing.T) {
	tests := []struct {
		name      string
		limit     int
		offset    int
		repoUsers []*domain.User
		repoError error
		wantErr   bool
	}{
		{
			name:   "successful user listing",
			limit:  10,
			offset: 0,
			repoUsers: []*domain.User{
				{ID: "user-1", Name: "User 1", Email: "user1@example.com"},
				{ID: "user-2", Name: "User 2", Email: "user2@example.com"},
			},
			wantErr: false,
		},
		{
			name:    "zero limit sets default",
			limit:   0,
			offset:  0,
			wantErr: false,
		},
		{
			name:    "negative offset sets to zero",
			limit:   10,
			offset:  -5,
			wantErr: false,
		},
		{
			name:    "limit too high gets capped",
			limit:   150,
			offset:  0,
			wantErr: false,
		},
		{
			name:      "repository error",
			limit:     10,
			offset:    0,
			repoError: errors.New("database error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := mocks.NewMockUserRepository()
			mockValidator := mocks.NewMockUserValidator()
			mockLogger := mocks.NewMockLogger()
			mockPublisher := mocks.NewMockEventPublisher()

			// Configure mock behaviors
			if tt.repoError != nil {
				mockRepo.ListError = tt.repoError
			}
			for _, user := range tt.repoUsers {
				mockRepo.Create(context.Background(), user)
			}

			// Create use case
			uc := NewUserUseCase(mockRepo, mockValidator, mockLogger, mockPublisher)

			// Execute
			users, err := uc.ListUsers(context.Background(), tt.limit, tt.offset)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && users == nil {
				t.Error("ListUsers() returned nil users when success expected")
			}

			// Check that limits are applied correctly
			if !tt.wantErr && len(mockRepo.ListCalls) > 0 {
				call := mockRepo.ListCalls[0]
				if tt.limit <= 0 && call.Limit != 10 {
					t.Errorf("Expected default limit 10, got %d", call.Limit)
				}
				if tt.limit > 100 && call.Limit != 100 {
					t.Errorf("Expected capped limit 100, got %d", call.Limit)
				}
				if tt.offset < 0 && call.Offset != 0 {
					t.Errorf("Expected offset 0, got %d", call.Offset)
				}
			}
		})
	}
}

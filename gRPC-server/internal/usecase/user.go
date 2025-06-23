package usecase

import (
	"context"
	"fmt"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
)

// userUseCase implements the UserUseCase interface
type userUseCase struct {
	repo      domain.UserRepository
	validator domain.UserValidator
	logger    domain.Logger
	publisher domain.EventPublisher
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(
	repo domain.UserRepository,
	validator domain.UserValidator,
	logger domain.Logger,
	publisher domain.EventPublisher,
) domain.UserUseCase {
	return &userUseCase{
		repo:      repo,
		validator: validator,
		logger:    logger,
		publisher: publisher,
	}
}

// CreateUser creates a new user
func (uc *userUseCase) CreateUser(ctx context.Context, user *domain.User) error {
	if user == nil {
		return domain.ErrInvalidInput
	}

	uc.logger.Info("Creating user", "user_id", user.ID)

	// Validate user
	if err := uc.validator.Validate(user); err != nil {
		uc.logger.Error("User validation failed", "error", err, "user_id", user.ID)
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if user already exists
	existingUser, err := uc.repo.GetByID(ctx, user.ID)
	if err == nil && existingUser != nil {
		uc.logger.Warn("User already exists", "user_id", user.ID)
		return domain.ErrUserAlreadyExists
	}

	// Create user
	if err := uc.repo.Create(ctx, user); err != nil {
		uc.logger.Error("Failed to create user", "error", err, "user_id", user.ID)
		return fmt.Errorf("failed to create user: %w", err)
	}

	uc.logger.Info("User created successfully", "user_id", user.ID)
	return nil
}

// GetUser retrieves a user by ID
func (uc *userUseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}

	uc.logger.Debug("Getting user", "user_id", id)

	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get user", "error", err, "user_id", id)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		uc.logger.Warn("User not found", "user_id", id)
		return nil, domain.ErrUserNotFound
	}

	uc.logger.Debug("User retrieved successfully", "user_id", id)
	return user, nil
}

// UpdateUser updates an existing user
func (uc *userUseCase) UpdateUser(ctx context.Context, user *domain.User) error {
	if user == nil {
		return domain.ErrInvalidInput
	}

	uc.logger.Info("Updating user", "user_id", user.ID)

	// Validate user
	if err := uc.validator.Validate(user); err != nil {
		uc.logger.Error("User validation failed", "error", err, "user_id", user.ID)
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if user exists
	existingUser, err := uc.repo.GetByID(ctx, user.ID)
	if err != nil {
		uc.logger.Error("Failed to get user for update", "error", err, "user_id", user.ID)
		return fmt.Errorf("failed to get user: %w", err)
	}

	if existingUser == nil {
		uc.logger.Warn("User not found for update", "user_id", user.ID)
		return domain.ErrUserNotFound
	}

	// Update user
	if err := uc.repo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update user", "error", err, "user_id", user.ID)
		return fmt.Errorf("failed to update user: %w", err)
	}

	uc.logger.Info("User updated successfully", "user_id", user.ID)
	return nil
}

// DeleteUser deletes a user by ID
func (uc *userUseCase) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidInput
	}

	uc.logger.Info("Deleting user", "user_id", id)

	// Check if user exists
	existingUser, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get user for deletion", "error", err, "user_id", id)
		return fmt.Errorf("failed to get user: %w", err)
	}

	if existingUser == nil {
		uc.logger.Warn("User not found for deletion", "user_id", id)
		return domain.ErrUserNotFound
	}

	// Delete user
	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete user", "error", err, "user_id", id)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	uc.logger.Info("User deleted successfully", "user_id", id)
	return nil
}

// ListUsers retrieves a list of users with pagination
func (uc *userUseCase) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}
	if limit > 100 {
		limit = 100 // max limit
	}
	if offset < 0 {
		offset = 0
	}

	uc.logger.Debug("Listing users", "limit", limit, "offset", offset)

	users, err := uc.repo.List(ctx, limit, offset)
	if err != nil {
		uc.logger.Error("Failed to list users", "error", err, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	uc.logger.Debug("Users retrieved successfully", "count", len(users))
	return users, nil
}

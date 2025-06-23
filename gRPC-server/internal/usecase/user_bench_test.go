package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/mocks"
)

func BenchmarkUserUseCase_CreateUser(b *testing.B) {
	mockRepo := mocks.NewMockUserRepository()
	mockValidator := mocks.NewMockUserValidator()
	mockLogger := mocks.NewMockLogger()
	mockPublisher := mocks.NewMockEventPublisher()

	uc := NewUserUseCase(mockRepo, mockValidator, mockLogger, mockPublisher)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &domain.User{
			ID:    fmt.Sprintf("bench-user-%d", i),
			Name:  fmt.Sprintf("Bench User %d", i),
			Email: fmt.Sprintf("bench%d@example.com", i),
		}
		uc.CreateUser(ctx, user)
	}
}

func BenchmarkUserUseCase_GetUser(b *testing.B) {
	mockRepo := mocks.NewMockUserRepository()
	mockValidator := mocks.NewMockUserValidator()
	mockLogger := mocks.NewMockLogger()
	mockPublisher := mocks.NewMockEventPublisher()

	uc := NewUserUseCase(mockRepo, mockValidator, mockLogger, mockPublisher)
	ctx := context.Background()

	// Pre-populate with users
	const numUsers = 1000
	for i := 0; i < numUsers; i++ {
		user := &domain.User{
			ID:    fmt.Sprintf("bench-user-%d", i),
			Name:  fmt.Sprintf("Bench User %d", i),
			Email: fmt.Sprintf("bench%d@example.com", i),
		}
		uc.CreateUser(ctx, user)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("bench-user-%d", i%numUsers)
		uc.GetUser(ctx, userID)
	}
}

func BenchmarkUserUseCase_UpdateUser(b *testing.B) {
	mockRepo := mocks.NewMockUserRepository()
	mockValidator := mocks.NewMockUserValidator()
	mockLogger := mocks.NewMockLogger()
	mockPublisher := mocks.NewMockEventPublisher()

	uc := NewUserUseCase(mockRepo, mockValidator, mockLogger, mockPublisher)
	ctx := context.Background()

	// Pre-populate with users
	const numUsers = 1000
	for i := 0; i < numUsers; i++ {
		user := &domain.User{
			ID:    fmt.Sprintf("bench-user-%d", i),
			Name:  fmt.Sprintf("Bench User %d", i),
			Email: fmt.Sprintf("bench%d@example.com", i),
		}
		uc.CreateUser(ctx, user)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &domain.User{
			ID:    fmt.Sprintf("bench-user-%d", i%numUsers),
			Name:  fmt.Sprintf("Updated Bench User %d", i),
			Email: fmt.Sprintf("updated-bench%d@example.com", i),
		}
		uc.UpdateUser(ctx, user)
	}
}

func BenchmarkUserUseCase_ListUsers(b *testing.B) {
	mockRepo := mocks.NewMockUserRepository()
	mockValidator := mocks.NewMockUserValidator()
	mockLogger := mocks.NewMockLogger()
	mockPublisher := mocks.NewMockEventPublisher()

	uc := NewUserUseCase(mockRepo, mockValidator, mockLogger, mockPublisher)
	ctx := context.Background()

	// Pre-populate with users
	const numUsers = 1000
	for i := 0; i < numUsers; i++ {
		user := &domain.User{
			ID:    fmt.Sprintf("bench-user-%d", i),
			Name:  fmt.Sprintf("Bench User %d", i),
			Email: fmt.Sprintf("bench%d@example.com", i),
		}
		uc.CreateUser(ctx, user)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uc.ListUsers(ctx, 50, 0)
	}
}

func BenchmarkUserValidation(b *testing.B) {
	user := &domain.User{
		ID:    "bench-user-validation",
		Name:  "Bench User",
		Email: "bench@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.Validate()
	}
}

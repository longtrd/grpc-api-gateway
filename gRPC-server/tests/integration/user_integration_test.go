//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/mocks"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/testutil"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/usecase"
)

func TestUserIntegration_FullWorkflow(t *testing.T) {
	// Setup - this would typically connect to a real database
	// For this example, we're using mocks but with more realistic scenarios
	mockRepo := mocks.NewMockUserRepository()
	mockValidator := mocks.NewMockUserValidator()
	mockLogger := mocks.NewMockLogger()
	mockPublisher := mocks.NewMockEventPublisher()

	userUC := usecase.NewUserUseCase(mockRepo, mockValidator, mockLogger, mockPublisher)
	ctx := context.Background()

	t.Run("complete user lifecycle", func(t *testing.T) {
		// Create a user
		user := testutil.CreateTestUser()
		err := userUC.CreateUser(ctx, user)
		testutil.AssertNoError(t, err)

		// Retrieve the user
		retrievedUser, err := userUC.GetUser(ctx, user.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertUserEqual(t, user, retrievedUser)

		// Update the user
		user.Name = "Updated Name"
		user.Email = "updated@example.com"
		err = userUC.UpdateUser(ctx, user)
		testutil.AssertNoError(t, err)

		// Verify the update
		updatedUser, err := userUC.GetUser(ctx, user.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertUserEqual(t, user, updatedUser)

		// Delete the user
		err = userUC.DeleteUser(ctx, user.ID)
		testutil.AssertNoError(t, err)

		// Verify deletion
		_, err = userUC.GetUser(ctx, user.ID)
		testutil.AssertErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("bulk operations", func(t *testing.T) {
		// Create multiple users
		users := testutil.CreateTestUsers(5)
		for _, user := range users {
			err := userUC.CreateUser(ctx, user)
			testutil.AssertNoError(t, err)
		}

		// List users
		retrievedUsers, err := userUC.ListUsers(ctx, 10, 0)
		testutil.AssertNoError(t, err)

		if len(retrievedUsers) < 5 {
			t.Errorf("Expected at least 5 users, got %d", len(retrievedUsers))
		}

		// Test pagination
		firstPage, err := userUC.ListUsers(ctx, 2, 0)
		testutil.AssertNoError(t, err)

		secondPage, err := userUC.ListUsers(ctx, 2, 2)
		testutil.AssertNoError(t, err)

		if len(firstPage) != 2 || len(secondPage) != 2 {
			t.Errorf("Pagination failed: first page %d, second page %d", len(firstPage), len(secondPage))
		}

		// Cleanup
		for _, user := range users {
			userUC.DeleteUser(ctx, user.ID)
		}
	})

	t.Run("concurrent operations", func(t *testing.T) {
		// Test concurrent user creation
		const numGoroutines = 10
		done := make(chan bool, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				user := testutil.CreateTestUserWithID(fmt.Sprintf("concurrent-user-%d", id))
				err := userUC.CreateUser(ctx, user)
				if err != nil {
					errors <- err
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Check for errors
		close(errors)
		for err := range errors {
			t.Errorf("Concurrent operation failed: %v", err)
		}

		// Verify all users were created
		users, err := userUC.ListUsers(ctx, 20, 0)
		testutil.AssertNoError(t, err)

		concurrentUsers := 0
		for _, user := range users {
			if strings.HasPrefix(user.ID, "concurrent-user-") {
				concurrentUsers++
			}
		}

		if concurrentUsers != numGoroutines {
			t.Errorf("Expected %d concurrent users, got %d", numGoroutines, concurrentUsers)
		}
	})
}

func TestUserIntegration_ErrorScenarios(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	mockValidator := mocks.NewMockUserValidator()
	mockLogger := mocks.NewMockLogger()
	mockPublisher := mocks.NewMockEventPublisher()

	userUC := usecase.NewUserUseCase(mockRepo, mockValidator, mockLogger, mockPublisher)
	ctx := context.Background()

	t.Run("validation errors", func(t *testing.T) {
		invalidUsers := []*domain.User{
			{ID: "", Name: "Test", Email: "test@example.com"},  // Empty ID
			{ID: "test", Name: "", Email: "test@example.com"},  // Empty name
			{ID: "test", Name: "Test", Email: ""},              // Empty email
			{ID: "test", Name: "Test", Email: "invalid-email"}, // Invalid email
		}

		for _, user := range invalidUsers {
			err := userUC.CreateUser(ctx, user)
			testutil.AssertError(t, err)
		}
	})

	t.Run("duplicate user creation", func(t *testing.T) {
		user := testutil.CreateTestUser()

		// First creation should succeed
		err := userUC.CreateUser(ctx, user)
		testutil.AssertNoError(t, err)

		// Second creation should fail
		err = userUC.CreateUser(ctx, user)
		testutil.AssertErrorIs(t, err, domain.ErrUserAlreadyExists)
	})

	t.Run("operations on non-existent user", func(t *testing.T) {
		nonExistentID := "non-existent-user"

		// Get non-existent user
		_, err := userUC.GetUser(ctx, nonExistentID)
		testutil.AssertErrorIs(t, err, domain.ErrUserNotFound)

		// Update non-existent user
		user := testutil.CreateTestUserWithID(nonExistentID)
		err = userUC.UpdateUser(ctx, user)
		testutil.AssertErrorIs(t, err, domain.ErrUserNotFound)

		// Delete non-existent user
		err = userUC.DeleteUser(ctx, nonExistentID)
		testutil.AssertErrorIs(t, err, domain.ErrUserNotFound)
	})
}

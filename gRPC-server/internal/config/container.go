package config

import (
	"fmt"
	"sync"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/mocks"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/usecase"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/pkg/logger"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/pkg/validator"
)

// Container manages all application dependencies
type Container struct {
	// Configuration
	Config *Config
	Logger domain.Logger

	// Repositories
	UserRepository domain.UserRepository

	// Use Cases
	UserUseCase domain.UserUseCase

	// Handlers
	// UserHandler will be added when gRPC handler is implemented

	// External services
	// DatabaseClient will be added when database integration is needed

	// Validators
	UserValidator domain.UserValidator

	// Event Publisher
	EventPublisher domain.EventPublisher

	// Sync for lazy loading
	mu sync.RWMutex
}

// NewContainer creates a new dependency injection container
func NewContainer(cfg *Config) (*Container, error) {
	container := &Container{
		Config: cfg,
	}

	if err := container.buildDependencies(); err != nil {
		return nil, fmt.Errorf("failed to build dependencies: %w", err)
	}

	return container, nil
}

// buildDependencies builds all application dependencies
func (c *Container) buildDependencies() error {
	// Build logger first
	if err := c.buildLogger(); err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}

	// Build validators
	if err := c.buildValidators(); err != nil {
		return fmt.Errorf("failed to build validators: %w", err)
	}

	// Build event publisher
	if err := c.buildEventPublisher(); err != nil {
		return fmt.Errorf("failed to build event publisher: %w", err)
	}

	// Build repositories
	if err := c.buildRepositories(); err != nil {
		return fmt.Errorf("failed to build repositories: %w", err)
	}

	// Build use cases
	if err := c.buildUseCases(); err != nil {
		return fmt.Errorf("failed to build use cases: %w", err)
	}

	// Build handlers would go here when implemented

	return nil
}

// buildLogger creates the logger based on configuration
func (c *Container) buildLogger() error {
	c.Logger = logger.New(c.Config.Logger.Level, c.Config.Logger.Format)
	return nil
}

// buildValidators creates validators for domain entities
func (c *Container) buildValidators() error {
	c.UserValidator = validator.NewUserValidator()
	return nil
}

// buildEventPublisher creates the event publisher
func (c *Container) buildEventPublisher() error {
	// For now, use a null publisher. Replace with real implementation as needed.
	c.EventPublisher = mocks.NewMockEventPublisher()
	return nil
}

// buildRepositories creates repository implementations
func (c *Container) buildRepositories() error {
	switch c.Config.Database.Type {
	case "memory":
		c.UserRepository = mocks.NewMockUserRepository()
	case "postgres":
		// TODO: Implement PostgreSQL repository
		return fmt.Errorf("postgres repository not implemented yet")
	default:
		return fmt.Errorf("unsupported database type: %s", c.Config.Database.Type)
	}
	return nil
}

// buildUseCases creates use case implementations
func (c *Container) buildUseCases() error {
	c.UserUseCase = usecase.NewUserUseCase(
		c.UserRepository,
		c.UserValidator,
		c.Logger,
		c.EventPublisher,
	)
	return nil
}

// Cleanup cleans up resources when shutting down
func (c *Container) Cleanup() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errors []error

	// Close database connections if any
	// if closer, ok := c.DatabaseClient.(io.Closer); ok {
	//     if err := closer.Close(); err != nil {
	//         errors = append(errors, err)
	//     }
	// }

	// Close event publisher if needed
	// if closer, ok := c.EventPublisher.(io.Closer); ok {
	//     if err := closer.Close(); err != nil {
	//         errors = append(errors, err)
	//     }
	// }

	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %v", errors)
	}

	return nil
}

// GetLogger returns the logger instance
func (c *Container) GetLogger() domain.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Logger
}

// GetUserUseCase returns the user use case instance
func (c *Container) GetUserUseCase() domain.UserUseCase {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.UserUseCase
}

// GetConfig returns the configuration
func (c *Container) GetConfig() *Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Config
}

// NewTestContainer creates a container for testing
func NewTestContainer() *Container {
	cfg := &Config{
		Logger: LoggerConfig{
			Level:  "debug",
			Format: "text",
		},
		Database: DatabaseConfig{
			Type: "memory",
		},
		Features: FeatureConfig{
			EnableHealthCheck: true,
		},
	}

	return &Container{
		Config:         cfg,
		Logger:         logger.NewNullLoggerInstance(), // Silent logger for tests
		UserRepository: mocks.NewMockUserRepository(),
		UserValidator:  mocks.NewMockUserValidator(),
		EventPublisher: mocks.NewMockEventPublisher(),
	}
}

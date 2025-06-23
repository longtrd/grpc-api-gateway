package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// User represents a user entity in the domain
type User struct {
	ID    string
	Name  string
	Email string
}

// Domain-specific errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")

	// Validation errors
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrInvalidUserName  = errors.New("invalid user name")
	ErrInvalidUserEmail = errors.New("invalid user email")
)

// emailRegex is a simple regex for email validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Validate validates the user entity
func (u User) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return errors.New("user ID cannot be empty")
	}

	if strings.TrimSpace(u.Name) == "" {
		return errors.New("user name cannot be empty")
	}

	if strings.TrimSpace(u.Email) == "" {
		return errors.New("user email cannot be empty")
	}

	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	return nil
}

// String returns a string representation of the user
func (u User) String() string {
	return fmt.Sprintf("User{ID: %s, Name: %s, Email: %s}", u.ID, u.Name, u.Email)
}

// IsEmpty checks if the user is empty
func (u User) IsEmpty() bool {
	return u.ID == "" && u.Name == "" && u.Email == ""
}

// Clone creates a copy of the user
func (u User) Clone() User {
	return User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

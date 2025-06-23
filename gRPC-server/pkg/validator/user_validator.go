package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
)

// UserValidator implements domain.UserValidator interface
type UserValidator struct {
	emailRegex *regexp.Regexp
}

// NewUserValidator creates a new user validator instance
func NewUserValidator() domain.UserValidator {
	// RFC 5322 compliant email regex (simplified)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return &UserValidator{
		emailRegex: emailRegex,
	}
}

// Validate validates a user entity according to business rules
func (v *UserValidator) Validate(user *domain.User) error {
	if err := v.validateID(user.ID); err != nil {
		return err
	}

	if err := v.validateName(user.Name); err != nil {
		return err
	}

	if err := v.validateEmail(user.Email); err != nil {
		return err
	}

	return nil
}

// validateID validates the user ID
func (v *UserValidator) validateID(id string) error {
	if id == "" {
		return domain.ErrInvalidUserID
	}

	// Trim whitespace
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.ErrInvalidUserID
	}

	// Check minimum length
	if len(id) < 1 {
		return domain.ErrInvalidUserID
	}

	// Check maximum length (reasonable limit)
	if len(id) > 255 {
		return fmt.Errorf("%w: id too long (max 255 characters)", domain.ErrInvalidUserID)
	}

	return nil
}

// validateName validates the user name
func (v *UserValidator) validateName(name string) error {
	if name == "" {
		return domain.ErrInvalidUserName
	}

	// Trim whitespace
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.ErrInvalidUserName
	}

	// Check minimum length
	if len(name) < 2 {
		return fmt.Errorf("%w: name too short (min 2 characters)", domain.ErrInvalidUserName)
	}

	// Check maximum length
	if len(name) > 100 {
		return fmt.Errorf("%w: name too long (max 100 characters)", domain.ErrInvalidUserName)
	}

	// Check for invalid characters (basic check)
	if strings.ContainsAny(name, "<>\"'&") {
		return fmt.Errorf("%w: name contains invalid characters", domain.ErrInvalidUserName)
	}

	return nil
}

// validateEmail validates the user email
func (v *UserValidator) validateEmail(email string) error {
	if email == "" {
		return domain.ErrInvalidUserEmail
	}

	// Trim whitespace
	email = strings.TrimSpace(email)
	if email == "" {
		return domain.ErrInvalidUserEmail
	}

	// Check maximum length
	if len(email) > 254 {
		return fmt.Errorf("%w: email too long (max 254 characters)", domain.ErrInvalidUserEmail)
	}

	// Validate email format using regex
	if !v.emailRegex.MatchString(email) {
		return domain.ErrInvalidUserEmail
	}

	return nil
}

// ValidateForUpdate validates a user for update operations
func (v *UserValidator) ValidateForUpdate(user *domain.User) error {
	// For updates, we might have different validation rules
	// For now, use the same validation as create
	return v.Validate(user)
}

// ValidateID validates just the user ID (useful for operations that only need ID)
func (v *UserValidator) ValidateID(id string) error {
	return v.validateID(id)
}

// ValidateEmail validates just the email (useful for email-specific operations)
func (v *UserValidator) ValidateEmail(email string) error {
	return v.validateEmail(email)
}

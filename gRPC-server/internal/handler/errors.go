package handler

import (
	"errors"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mapDomainError maps domain errors to gRPC status codes
func mapDomainError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Errorf(codes.NotFound, "user not found: %v", err)
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Errorf(codes.InvalidArgument, "invalid input: %v", err)
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return status.Errorf(codes.AlreadyExists, "user already exists: %v", err)
	case errors.Is(err, domain.ErrUnauthorized):
		return status.Errorf(codes.Unauthenticated, "unauthorized: %v", err)
	case errors.Is(err, domain.ErrInvalidUserID):
		return status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	case errors.Is(err, domain.ErrInvalidUserName):
		return status.Errorf(codes.InvalidArgument, "invalid user name: %v", err)
	case errors.Is(err, domain.ErrInvalidUserEmail):
		return status.Errorf(codes.InvalidArgument, "invalid user email: %v", err)
	case errors.Is(err, domain.ErrInvalidEmail):
		return status.Errorf(codes.InvalidArgument, "invalid email format: %v", err)
	default:
		return status.Errorf(codes.Internal, "internal error: %v", err)
	}
}

// validateGetUserRequest validates GetUserRequest
func validateGetUserRequest(req interface{}) error {
	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}
	return nil
}

// validateCreateUserRequest validates CreateUserRequest
func validateCreateUserRequest(req interface{}) error {
	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}
	return nil
}

// validateUpdateUserRequest validates UpdateUserRequest
func validateUpdateUserRequest(req interface{}) error {
	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}
	return nil
}

// validateDeleteUserRequest validates DeleteUserRequest
func validateDeleteUserRequest(req interface{}) error {
	if req == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}
	return nil
}

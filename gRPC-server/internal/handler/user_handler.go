package handler

import (
	"context"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/adapter"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
	pb "github.com/longtrd/grpc-api-gateway/gRPC-server/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserHandler implements the gRPC UserServiceServer
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	useCase domain.UserUseCase
	adapter *adapter.UserAdapter
	logger  domain.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(
	useCase domain.UserUseCase,
	adapter *adapter.UserAdapter,
	logger domain.Logger,
) *UserHandler {
	return &UserHandler{
		useCase: useCase,
		adapter: adapter,
		logger:  logger,
	}
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Check if context is already cancelled
	if ctx.Err() != nil {
		h.logger.Warn("request context cancelled", "method", "GetUser", "error", ctx.Err())
		return nil, status.Errorf(codes.Canceled, "request cancelled")
	}

	// Validate request
	if req == nil {
		h.logger.Warn("received nil request", "method", "GetUser")
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Id == "" {
		h.logger.Warn("received empty user ID", "method", "GetUser")
		return nil, status.Errorf(codes.InvalidArgument, "user ID cannot be empty")
	}

	h.logger.Debug("getting user", "user_id", req.Id)

	// Call use case
	user, err := h.useCase.GetUser(ctx, req.Id)
	if err != nil {
		h.logger.Error("failed to get user", "user_id", req.Id, "error", err)
		return nil, mapDomainError(err)
	}

	// Convert to proto and return
	pbUser := h.adapter.ToProto(user)
	h.logger.Debug("user retrieved successfully", "user_id", req.Id)

	return &pb.GetUserResponse{User: pbUser}, nil
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Check if context is already cancelled
	if ctx.Err() != nil {
		h.logger.Warn("request context cancelled", "method", "CreateUser", "error", ctx.Err())
		return nil, status.Errorf(codes.Canceled, "request cancelled")
	}

	// Validate request
	if req == nil {
		h.logger.Warn("received nil request", "method", "CreateUser")
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.User == nil {
		h.logger.Warn("received nil user", "method", "CreateUser")
		return nil, status.Errorf(codes.InvalidArgument, "user cannot be nil")
	}

	// Convert from proto
	user := h.adapter.FromProto(req.User)
	h.logger.Debug("creating user", "user_id", user.ID)

	// Call use case
	err := h.useCase.CreateUser(ctx, user)
	if err != nil {
		h.logger.Error("failed to create user", "user_id", user.ID, "error", err)
		return nil, mapDomainError(err)
	}

	// Convert to proto and return
	pbUser := h.adapter.ToProto(user)
	h.logger.Info("user created successfully", "user_id", user.ID)

	return &pb.CreateUserResponse{User: pbUser}, nil
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// Check if context is already cancelled
	if ctx.Err() != nil {
		h.logger.Warn("request context cancelled", "method", "UpdateUser", "error", ctx.Err())
		return nil, status.Errorf(codes.Canceled, "request cancelled")
	}

	// Validate request
	if req == nil {
		h.logger.Warn("received nil request", "method", "UpdateUser")
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.User == nil {
		h.logger.Warn("received nil user", "method", "UpdateUser")
		return nil, status.Errorf(codes.InvalidArgument, "user cannot be nil")
	}

	// Convert from proto
	user := h.adapter.FromProto(req.User)
	h.logger.Debug("updating user", "user_id", user.ID)

	// Call use case
	err := h.useCase.UpdateUser(ctx, user)
	if err != nil {
		h.logger.Error("failed to update user", "user_id", user.ID, "error", err)
		return nil, mapDomainError(err)
	}

	// Convert to proto and return
	pbUser := h.adapter.ToProto(user)
	h.logger.Info("user updated successfully", "user_id", user.ID)

	return &pb.UpdateUserResponse{User: pbUser}, nil
}

// DeleteUser deletes a user by ID
func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Check if context is already cancelled
	if ctx.Err() != nil {
		h.logger.Warn("request context cancelled", "method", "DeleteUser", "error", ctx.Err())
		return nil, status.Errorf(codes.Canceled, "request cancelled")
	}

	// Validate request
	if req == nil {
		h.logger.Warn("received nil request", "method", "DeleteUser")
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Id == "" {
		h.logger.Warn("received empty user ID", "method", "DeleteUser")
		return nil, status.Errorf(codes.InvalidArgument, "user ID cannot be empty")
	}

	h.logger.Debug("deleting user", "user_id", req.Id)

	// Call use case
	err := h.useCase.DeleteUser(ctx, req.Id)
	if err != nil {
		h.logger.Error("failed to delete user", "user_id", req.Id, "error", err)
		return nil, mapDomainError(err)
	}

	h.logger.Info("user deleted successfully", "user_id", req.Id)

	return &pb.DeleteUserResponse{Id: req.Id}, nil
}

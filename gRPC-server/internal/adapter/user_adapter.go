package adapter

import (
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
	pb "github.com/longtrd/grpc-api-gateway/gRPC-server/proto"
)

// UserAdapter handles conversion between protobuf and domain models
type UserAdapter struct{}

// NewUserAdapter creates a new UserAdapter
func NewUserAdapter() *UserAdapter {
	return &UserAdapter{}
}

// ToProto converts a domain User to protobuf User
func (a *UserAdapter) ToProto(user *domain.User) *pb.User {
	if user == nil {
		return nil
	}

	return &pb.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

// FromProto converts a protobuf User to domain User
func (a *UserAdapter) FromProto(pbUser *pb.User) *domain.User {
	if pbUser == nil {
		return nil
	}

	return &domain.User{
		ID:    pbUser.Id,
		Name:  pbUser.Name,
		Email: pbUser.Email,
	}
}

// ToProtoUsers converts a slice of domain Users to protobuf Users
func (a *UserAdapter) ToProtoUsers(users []*domain.User) []*pb.User {
	if users == nil {
		return nil
	}

	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = a.ToProto(user)
	}
	return pbUsers
}

// FromProtoUsers converts a slice of protobuf Users to domain Users
func (a *UserAdapter) FromProtoUsers(pbUsers []*pb.User) []*domain.User {
	if pbUsers == nil {
		return nil
	}

	users := make([]*domain.User, len(pbUsers))
	for i, pbUser := range pbUsers {
		users[i] = a.FromProto(pbUser)
	}
	return users
}

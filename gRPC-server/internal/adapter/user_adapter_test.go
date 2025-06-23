package adapter

import (
	"testing"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
	pb "github.com/longtrd/grpc-api-gateway/gRPC-server/proto"
)

func TestUserAdapter_ToProto(t *testing.T) {
	adapter := NewUserAdapter()

	tests := []struct {
		name     string
		user     *domain.User
		expected *pb.User
	}{
		{
			name: "valid user",
			user: &domain.User{
				ID:    "123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			expected: &pb.User{
				Id:    "123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
		},
		{
			name:     "nil user",
			user:     nil,
			expected: nil,
		},
		{
			name: "empty user",
			user: &domain.User{},
			expected: &pb.User{
				Id:    "",
				Name:  "",
				Email: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.ToProto(tt.user)

			if tt.expected == nil && result != nil {
				t.Errorf("ToProto() = %v, expected nil", result)
				return
			}

			if tt.expected != nil && result == nil {
				t.Errorf("ToProto() = nil, expected %v", tt.expected)
				return
			}

			if tt.expected != nil && result != nil {
				if result.Id != tt.expected.Id {
					t.Errorf("ToProto().Id = %v, expected %v", result.Id, tt.expected.Id)
				}
				if result.Name != tt.expected.Name {
					t.Errorf("ToProto().Name = %v, expected %v", result.Name, tt.expected.Name)
				}
				if result.Email != tt.expected.Email {
					t.Errorf("ToProto().Email = %v, expected %v", result.Email, tt.expected.Email)
				}
			}
		})
	}
}

func TestUserAdapter_FromProto(t *testing.T) {
	adapter := NewUserAdapter()

	tests := []struct {
		name     string
		pbUser   *pb.User
		expected *domain.User
	}{
		{
			name: "valid proto user",
			pbUser: &pb.User{
				Id:    "123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			expected: &domain.User{
				ID:    "123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
		},
		{
			name:     "nil proto user",
			pbUser:   nil,
			expected: nil,
		},
		{
			name:   "empty proto user",
			pbUser: &pb.User{},
			expected: &domain.User{
				ID:    "",
				Name:  "",
				Email: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.FromProto(tt.pbUser)

			if tt.expected == nil && result != nil {
				t.Errorf("FromProto() = %v, expected nil", result)
				return
			}

			if tt.expected != nil && result == nil {
				t.Errorf("FromProto() = nil, expected %v", tt.expected)
				return
			}

			if tt.expected != nil && result != nil {
				if result.ID != tt.expected.ID {
					t.Errorf("FromProto().ID = %v, expected %v", result.ID, tt.expected.ID)
				}
				if result.Name != tt.expected.Name {
					t.Errorf("FromProto().Name = %v, expected %v", result.Name, tt.expected.Name)
				}
				if result.Email != tt.expected.Email {
					t.Errorf("FromProto().Email = %v, expected %v", result.Email, tt.expected.Email)
				}
			}
		})
	}
}

func TestUserAdapter_ToProtoUsers(t *testing.T) {
	adapter := NewUserAdapter()

	tests := []struct {
		name     string
		users    []*domain.User
		expected []*pb.User
	}{
		{
			name: "multiple users",
			users: []*domain.User{
				{ID: "1", Name: "User 1", Email: "user1@example.com"},
				{ID: "2", Name: "User 2", Email: "user2@example.com"},
			},
			expected: []*pb.User{
				{Id: "1", Name: "User 1", Email: "user1@example.com"},
				{Id: "2", Name: "User 2", Email: "user2@example.com"},
			},
		},
		{
			name:     "nil slice",
			users:    nil,
			expected: nil,
		},
		{
			name:     "empty slice",
			users:    []*domain.User{},
			expected: []*pb.User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.ToProtoUsers(tt.users)

			if tt.expected == nil && result != nil {
				t.Errorf("ToProtoUsers() = %v, expected nil", result)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("ToProtoUsers() length = %v, expected %v", len(result), len(tt.expected))
				return
			}

			for i, pbUser := range result {
				expected := tt.expected[i]
				if pbUser.Id != expected.Id || pbUser.Name != expected.Name || pbUser.Email != expected.Email {
					t.Errorf("ToProtoUsers()[%d] = %v, expected %v", i, pbUser, expected)
				}
			}
		})
	}
}

func TestUserAdapter_FromProtoUsers(t *testing.T) {
	adapter := NewUserAdapter()

	tests := []struct {
		name     string
		pbUsers  []*pb.User
		expected []*domain.User
	}{
		{
			name: "multiple proto users",
			pbUsers: []*pb.User{
				{Id: "1", Name: "User 1", Email: "user1@example.com"},
				{Id: "2", Name: "User 2", Email: "user2@example.com"},
			},
			expected: []*domain.User{
				{ID: "1", Name: "User 1", Email: "user1@example.com"},
				{ID: "2", Name: "User 2", Email: "user2@example.com"},
			},
		},
		{
			name:     "nil slice",
			pbUsers:  nil,
			expected: nil,
		},
		{
			name:     "empty slice",
			pbUsers:  []*pb.User{},
			expected: []*domain.User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.FromProtoUsers(tt.pbUsers)

			if tt.expected == nil && result != nil {
				t.Errorf("FromProtoUsers() = %v, expected nil", result)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("FromProtoUsers() length = %v, expected %v", len(result), len(tt.expected))
				return
			}

			for i, user := range result {
				expected := tt.expected[i]
				if user.ID != expected.ID || user.Name != expected.Name || user.Email != expected.Email {
					t.Errorf("FromProtoUsers()[%d] = %v, expected %v", i, user, expected)
				}
			}
		})
	}
}

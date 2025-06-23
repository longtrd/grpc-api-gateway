package domain

import (
	"testing"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user",
			user: User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			user: User{
				ID:    "",
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			wantErr: true,
			errMsg:  "user ID cannot be empty",
		},
		{
			name: "empty name",
			user: User{
				ID:    "user-123",
				Name:  "",
				Email: "john.doe@example.com",
			},
			wantErr: true,
			errMsg:  "user name cannot be empty",
		},
		{
			name: "empty email",
			user: User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "",
			},
			wantErr: true,
			errMsg:  "user email cannot be empty",
		},
		{
			name: "invalid email format",
			user: User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("User.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUser_String(t *testing.T) {
	user := User{
		ID:    "user-123",
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	expected := "User{ID: user-123, Name: John Doe, Email: john.doe@example.com}"
	if got := user.String(); got != expected {
		t.Errorf("User.String() = %v, want %v", got, expected)
	}
}

func TestUser_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		user User
		want bool
	}{
		{
			name: "empty user",
			user: User{},
			want: true,
		},
		{
			name: "user with only ID",
			user: User{ID: "user-123"},
			want: false,
		},
		{
			name: "complete user",
			user: User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.IsEmpty(); got != tt.want {
				t.Errorf("User.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

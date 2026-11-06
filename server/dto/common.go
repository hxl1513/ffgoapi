package dto

import (
	"fmt"
	"time"

	"github.com/your-org/go-web-api-framework/model"
)

// FromUser converts a model.User (domain entity) into a UserResponse DTO.
// It never exposes the password hash.
func FromUser(user *model.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}

// MustFromUser is like FromUser but panics on a nil user.  Useful in tests.
func MustFromUser(user *model.User) *UserResponse {
	resp := FromUser(user)
	if resp == nil {
		panic("FromUser: user is nil")
	}
	return resp
}

// StringPtrToUint converts a *uint (often from a DB query) to uint,
// returning 0 when the pointer is nil.
func StringPtrToUint(s *string) uint {
	if s == nil {
		return 0
	}
	var n uint
	fmt.Sscanf(*s, "%d", &n)
	return n
}

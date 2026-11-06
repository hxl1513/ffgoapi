package dto

import "github.com/go-playground/validator/v10"

// CreateUserRequest represents the payload for creating a new user.
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Validate runs struct-field validation on the request and returns a descriptive
// error if any rule is violated.
func (r *CreateUserRequest) Validate() error {
	v := validator.New()
	if err := v.Struct(r); err != nil {
		return err
	}
	return nil
}

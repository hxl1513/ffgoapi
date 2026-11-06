package dto

// UserResponse is the safe, serialisable representation of a user returned by
// the API. It deliberately omits the password hash.
type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      int    `json:"role"`
	CreatedAt string `json:"created_at"`
}

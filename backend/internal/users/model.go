package users

import "time"

// User represents the full database entity
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never expose password in JSON
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest DTO
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// UpdateUserRequest DTO
type UpdateUserRequest struct {
	FullName string  `json:"full_name"`
	Role     string  `json:"role"`
	Password *string `json:"password,omitempty"` // Optional new password
}

// ToggleStatusRequest DTO
type ToggleStatusRequest struct {
	IsActive bool `json:"is_active"`
}


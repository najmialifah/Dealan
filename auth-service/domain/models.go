package domain

import "time"

type UserRole string

const (
	RoleUser   UserRole = "user"
	RoleDriver UserRole = "driver"
)

type AuthUser struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

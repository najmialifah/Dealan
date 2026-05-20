package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Photo     string    `json:"photo,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Photo string `json:"photo"`
}

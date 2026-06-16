package repository

import (
	"context"

	"github.com/shakilaaulia/Dealan/user-service/domain"
)

// UserRepository mendefinisikan kontrak akses data GORM untuk user-service
type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	CreateUser(ctx context.Context, user *domain.User) error
}

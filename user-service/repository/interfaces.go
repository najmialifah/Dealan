package repository

import (
	"context"

	"github.com/shakilaaulia/Dealan/user-service/domain"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mock_user_repository.go -package=mocks

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
}

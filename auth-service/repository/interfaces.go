package repository

import (
	"context"

	"github.com/shakilaaulia/Dealan/auth-service/domain"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mock_auth_repository.go -package=mocks

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.AuthUser, error)
	CreateUser(ctx context.Context, user *domain.AuthUser) error
}

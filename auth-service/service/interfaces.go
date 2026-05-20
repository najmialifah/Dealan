package service

import (
	"context"

	"github.com/shakilaaulia/Dealan/auth-service/domain"
)

type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*domain.AuthUser, error)
}

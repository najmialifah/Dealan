package service

import (
	"context"
	"errors"

	"github.com/shakilaaulia/Dealan/auth-service/domain"
	"github.com/shakilaaulia/Dealan/auth-service/repository"
)

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Draft logic
	return nil, errors.New("not implemented")
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	// Draft logic: get user
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// Draft logic: compare password
	if user.Password != req.Password {
		return nil, errors.New("invalid credentials")
	}

	// Draft logic: create token
	return &domain.AuthResponse{Token: "draft-jwt-token"}, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*domain.AuthUser, error) {
	// Draft logic
	return nil, errors.New("not implemented")
}

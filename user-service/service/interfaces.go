package service

import (
	"context"

	"github.com/shakilaaulia/Dealan/user-service/domain"
)

type UserService interface {
	GetProfile(ctx context.Context, id string) (*domain.User, error)
	UpdateProfile(ctx context.Context, id string, req domain.UpdateProfileRequest) error
	GetInternalName(ctx context.Context, id string) (string, error)
}

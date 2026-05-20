package service

import (
	"context"
	"errors"

	"github.com/shakilaaulia/Dealan/user-service/domain"
	"github.com/shakilaaulia/Dealan/user-service/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	// Draft logic
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, id string, req domain.UpdateProfileRequest) error {
	// Draft logic
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	user.Name = req.Name
	user.Phone = req.Phone
	user.Photo = req.Photo

	return s.repo.UpdateUser(ctx, user)
}

func (s *userService) GetInternalName(ctx context.Context, id string) (string, error) {
	// Draft logic
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	return user.Name, nil
}

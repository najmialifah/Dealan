package service

import (
	"context"
	"errors"

	"github.com/najmialifah/Dealan/user-service/domain"
	"github.com/najmialifah/Dealan/user-service/repository"
)

type userService struct {
	repo repository.UserRepository
}

// NewUserService membuat instance baru dari bisnis service user
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// GetProfile mengambil profil lengkap milik user
func (s *userService) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateProfile mengubah field profil yang diperbolehkan
func (s *userService) UpdateProfile(ctx context.Context, id string, req domain.UpdateProfileRequest) error {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	user.Nama = req.Nama
	user.NomorHP = req.NomorHP
	user.Alamat = req.Alamat
	user.FotoProfil = req.FotoProfil

	return s.repo.UpdateUser(ctx, user)
}

// GetInternalName mengembalikan nama user untuk konsumsi internal service lain
func (s *userService) GetInternalName(ctx context.Context, id string) (string, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user tidak ditemukan")
	}
	return user.Nama, nil
}

// CreateUser mendaftarkan user baru yang dipicu oleh event Kafka dari auth-service
func (s *userService) CreateUser(ctx context.Context, id, name, email, phone string) error {
	// Cek apakah user sudah terdaftar di DB
	existing, _ := s.repo.GetUserByID(ctx, id)
	if existing != nil {
		return nil // idempotent
	}

	user := &domain.User{
		ID:      id,
		Nama:    name,
		Email:   email,
		NomorHP: phone,
		Status:  "active",
	}
	return s.repo.CreateUser(ctx, user)
}

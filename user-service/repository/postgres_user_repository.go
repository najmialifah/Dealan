package repository

import (
	"context"

	"github.com/shakilaaulia/Dealan/user-service/domain"
	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

// NewPostgresUserRepository membuat instance baru UserRepository berbasis GORM PostgreSQL
func NewPostgresUserRepository(db *gorm.DB) UserRepository {
	return &postgresUserRepository{
		db: db,
	}
}

// GetUserByID mengambil data user beserta relasi profile dan rating menggunakan Preload GORM
func (r *postgresUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Preload("UserProfile").Preload("UserRating").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser memperbarui data user ke database
func (r *postgresUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// CreateUser mendaftarkan user baru secara transaksional beserta profile dan rating default-nya
func (r *postgresUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Simpan user utama
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// Buat record UserProfile default
		profile := &domain.UserProfile{
			UserID:     user.ID,
			Preferensi: make(domain.JSONMap),
			TotalOrder: 0,
			TotalSpend: 0.0,
		}
		if err := tx.Create(profile).Error; err != nil {
			return err
		}

		// Buat record UserRating default
		rating := &domain.UserRating{
			UserID:      user.ID,
			TotalRating: 5.00,
			TotalReview: 0,
		}
		if err := tx.Create(rating).Error; err != nil {
			return err
		}

		return nil
	})
}

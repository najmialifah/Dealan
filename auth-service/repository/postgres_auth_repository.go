package repository

import (
	"context"
	"time"

	"github.com/najmialifah/Dealan/auth-service/domain"
	"gorm.io/gorm"
)

type postgresAuthRepository struct {
	db *gorm.DB
}

// NewPostgresAuthRepository membuat instance baru dari repository GORM PostgreSQL
func NewPostgresAuthRepository(db *gorm.DB) AuthRepository {
	return &postgresAuthRepository{
		db: db,
	}
}

// GetCredentialByEmail mencari data kredensial berdasarkan email
func (r *postgresAuthRepository) GetCredentialByEmail(ctx context.Context, email string) (*domain.AuthCredential, error) {
	var cred domain.AuthCredential
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&cred).Error
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

// CreateCredential menyimpan data kredensial baru ke database
func (r *postgresAuthRepository) CreateCredential(ctx context.Context, cred *domain.AuthCredential) error {
	return r.db.WithContext(ctx).Create(cred).Error
}

// CreateRefreshToken menyimpan token refresh baru
func (r *postgresAuthRepository) CreateRefreshToken(ctx context.Context, rt *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Create(rt).Error
}

// RevokeRefreshToken mencabut token refresh dengan memberikan timestamp pada revoked_at
func (r *postgresAuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked_at", &now).Error
}

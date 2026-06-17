package repository

import (
	"context"

	"github.com/najmialifah/Dealan/auth-service/domain"
)

// AuthRepository mendefinisikan interface untuk interaksi data auth di PostgreSQL menggunakan GORM
type AuthRepository interface {
	GetCredentialByEmail(ctx context.Context, email string) (*domain.AuthCredential, error)
	CreateCredential(ctx context.Context, cred *domain.AuthCredential) error
	CreateRefreshToken(ctx context.Context, rt *domain.RefreshToken) error
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}

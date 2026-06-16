package service

import (
	"context"

	"github.com/shakilaaulia/Dealan/auth-service/domain"
)

// AuthService mendefinisikan interface bisnis untuk layanan autentikasi
type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*domain.AuthCredential, error)
}

// EventProducer mendefinisikan interface untuk mengirim event pendaftaran ke Apache Kafka
type EventProducer interface {
	PublishUserCreated(ctx context.Context, id, name, email, phone string) error
	PublishDriverCreated(ctx context.Context, id, name, email, phone string) error
}

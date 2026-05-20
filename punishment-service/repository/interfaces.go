package repository

import (
	"context"
	"github.com/shakilaaulia/Dealan/punishment-service/domain"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mock_punishment_repository.go -package=mocks

type PunishmentRepository interface {
	StoreViolation(ctx context.Context, req domain.PunishmentRequest) (string, error)
	UpdateAccountStatus(ctx context.Context, id string, status string) error
}
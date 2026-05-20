package service

import (
	"context"
	"github.com/shakilaaulia/Dealan/punishment-service/domain"
)

type PunishmentService interface {
	ApplyPunishment(ctx context.Context, req domain.PunishmentRequest) (*domain.PunishmentResponse, error)
}
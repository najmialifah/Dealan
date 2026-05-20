package service

import (
	"context"
	"errors"
	"github.com/shakilaaulia/Dealan/punishment-service/domain"
	"github.com/shakilaaulia/Dealan/punishment-service/repository"
)

type punishmentService struct {
	repo repository.PunishmentRepository
}

func NewPunishmentService(repo repository.PunishmentRepository) PunishmentService {
	return &punishmentService{repo: repo}
}

func (s *punishmentService) ApplyPunishment(ctx context.Context, req domain.PunishmentRequest) (*domain.PunishmentResponse, error) {

	return nil, errors.New("not implemented")
}
package service

import (
	"context"
	"errors"

	"github.com/najmialifah/Dealan/punishment-service/domain"
	"github.com/najmialifah/Dealan/punishment-service/repository"
)

type punishmentService struct {
	repo repository.PunishmentRepository
}

// NewPunishmentService membuat instance baru dari service punishment
func NewPunishmentService(repo repository.PunishmentRepository) PunishmentService {
	return &punishmentService{repo: repo}
}

// ApplyPunishment memproses pemberian sanksi, menyimpan log pelanggaran, dan mengupdate status akun
func (s *punishmentService) ApplyPunishment(ctx context.Context, req domain.PunishmentRequest) (*domain.PunishmentResponse, error) {
	// Validasi input
	if req.AccountID == "" || req.ReasonCode == "" {
		return nil, errors.New("account_id dan reason_code tidak boleh kosong")
	}

	// Tentukan status baru akun berdasarkan durasi sanksi
	var newStatus string
	if req.Duration == -1 {
		newStatus = "Banned"
	} else if req.Duration > 0 {
		newStatus = "Suspended"
	} else {
		newStatus = "Warning"
	}

	// 1. Simpan data pelanggaran ke DB
	penaltyID, err := s.repo.StoreViolation(ctx, req)
	if err != nil {
		return nil, err
	}

	// 2. Update status akun pengguna/driver ke database
	err = s.repo.UpdateAccountStatus(ctx, req.AccountID, newStatus)
	if err != nil {
		return nil, err
	}

	return &domain.PunishmentResponse{
		PenaltyID:        penaltyID,
		NewAccountStatus: newStatus,
		Status:           "Success",
	}, nil
}
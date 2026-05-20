package service

import (
	"errors"
	"matching-service/domain"
)

type MatchingService struct {
	Repo domain.MatchingRepository
}

func (s *MatchingService) FindDriver(req *domain.MatchingRequest) (*domain.MatchingResponse, error) {
	if req.OrderID == "" {
		return nil, errors.New("order ID tidak boleh kosong")
	}
	// Panggil database bohongan untuk cari driver
	return s.Repo.FindNearestDriver(req)
}
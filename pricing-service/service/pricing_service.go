package service

import (
	"context"

	"github.com/shakilaaulia/Dealan/pricing-service/domain"
	"github.com/shakilaaulia/Dealan/pricing-service/repository"
)

type PricingService interface {
	Calculate(ctx context.Context, req domain.PricingRequest) (domain.PricingResponse, error)
	Negotiate(ctx context.Context, req domain.NegotiationData) error
}

type pricingServiceImpl struct {
	repo repository.PricingRepository
}

func NewPricingService(repo repository.PricingRepository) PricingService {
	return &pricingServiceImpl{repo: repo}
}

func (s *pricingServiceImpl) Calculate(ctx context.Context, req domain.PricingRequest) (domain.PricingResponse, error) {
	basePrice, err := s.repo.GetBasePrice(ctx, req.ServiceType)
	if err != nil {
		return domain.PricingResponse{}, err
	}

	harga := basePrice * req.Jarak
	if req.ServiceType == "send" {
		harga += req.BeratBarang * 1000
	}

	return domain.PricingResponse{Harga: harga}, nil
}

func (s *pricingServiceImpl) Negotiate(ctx context.Context, req domain.NegotiationData) error {
	return s.repo.SaveNegotiation(ctx, req)
}

package service

import (
	"context"

	"promo-service/domain"
	"promo-service/repository"
)

type promoService struct {
	repo repository.PromoRepository
}

func NewPromoService(r repository.PromoRepository) PromoService {
	return &promoService{r}
}

func (s *promoService) ApplyPromo(ctx context.Context, req domain.PromoRequest) (*domain.PromoResponse, error) {

	if s.repo == nil {
		return &domain.PromoResponse{
			Discount: 10000,
			IsValid:  true,
		}, nil
	}

	discount, err := s.repo.CheckPromo(ctx, req.Code)
	if err != nil {
		return &domain.PromoResponse{IsValid: false}, nil
	}

	return &domain.PromoResponse{
		Discount: discount,
		IsValid:  true,
	}, nil
}
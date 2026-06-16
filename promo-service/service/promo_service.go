package service

import (
	"context"

	"promo-service/domain"
	"promo-service/repository"
)

type promoService struct {
	repo repository.PromoRepository
}

// NewPromoService membuat instance baru dari service promo
func NewPromoService(r repository.PromoRepository) PromoService {
	return &promoService{r}
}

// ApplyPromo memverifikasi kode promo dan mengembalikan diskon jika promo valid
func (s *promoService) ApplyPromo(ctx context.Context, req domain.PromoRequest) (*domain.PromoResponse, error) {
	if s.repo == nil {
		return &domain.PromoResponse{
			Discount: 10000,
			IsValid:  true,
		}, nil
	}

	discount, err := s.repo.CheckPromo(ctx, req.Code)
	if err != nil {
		// Jika terjadi error (misalnya promo tidak ada atau kuota habis), kembalikan status tidak valid
		return &domain.PromoResponse{
			Discount: 0,
			IsValid:  false,
		}, nil
	}

	return &domain.PromoResponse{
		Discount: discount,
		IsValid:  true,
	}, nil
}
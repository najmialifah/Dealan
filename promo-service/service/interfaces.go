package service

import (
	"context"
	"promo-service/domain"
)

type PromoService interface {
	ApplyPromo(ctx context.Context, req domain.PromoRequest) (*domain.PromoResponse, error)
}
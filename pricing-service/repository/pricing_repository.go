package repository

import (
	"context"

	"github.com/shakilaaulia/Dealan/pricing-service/domain"
)

//go:generate mockgen -source=pricing_repository.go -destination=../mocks/mock_pricing_repository.go -package=mocks

type PricingRepository interface {
	GetBasePrice(ctx context.Context, serviceType string) (float64, error)
	SaveNegotiation(ctx context.Context, data domain.NegotiationData) error
}

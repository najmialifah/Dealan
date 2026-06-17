package service

import (
	"context"

	"github.com/najmialifah/Dealan/rating-review-service/domain"
)

// RatingService mendefinisikan kontrak logika bisnis untuk Rating & Review
type RatingService interface {
	SubmitReview(ctx context.Context, req domain.RatingRequest) (*domain.RatingResponse, error)
}
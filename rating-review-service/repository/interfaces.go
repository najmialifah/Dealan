package repository

import (
	"context"

	"github.com/najmialifah/Dealan/rating-review-service/domain"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mock_rating_repository.go -package=mocks

type RatingRepository interface {
	SaveReview(ctx context.Context, req domain.RatingRequest) (string, error)
	GetAverageRating(ctx context.Context, targetID string) (float64, error)
}
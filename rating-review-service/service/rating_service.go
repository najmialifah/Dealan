package service

import (
	"context"
	"errors"
	"github.com/shakilaaulia/Dealan/rating-review-service/domain"
	"github.com/shakilaaulia/Dealan/rating-review-service/repository"
)


type ratingService struct {
	repo repository.RatingRepository
}

func NewRatingService(repo repository.RatingRepository) RatingService {
	return &ratingService{repo: repo}
}

func (s *ratingService) SubmitReview(ctx context.Context, req domain.RatingRequest) (*domain.RatingResponse, error) {
	// Masih Draft Tahap 2
	return nil, errors.New("not implemented")
}
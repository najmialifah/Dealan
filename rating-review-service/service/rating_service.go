package service

import (
	"context"
	"errors"

	"github.com/najmialifah/Dealan/rating-review-service/domain"
	"github.com/najmialifah/Dealan/rating-review-service/repository"
)

type ratingService struct {
	repo repository.RatingRepository
}

// NewRatingService membuat instance baru dari service rating
func NewRatingService(repo repository.RatingRepository) RatingService {
	return &ratingService{repo: repo}
}

// SubmitReview memproses penyimpanan review dan menghitung rating rata-rata instan
func (s *ratingService) SubmitReview(ctx context.Context, req domain.RatingRequest) (*domain.RatingResponse, error) {
	// Validasi input sederhana
	if req.OrderID == "" || req.ReviewerID == "" || req.TargetID == "" {
		return nil, errors.New("order_id, reviewer_id, dan target_id tidak boleh kosong")
	}
	if req.RatingScore < 1 || req.RatingScore > 5 {
		return nil, errors.New("skor rating harus antara 1 dan 5")
	}

	// 1. Simpan review baru ke DB
	reviewID, err := s.repo.SaveReview(ctx, req)
	if err != nil {
		return nil, err
	}

	// 2. Dapatkan nilai rata-rata rating target secara instan menggunakan query AVG()
	avg, err := s.repo.GetAverageRating(ctx, req.TargetID)
	if err != nil {
		return nil, err
	}

	return &domain.RatingResponse{
		AverageRating: avg,
		ReviewID:      reviewID,
		Status:        "Success",
	}, nil
}
package service_test

import (
	"context"
	"testing"

	"github.com/najmialifah/Dealan/rating-review-service/domain"
	"github.com/najmialifah/Dealan/rating-review-service/mocks"
	"github.com/najmialifah/Dealan/rating-review-service/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRatingService_SubmitReview_Success(t *testing.T) {
	// 1. Setup Mock Controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 2. Setup Mock Repository
	mockRepo := mocks.NewMockRatingRepository(ctrl)

	// 3. Initialize Service dengan Mock Repository
	svc := service.NewRatingService(mockRepo)

	req := domain.RatingRequest{
		OrderID:      "ORDER-100",
		ReviewerID:   "USER-001",
		ReviewerRole: "user",
		TargetID:     "DRIVER-001",
		TargetRole:   "driver",
		RatingScore:  5,
		Comment:      "Driver sangat ramah!",
	}

	// Ekspektasikan pemanggilan SaveReview dan GetAverageRating
	mockRepo.EXPECT().
		SaveReview(gomock.Any(), req).
		Return("review-uuid-123", nil).
		Times(1)

	mockRepo.EXPECT().
		GetAverageRating(gomock.Any(), "DRIVER-001").
		Return(4.8, nil).
		Times(1)

	// 4. Eksekusi
	res, err := svc.SubmitReview(context.Background(), req)

	// 5. Assertion (Pengecekan)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 4.8, res.AverageRating)
	assert.Equal(t, "review-uuid-123", res.ReviewID)
	assert.Equal(t, "Success", res.Status)
}

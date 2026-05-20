package service

import (
	"context"
	"testing"

	"github.com/shakilaaulia/Dealan/rating-review-service/domain"
	"github.com/shakilaaulia/Dealan/rating-review-service/mocks"
	"go.uber.org/mock/gomock"
)

func TestRatingService_SubmitReview_ExpectedFail(t *testing.T) {
	// 1. Setup Mock Controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 2. Setup Mock Repository
	mockRepo := mocks.NewMockRatingRepository(ctrl)
	
	// 3. Initialize Service dengan Mock Repository
	svc := NewRatingService(mockRepo)

	// 4. Data Dummy untuk Testing
	req := domain.RatingRequest{
		OrderID:     "ORDER-100",
		ReviewerID:  "USER-001",
		RatingScore: 5,
		Comment:     "Driver sangat ramah!",
	}

	// 5. Eksekusi Fungsi yang Dites
	ctx := context.Background()
	res, err := svc.SubmitReview(ctx, req)

	// 6. Assertion (Pengecekan)
	if err == nil {
		t.Error("Ekspektasi: Error 'not implemented', Realita: Tidak ada error")
	}

	if err != nil && err.Error() != "not implemented" {
		t.Errorf("Ekspektasi error message 'not implemented', Realita: %v", err.Error())
	}

	if res != nil {
		t.Error("Ekspektasi response nil karena masih draft, Realita: Response tidak nil")
	}
}



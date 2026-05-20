package service

import (
	"context"
	"testing"

	"promo-service/domain"
	"promo-service/mocks"

	"go.uber.org/mock/gomock"
)

func TestPromo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPromoRepository(ctrl)

	mockRepo.EXPECT().
		CheckPromo(gomock.Any(), "DISKON").
		Return(10000.0, nil)

	svc := NewPromoService(mockRepo)

	res, _ := svc.ApplyPromo(context.Background(), domain.PromoRequest{
		Code: "DISKON",
	})

	if res.Discount != 10000 {
		t.Errorf("wrong discount")
	}
}
package service_test

import (
	"context"
	"errors"
	"testing"

	"promo-service/domain"
	"promo-service/mocks"
	"promo-service/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPromoService_ApplyPromo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPromoRepository(ctrl)
	svc := service.NewPromoService(mockRepo)

	mockRepo.EXPECT().
		CheckPromo(gomock.Any(), "DISKON50").
		Return(15000.0, nil).
		Times(1)

	res, err := svc.ApplyPromo(context.Background(), domain.PromoRequest{
		Code: "DISKON50",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.IsValid)
	assert.Equal(t, 15000.0, res.Discount)
}

func TestPromoService_ApplyPromo_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPromoRepository(ctrl)
	svc := service.NewPromoService(mockRepo)

	mockRepo.EXPECT().
		CheckPromo(gomock.Any(), "INVALIDCODE").
		Return(0.0, errors.New("promo not found")).
		Times(1)

	res, err := svc.ApplyPromo(context.Background(), domain.PromoRequest{
		Code: "INVALIDCODE",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.False(t, res.IsValid)
	assert.Equal(t, 0.0, res.Discount)
}
package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shakilaaulia/Dealan/pricing-service/domain"
	"github.com/shakilaaulia/Dealan/pricing-service/mocks"
	"github.com/shakilaaulia/Dealan/pricing-service/service"
	"go.uber.org/mock/gomock"
)

func TestPricingService_Calculate_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPricingRepository(ctrl)
	srv := service.NewPricingService(mockRepo)

	ctx := context.Background()
	req := domain.PricingRequest{
		ServiceType: "ride",
		Jarak:       10,
	}

	expectedErr := errors.New("failed to get base price")
	mockRepo.EXPECT().GetBasePrice(gomock.Any(), "ride").Return(0.0, expectedErr).Times(1)

	_, err := srv.Calculate(ctx, req)

	if err == nil {
		t.Errorf("Expected error but got nil")
	}
	if !errors.Is(err, expectedErr) && err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

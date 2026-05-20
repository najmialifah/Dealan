package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shakilaaulia/Dealan/payment-service/domain"
	"github.com/shakilaaulia/Dealan/payment-service/mocks"
	"github.com/shakilaaulia/Dealan/payment-service/service"
	"go.uber.org/mock/gomock"
)

func TestPaymentService_Process_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPaymentRepository(ctrl)
	srv := service.NewPaymentService(mockRepo)

	ctx := context.Background()
	req := domain.PaymentRequest{
		OrderID:          "ORD-123",
		Nominal:          100000,
		MetodePembayaran: "QRIS",
		PromoID:          "",
	}

	expectedErr := errors.New("failed to save transaction")
	mockRepo.EXPECT().SaveTransaction(gomock.Any(), req).Return("", expectedErr).Times(1)

	_, err := srv.Process(ctx, req)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
	if !errors.Is(err, expectedErr) && err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

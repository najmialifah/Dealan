package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shakilaaulia/Dealan/shipment-service/domain"
	"github.com/shakilaaulia/Dealan/shipment-service/mocks"
	"github.com/shakilaaulia/Dealan/shipment-service/service"
	"go.uber.org/mock/gomock"
)

func TestShipmentService_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockShipmentRepository(ctrl)
	srv := service.NewShipmentService(mockRepo)

	ctx := context.Background()
	req := domain.ShipmentRequest{
		OrderID:        "ORD-123",
		KategoriBarang: "Electronics",
		BeratBarang:    2.5,
		NamaPenerima:   "John Doe",
		NomorPenerima:  "081234567890",
		CatatanPickup:  "Fragile",
	}

	expectedErr := errors.New("failed to create shipment")
	mockRepo.EXPECT().CreateShipment(gomock.Any(), req).Return("", "", expectedErr).Times(1)

	_, err := srv.Create(ctx, req)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
	if !errors.Is(err, expectedErr) && err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

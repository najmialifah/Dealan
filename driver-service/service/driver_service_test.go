package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shakilaaulia/Dealan/driver-service/domain"
	"github.com/shakilaaulia/Dealan/driver-service/mocks"
	"github.com/shakilaaulia/Dealan/driver-service/service"
	"go.uber.org/mock/gomock"
)

func TestDriverService_UpdateLocation_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDriverRepository(ctrl)
	driverSvc := service.NewDriverService(mockRepo)

	expectedErr := errors.New("db error")
	mockRepo.EXPECT().
		GetDriverByID(gomock.Any(), "drv-1").
		Return(nil, expectedErr).
		Times(1)

	req := domain.UpdateLocationRequest{Lat: -6.2, Long: 106.8}
	err := driverSvc.UpdateLocation(context.Background(), "drv-1", req)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) && err.Error() != expectedErr.Error() {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestDriverService_UpdateLocation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDriverRepository(ctrl)
	driverSvc := service.NewDriverService(mockRepo)

	existingDriver := &domain.Driver{ID: "drv-1", Lat: 0, Long: 0}
	
	mockRepo.EXPECT().
		GetDriverByID(gomock.Any(), "drv-1").
		Return(existingDriver, nil).
		Times(1)

	mockRepo.EXPECT().
		UpdateDriver(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	req := domain.UpdateLocationRequest{Lat: -6.2, Long: 106.8}
	err := driverSvc.UpdateLocation(context.Background(), "drv-1", req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

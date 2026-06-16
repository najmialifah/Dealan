package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shakilaaulia/Dealan/driver-service/domain"
	"github.com/shakilaaulia/Dealan/driver-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDriverRepository adalah mock untuk repository.DriverRepository menggunakan testify/mock
type MockDriverRepository struct {
	mock.Mock
}

func (m *MockDriverRepository) GetDriverByID(ctx context.Context, id string) (*domain.Driver, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *MockDriverRepository) UpdateDriver(ctx context.Context, driver *domain.Driver) error {
	args := m.Called(ctx, driver)
	return args.Error(0)
}

func (m *MockDriverRepository) CreateDriver(ctx context.Context, driver *domain.Driver) error {
	args := m.Called(ctx, driver)
	return args.Error(0)
}

func (m *MockDriverRepository) UpdateDriverLocation(ctx context.Context, driverID string, lat, long float64) error {
	args := m.Called(ctx, driverID, lat, long)
	return args.Error(0)
}

func (m *MockDriverRepository) UpdateDriverStatus(ctx context.Context, driverID string, isOnline bool, layananAktif string) error {
	args := m.Called(ctx, driverID, isOnline, layananAktif)
	return args.Error(0)
}

func TestDriverService_GetProfile_Success(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	driverSvc := service.NewDriverService(mockRepo)

	expectedDriver := &domain.Driver{
		ID:      "driver-123",
		Nama:    "Joko Susilo",
		NomorHP: "0855443322",
	}

	mockRepo.On("GetDriverByID", mock.Anything, "driver-123").Return(expectedDriver, nil)

	profile, err := driverSvc.GetProfile(context.Background(), "driver-123")

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, expectedDriver.Nama, profile.Nama)
	mockRepo.AssertExpectations(t)
}

func TestDriverService_UpdateLocation_Success(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	driverSvc := service.NewDriverService(mockRepo)

	req := domain.UpdateLocationRequest{
		Lat:  -6.1754,
		Long: 106.8272,
	}

	mockRepo.On("UpdateDriverLocation", mock.Anything, "driver-123", req.Lat, req.Long).Return(nil)

	err := driverSvc.UpdateLocation(context.Background(), "driver-123", req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDriverService_UpdateStatus_Success(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	driverSvc := service.NewDriverService(mockRepo)

	req := domain.UpdateStatusRequest{
		IsOnline:     true,
		LayananAktif: "ride,send",
	}

	mockRepo.On("UpdateDriverStatus", mock.Anything, "driver-123", req.IsOnline, req.LayananAktif).Return(nil)

	err := driverSvc.UpdateStatus(context.Background(), "driver-123", req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDriverService_CreateDriver_Success(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	driverSvc := service.NewDriverService(mockRepo)

	mockRepo.On("GetDriverByID", mock.Anything, "new-driver-id").Return(nil, errors.New("record not found"))
	mockRepo.On("CreateDriver", mock.Anything, mock.MatchedBy(func(driver *domain.Driver) bool {
		return driver.ID == "new-driver-id" && driver.Nama == "Supriadi"
	})).Return(nil)

	err := driverSvc.CreateDriver(context.Background(), "new-driver-id", "Supriadi", "supri@example.com", "08111222333")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDriverService_AddVehicle_Success(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	driverSvc := service.NewDriverService(mockRepo)

	existingDriver := &domain.Driver{
		ID:   "driver-123",
		Nama: "Supriadi",
	}

	vehicle := domain.Vehicle{
		Jenis:     "motor",
		PlatNomor: "B 1234 ABC",
	}

	mockRepo.On("GetDriverByID", mock.Anything, "driver-123").Return(existingDriver, nil)
	mockRepo.On("UpdateDriver", mock.Anything, mock.MatchedBy(func(d *domain.Driver) bool {
		return len(d.Vehicles) == 1 && d.Vehicles[0].PlatNomor == "B 1234 ABC"
	})).Return(nil)

	err := driverSvc.AddVehicle(context.Background(), "driver-123", vehicle)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

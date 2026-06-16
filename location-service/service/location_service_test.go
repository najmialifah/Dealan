package service

import (
	"context"
	"location-service/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) UpsertLocation(ctx context.Context, driverID uint, lat, lon float64) error {
	args := m.Called(ctx, driverID, lat, lon)
	return args.Error(0)
}

func (m *MockLocationRepository) FindNearbyDrivers(ctx context.Context, lat, lon float64, radiusMeters float64) ([]models.NearbyDriver, error) {
	args := m.Called(ctx, lat, lon, radiusMeters)
	return args.Get(0).([]models.NearbyDriver), args.Error(1)
}

func TestLocationService_UpdateLocation(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	svc := NewLocationService(mockRepo)

	// Ekspektasi
	mockRepo.On("UpsertLocation", mock.Anything, uint(1), -6.2, 106.8).Return(nil)

	err := svc.UpdateLocation(context.Background(), 1, -6.2, 106.8)
	assert.NoError(t, err)

	// Beri sedikit jeda karena menggunakan goroutine fire-and-forget
	time.Sleep(10 * time.Millisecond)
	mockRepo.AssertExpectations(t)
}

func TestLocationService_GetNearbyDrivers(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	svc := NewLocationService(mockRepo)

	expectedDrivers := []models.NearbyDriver{
		{DriverID: 1, Latitude: -6.2, Longitude: 106.8, Distance: 500},
	}

	mockRepo.On("FindNearbyDrivers", mock.Anything, -6.2, 106.8, float64(5000)).Return(expectedDrivers, nil)

	drivers, err := svc.GetNearbyDrivers(context.Background(), -6.2, 106.8, 5000)

	assert.NoError(t, err)
	assert.Len(t, drivers, 1)
	assert.Equal(t, uint(1), drivers[0].DriverID)
	mockRepo.AssertExpectations(t)
}
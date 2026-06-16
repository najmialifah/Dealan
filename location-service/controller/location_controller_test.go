package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"location-service/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLocationService struct {
	mock.Mock
}

func (m *MockLocationService) UpdateLocation(ctx context.Context, driverID uint, lat, lon float64) error {
	args := m.Called(ctx, driverID, lat, lon)
	return args.Error(0)
}

func (m *MockLocationService) GetNearbyDrivers(ctx context.Context, lat, lon float64, radius float64) ([]models.NearbyDriver, error) {
	args := m.Called(ctx, lat, lon, radius)
	return args.Get(0).([]models.NearbyDriver), args.Error(1)
}

func TestLocationController_UpdateLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockLocationService)
	ctrl := NewLocationController(mockSvc)

	r := gin.Default()
	r.POST("/api/v1/locations/update", ctrl.UpdateLocation)

	reqBody := map[string]interface{}{
		"driver_id": 1,
		"latitude":  -6.2,
		"longitude": 106.8,
	}
	jsonBody, _ := json.Marshal(reqBody)

	mockSvc.On("UpdateLocation", mock.Anything, uint(1), -6.2, 106.8).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/locations/update", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestLocationController_FindNearby(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockLocationService)
	ctrl := NewLocationController(mockSvc)

	r := gin.Default()
	r.GET("/api/v1/locations/nearby", ctrl.FindNearby)

	expectedDrivers := []models.NearbyDriver{
		{DriverID: 1, Latitude: -6.2, Longitude: 106.8, Distance: 500},
	}
	mockSvc.On("GetNearbyDrivers", mock.Anything, -6.2, 106.8, float64(5000)).Return(expectedDrivers, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/locations/nearby?lat=-6.2&lon=106.8&radius=5000", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Sukses", response["message"])
	mockSvc.AssertExpectations(t)
}

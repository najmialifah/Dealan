package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"matching-service/domain" // Sudah diganti ke domain
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMatchingService struct {
	mock.Mock
}

func (m *MockMatchingService) MatchOrder(ctx context.Context, req *domain.MatchRequest) (*domain.MatchedDriver, error) {
	args := m.Called(ctx, req)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.MatchedDriver), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestMatchingController_MatchDriver_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockMatchingService)
	ctrl := NewMatchingController(mockSvc)

	r := gin.Default()
	r.POST("/api/v1/match/", ctrl.MatchDriver)

	reqBody := domain.MatchRequest{
		OrderID:   1,
		Latitude:  -6.2,
		Longitude: 106.8,
	}
	jsonBody, _ := json.Marshal(reqBody)

	expectedDriver := &domain.MatchedDriver{
		DriverID:  "driver-100", // Diubah menjadi string sesuai kontrak domain
		Distance:  10,
	}

	mockSvc.On("MatchOrder", mock.Anything, &reqBody).Return(expectedDriver, nil)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/match/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Berhasil menemukan driver terdekat", response["message"])
	mockSvc.AssertExpectations(t)
}

func TestMatchingController_MatchDriver_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockMatchingService)
	ctrl := NewMatchingController(mockSvc)

	r := gin.Default()
	r.POST("/api/v1/match/", ctrl.MatchDriver)

	reqBody := domain.MatchRequest{
		OrderID:   1,
		Latitude:  -6.2,
		Longitude: 106.8,
	}
	jsonBody, _ := json.Marshal(reqBody)

	mockSvc.On("MatchOrder", mock.Anything, &reqBody).Return((*domain.MatchedDriver)(nil), errors.New("tidak ada driver yang ditemukan di sekitar lokasi"))

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/match/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}
package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"order-service/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderService adalah mock untuk OrderService
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	if args.Error(0) == nil {
		order.ID = 1
		order.Status = "PENDING"
	}
	return args.Error(0)
}

func TestOrderController_CreateOrder_Success(t *testing.T) {
	// Setup Gin menjadi test mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockOrderService)
	controller := NewOrderController(mockService)

	router := gin.Default()
	router.POST("/api/v1/orders/", controller.CreateOrder)

	// Persiapkan request body
	requestBody := map[string]interface{}{
		"user_id": 123,
		"detail_paket": map[string]interface{}{
			"barang": "Buku",
		},
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Ekspektasi: service dipanggil
	mockService.On("CreateOrder", mock.Anything, mock.AnythingOfType("*models.Order")).Return(nil)

	// Buat HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Rekam response
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	// Verifikasi isi response
	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Pesanan berhasil dibuat", response["message"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(123), data["user_id"])
	assert.Equal(t, "PENDING", data["status"])

	mockService.AssertExpectations(t)
}

func TestOrderController_CreateOrder_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockOrderService)
	controller := NewOrderController(mockService)

	router := gin.Default()
	router.POST("/api/v1/orders/", controller.CreateOrder)

	// Request body kosong (tanpa user_id)
	requestBody := map[string]interface{}{
		"detail_paket": map[string]interface{}{},
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Seharusnya BadRequest karena user_id tidak ada
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

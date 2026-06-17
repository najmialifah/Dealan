package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/najmialifah/Dealan/payment-service/controller"
	"github.com/najmialifah/Dealan/payment-service/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPaymentService mendefinisikan mock untuk service.PaymentService.
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) Process(ctx context.Context, req domain.PaymentRequest) (domain.PaymentResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(domain.PaymentResponse), args.Error(1)
}

func (m *MockPaymentService) ProcessWebhook(ctx context.Context, trxID string, status string) error {
	args := m.Called(ctx, trxID, status)
	return args.Error(0)
}

func (m *MockPaymentService) GetStatus(ctx context.Context, transactionID string) (domain.PaymentResponse, error) {
	args := m.Called(ctx, transactionID)
	return args.Get(0).(domain.PaymentResponse), args.Error(1)
}

func (m *MockPaymentService) GetDriverWallet(ctx context.Context, driverID string) (*domain.DriverWallet, error) {
	args := m.Called(ctx, driverID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DriverWallet), args.Error(1)
}

func setupGinRouter(ctrl *controller.PaymentController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	paymentGroup := router.Group("/payments")
	{
		paymentGroup.POST("/create", ctrl.Create)
		paymentGroup.POST("/webhook", ctrl.Webhook)
		paymentGroup.GET("/:transaction_id", ctrl.GetStatus)
		paymentGroup.GET("/driver/:driver_id/wallet", ctrl.GetDriverWallet)
	}
	return router
}

func TestPaymentController_Create(t *testing.T) {
	mockSrv := new(MockPaymentService)
	ctrl := controller.NewPaymentController(mockSrv)
	router := setupGinRouter(ctrl)

	t.Run("✅ StatusCreated jika request valid", func(t *testing.T) {
		req := domain.PaymentRequest{
			OrderID:          "ORD-123",
			Nominal:          25000,
			MetodePembayaran: "qris",
			UserID:           "b58d0426-8c46-4c7c-b391-a185bbf90f05",
			DriverID:         "cfc846bf-efb0-4dbf-8181-79ef88be4fbf",
		}
		expectedRes := domain.PaymentResponse{
			TransactionID: "TRX-MOCK-123",
			PaymentStatus: "PENDING",
			InvoiceURL:    "https://checkout.dealan.id/TRX-MOCK-123",
		}

		mockSrv.On("Process", mock.Anything, req).Return(expectedRes, nil).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/payments/create", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "TRX-MOCK-123")
		mockSrv.AssertExpectations(t)
	})

	t.Run("❌ StatusBadRequest jika JSON tidak valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/payments/create", bytes.NewBufferString(`{"order_id": ""}`))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPaymentController_GetStatus(t *testing.T) {
	mockSrv := new(MockPaymentService)
	ctrl := controller.NewPaymentController(mockSrv)
	router := setupGinRouter(ctrl)

	t.Run("✅ StatusOK jika transaksi ditemukan", func(t *testing.T) {
		expectedRes := domain.PaymentResponse{
			TransactionID: "TRX-123",
			PaymentStatus: "success",
			InvoiceURL:    "https://checkout.dealan.id/TRX-123",
		}
		mockSrv.On("GetStatus", mock.Anything, "TRX-123").Return(expectedRes, nil).Once()

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/payments/TRX-123", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
		mockSrv.AssertExpectations(t)
	})
}

func TestPaymentController_Webhook(t *testing.T) {
	mockSrv := new(MockPaymentService)
	ctrl := controller.NewPaymentController(mockSrv)
	router := setupGinRouter(ctrl)

	t.Run("✅ StatusOK saat memproses webhook valid", func(t *testing.T) {
		payload := map[string]string{
			"transaction_id": "TRX-123",
			"status":         "success",
		}
		mockSrv.On("ProcessWebhook", mock.Anything, "TRX-123", "success").Return(nil).Once()

		body, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/payments/webhook", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Webhook berhasil diproses")
		mockSrv.AssertExpectations(t)
	})
}

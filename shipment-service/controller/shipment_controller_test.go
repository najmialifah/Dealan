package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/najmialifah/Dealan/shipment-service/controller"
	"github.com/najmialifah/Dealan/shipment-service/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockShipmentService mendefinisikan mock untuk service.ShipmentService.
type MockShipmentService struct {
	mock.Mock
}

func (m *MockShipmentService) CreateShipment(ctx context.Context, req domain.ShipmentRequest) (domain.ShipmentResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(domain.ShipmentResponse), args.Error(1)
}

func (m *MockShipmentService) GetShipmentByID(ctx context.Context, id string) (*domain.Shipment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Shipment), args.Error(1)
}

func (m *MockShipmentService) GetShipmentByTrackingCode(ctx context.Context, code string) (*domain.Shipment, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Shipment), args.Error(1)
}

func (m *MockShipmentService) UpdateShipment(ctx context.Context, id string, req domain.ShipmentRequest) (*domain.Shipment, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Shipment), args.Error(1)
}

func (m *MockShipmentService) DeleteShipment(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockShipmentService) ListShipments(ctx context.Context) ([]domain.Shipment, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Shipment), args.Error(1)
}

func (m *MockShipmentService) UploadProof(ctx context.Context, id string, proof domain.ProofData) error {
	args := m.Called(ctx, id, proof)
	return args.Error(0)
}

func setupGinRouter(ctrl *controller.ShipmentController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	shipmentGroup := router.Group("/shipments")
	{
		shipmentGroup.POST("", ctrl.Create)
		shipmentGroup.GET("", ctrl.List)
		shipmentGroup.GET("/:id", ctrl.GetByID)
		shipmentGroup.GET("/track/:tracking_code", ctrl.GetByTrackingCode)
		shipmentGroup.PUT("/:id", ctrl.Update)
		shipmentGroup.DELETE("/:id", ctrl.Delete)
		shipmentGroup.POST("/:id/proof", ctrl.UploadProof)
	}
	return router
}

func TestShipmentController_Create(t *testing.T) {
	mockSrv := new(MockShipmentService)
	ctrl := controller.NewShipmentController(mockSrv)
	router := setupGinRouter(ctrl)

	t.Run("✅ StatusCreated jika payload pengiriman valid", func(t *testing.T) {
		req := domain.ShipmentRequest{
			OrderID:        "ORD-GOSEND-200",
			KategoriBarang: "makanan",
			BeratBarang:    2.0,
			NamaPenerima:   "Dewi",
			NomorPenerima:  "081223344",
		}
		expectedRes := domain.ShipmentResponse{
			ShipmentID:    "shp-uuid-200",
			KodeTracking:  "SHP-20240616-999",
			LabelShipping: "https://label.dealan.id/SHP-20240616-999",
		}

		mockSrv.On("CreateShipment", mock.Anything, req).Return(expectedRes, nil).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/shipments", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "shp-uuid-200")
		mockSrv.AssertExpectations(t)
	})
}

func TestShipmentController_GetByID(t *testing.T) {
	mockSrv := new(MockShipmentService)
	ctrl := controller.NewShipmentController(mockSrv)
	router := setupGinRouter(ctrl)

	t.Run("✅ StatusOK jika shipment ID ditemukan", func(t *testing.T) {
		expectedShp := &domain.Shipment{
			ID:           "shp-1",
			OrderID:      "ORD-1",
			TrackingCode: "SHP-1",
			Status:       "pending",
		}
		mockSrv.On("GetShipmentByID", mock.Anything, "shp-1").Return(expectedShp, nil).Once()

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/shipments/shp-1", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "shp-1")
		mockSrv.AssertExpectations(t)
	})
}

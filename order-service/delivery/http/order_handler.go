package http

import (
	"net/http"
	"order-service/domain"

	"github.com/gin-gonic/gin"
)

// OrderHandler menangani request HTTP untuk order
type OrderHandler struct {
	orderService domain.OrderService
}

// NewOrderHandler menginisialisasi OrderHandler
func NewOrderHandler(orderService domain.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder adalah handler untuk membuat pesanan baru
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req domain.CreateOrderRequest

	// Bind request body (JSON) ke struct CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid: " + err.Error()})
		return
	}

	// Buat domain object
	order := domain.Order{
		UserID:      req.UserID,
		DetailPaket: req.DetailPaket,
	}

	// Panggil service untuk memproses pembuatan pesanan
	if err := h.orderService.CreateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pesanan: " + err.Error()})
		return
	}

	// Postman expects an id or order_id
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Pesanan berhasil dibuat",
		"id":       order.ID, // So that it matches Postman check: json.id
		"order_id": order.ID,
		"data":     order,
	})
}

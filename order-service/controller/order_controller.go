package controller

import (
	"net/http"
	"order-service/models"
	"order-service/service"

	"github.com/gin-gonic/gin"
)

// OrderController menangani request HTTP untuk order
type OrderController struct {
	orderService service.OrderService
}

// NewOrderController menginisialisasi OrderController
func NewOrderController(orderService service.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// CreateOrder adalah handler untuk membuat pesanan baru
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var order models.Order

	// Bind request body (JSON) ke struct Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid: " + err.Error()})
		return
	}

	// Validasi input sederhana
	if order.UserID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id tidak boleh kosong"})
		return
	}

	// Panggil service untuk memproses pembuatan pesanan
	if err := c.orderService.CreateOrder(ctx.Request.Context(), &order); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pesanan: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Pesanan berhasil dibuat",
		"data":    order,
	})
}

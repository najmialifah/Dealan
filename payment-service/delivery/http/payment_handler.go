package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shakilaaulia/Dealan/payment-service/domain"
	"github.com/shakilaaulia/Dealan/payment-service/service"
	"gorm.io/gorm"
)

// PaymentHandler menangani request HTTP untuk transaksi pembayaran dan dompet driver.
type PaymentHandler struct {
	paymentService service.PaymentService
}

// NewPaymentHandler membuat instance baru dari PaymentHandler.
func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// Create menangani endpoint POST /payments/create
func (h *PaymentHandler) Create(c *gin.Context) {
	var req domain.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid: " + err.Error(),
		})
		return
	}

	// Tangkap key idempotensi dari header (jika ada) dan pasang ke request
	idempotencyHeader := c.GetHeader("X-Idempotency-Key")
	if idempotencyHeader != "" {
		req.IdempotencyKey = idempotencyHeader
	}

	res, err := h.paymentService.Process(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memproses transaksi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   res,
	})
}

// GetStatus menangani endpoint GET /payments/:transaction_id
func (h *PaymentHandler) GetStatus(c *gin.Context) {
	trxID := c.Param("transaction_id")
	if trxID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Transaction ID wajib disertakan",
		})
		return
	}

	res, err := h.paymentService.GetStatus(c.Request.Context(), trxID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Transaksi tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data transaksi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   res,
	})
}

// Webhook menangani endpoint POST /payments/webhook (simulasi webhook dari Payment Gateway)
func (h *PaymentHandler) Webhook(c *gin.Context) {
	var payload struct {
		TransactionID string `json:"transaction_id" binding:"required"`
		Status        string `json:"status" binding:"required,oneof=success failed"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload webhook tidak valid: " + err.Error(),
		})
		return
	}

	err := h.paymentService.ProcessWebhook(c.Request.Context(), payload.TransactionID, payload.Status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Transaksi tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memproses webhook: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Webhook berhasil diproses",
	})
}

// GetDriverWallet menangani endpoint GET /payments/driver/:driver_id/wallet
func (h *PaymentHandler) GetDriverWallet(c *gin.Context) {
	driverID := c.Param("driver_id")
	if driverID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Driver ID wajib disertakan",
		})
		return
	}

	wallet, err := h.paymentService.GetDriverWallet(c.Request.Context(), driverID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Dompet driver belum aktif atau tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil dompet driver: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   wallet,
	})
}

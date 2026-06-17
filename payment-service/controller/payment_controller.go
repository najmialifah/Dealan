package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/najmialifah/Dealan/payment-service/domain"
	"github.com/najmialifah/Dealan/payment-service/service"
	"gorm.io/gorm"
)

// PaymentController menangani request HTTP untuk transaksi pembayaran dan dompet driver.
type PaymentController struct {
	paymentService service.PaymentService
}

// NewPaymentController membuat instance baru dari PaymentController.
func NewPaymentController(paymentService service.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// Create menangani endpoint POST /payments/create
func (ctrl *PaymentController) Create(c *gin.Context) {
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

	res, err := ctrl.paymentService.Process(c.Request.Context(), req)
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
func (ctrl *PaymentController) GetStatus(c *gin.Context) {
	trxID := c.Param("transaction_id")
	if trxID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Transaction ID wajib disertakan",
		})
		return
	}

	res, err := ctrl.paymentService.GetStatus(c.Request.Context(), trxID)
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
func (ctrl *PaymentController) Webhook(c *gin.Context) {
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

	err := ctrl.paymentService.ProcessWebhook(c.Request.Context(), payload.TransactionID, payload.Status)
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
func (ctrl *PaymentController) GetDriverWallet(c *gin.Context) {
	driverID := c.Param("driver_id")
	if driverID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Driver ID wajib disertakan",
		})
		return
	}

	wallet, err := ctrl.paymentService.GetDriverWallet(c.Request.Context(), driverID)
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

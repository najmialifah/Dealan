package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/najmialifah/Dealan/shipment-service/domain"
	"github.com/najmialifah/Dealan/shipment-service/service"
	"gorm.io/gorm"
)

// ShipmentController menangani HTTP request untuk CRUD data pengiriman barang (Shipment).
type ShipmentController struct {
	shipmentService service.ShipmentService
}

// NewShipmentController membuat instance baru dari ShipmentController.
func NewShipmentController(shipmentService service.ShipmentService) *ShipmentController {
	return &ShipmentController{
		shipmentService: shipmentService,
	}
}

// Create menangani endpoint POST /shipments
func (ctrl *ShipmentController) Create(c *gin.Context) {
	var req domain.ShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid: " + err.Error(),
		})
		return
	}

	res, err := ctrl.shipmentService.CreateShipment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat pengiriman: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   res,
	})
}

// GetByID menangani endpoint GET /shipments/:id
func (ctrl *ShipmentController) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Shipment ID wajib disertakan",
		})
		return
	}

	shipment, err := ctrl.shipmentService.GetShipmentByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Pengiriman tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data pengiriman: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   shipment,
	})
}

// GetByTrackingCode menangani endpoint GET /shipments/track/:tracking_code
func (ctrl *ShipmentController) GetByTrackingCode(c *gin.Context) {
	code := c.Param("tracking_code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kode tracking wajib disertakan",
		})
		return
	}

	shipment, err := ctrl.shipmentService.GetShipmentByTrackingCode(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Pengiriman dengan kode tracking tersebut tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal melacak pengiriman: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   shipment,
	})
}

// Update menangani endpoint PUT /shipments/:id
func (ctrl *ShipmentController) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Shipment ID wajib disertakan",
		})
		return
	}

	var req domain.ShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid: " + err.Error(),
		})
		return
	}

	shipment, err := ctrl.shipmentService.UpdateShipment(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Pengiriman tidak ditemukan untuk diperbarui",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui pengiriman: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   shipment,
	})
}

// Delete menangani endpoint DELETE /shipments/:id
func (ctrl *ShipmentController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Shipment ID wajib disertakan",
		})
		return
	}

	err := ctrl.shipmentService.DeleteShipment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus pengiriman: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pengiriman berhasil dihapus",
	})
}

// List menangani endpoint GET /shipments
func (ctrl *ShipmentController) List(c *gin.Context) {
	shipments, err := ctrl.shipmentService.ListShipments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil daftar pengiriman: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   shipments,
	})
}

// UploadProof menangani endpoint POST /shipments/:id/proof
func (ctrl *ShipmentController) UploadProof(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Shipment ID wajib disertakan",
		})
		return
	}

	var proof domain.ProofData
	if err := c.ShouldBindJSON(&proof); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload bukti pengiriman tidak valid: " + err.Error(),
		})
		return
	}

	err := ctrl.shipmentService.UploadProof(c.Request.Context(), id, proof)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Pengiriman tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengunggah bukti serah terima: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Bukti serah terima berhasil diunggah, pengiriman ditandai selesai (delivered)",
	})
}

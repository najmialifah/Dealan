package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/najmialifah/Dealan/driver-service/domain"
	"github.com/najmialifah/Dealan/driver-service/service"
)

// DriverHandler adalah controller REST API untuk driver-service menggunakan Gin
type DriverHandler struct {
	svc service.DriverService
}

// NewDriverHandler membuat instance baru dari DriverHandler
func NewDriverHandler(svc service.DriverService) *DriverHandler {
	return &DriverHandler{svc: svc}
}

// UpdateLocation menangani pembaruan koordinat GPS driver lewat X-Driver-ID header
func (h *DriverHandler) UpdateLocation(c *gin.Context) {
	id := c.GetHeader("X-Driver-ID")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Driver-ID header tidak ditemukan"})
		return
	}

	var req domain.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateLocation(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "lokasi berhasil diperbarui"})
}

// UpdateStatus menangani pembaruan status online dan jenis layanan aktif driver
func (h *DriverHandler) UpdateStatus(c *gin.Context) {
	id := c.GetHeader("X-Driver-ID")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Driver-ID header tidak ditemukan"})
		return
	}

	var req domain.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateStatus(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status berhasil diperbarui"})
}

// GetProfile mengambil data profil lengkap pengemudi beserta kendaraan dan statusnya
func (h *DriverHandler) GetProfile(c *gin.Context) {
	id := c.GetHeader("X-Driver-ID")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Driver-ID header tidak ditemukan"})
		return
	}

	profile, err := h.svc.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// AddVehicle menambahkan mobil/motor baru milik driver
func (h *DriverHandler) AddVehicle(c *gin.Context) {
	id := c.GetHeader("X-Driver-ID")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Driver-ID header tidak ditemukan"})
		return
	}

	var req domain.Vehicle
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.AddVehicle(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "kendaraan berhasil ditambahkan"})
}

package controller

import (
	"location-service/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LocationController struct {
	locationService service.LocationService
}

func NewLocationController(svc service.LocationService) *LocationController {
	return &LocationController{locationService: svc}
}

// UpdateLocationRequest struktur request JSON
type UpdateLocationRequest struct {
	DriverID  uint    `json:"driver_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// UpdateLocation handler untuk menerima push location dari app driver
func (c *LocationController) UpdateLocation(ctx *gin.Context) {
	var req UpdateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah: " + err.Error()})
		return
	}

	// Memanggil service untuk update secara asinkron
	if err := c.locationService.UpdateLocation(ctx.Request.Context(), req.DriverID, req.Latitude, req.Longitude); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update lokasi"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Lokasi berhasil diantrekan untuk update"})
}

// FindNearby handler untuk mencari driver terdekat dari user
func (c *LocationController) FindNearby(ctx *gin.Context) {
	// Ambil parameter query lat, lon, radius
	latStr := ctx.Query("lat")
	lonStr := ctx.Query("lon")
	radiusStr := ctx.DefaultQuery("radius", "5000") // default radius 5km

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lon, errLon := strconv.ParseFloat(lonStr, 64)
	radius, errRad := strconv.ParseFloat(radiusStr, 64)

	if errLat != nil || errLon != nil || errRad != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parameter lat, lon, atau radius tidak valid"})
		return
	}

	drivers, err := c.locationService.GetNearbyDrivers(ctx.Request.Context(), lat, lon, radius)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari driver: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    drivers,
	})
}

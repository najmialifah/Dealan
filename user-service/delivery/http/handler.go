package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shakilaaulia/Dealan/user-service/domain"
	"github.com/shakilaaulia/Dealan/user-service/service"
)

// UserHandler adalah controller REST API untuk user-service menggunakan Gin
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler membuat instance baru dari UserHandler
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetProfile menangani pengambilan profil user lewat X-User-ID header
func (h *UserHandler) GetProfile(c *gin.Context) {
	id := c.GetHeader("X-User-ID")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header tidak ditemukan"})
		return
	}

	profile, err := h.svc.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile memperbarui profil user lewat X-User-ID header
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id := c.GetHeader("X-User-ID")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header tidak ditemukan"})
		return
	}

	var req domain.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateProfile(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profil berhasil diperbarui"})
}

// GetInternalName menangani API internal untuk query nama user
func (h *UserHandler) GetInternalName(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter tidak ditemukan"})
		return
	}

	name, err := h.svc.GetInternalName(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": name})
}

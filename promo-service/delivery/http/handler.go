package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"promo-service/domain"
	"promo-service/service"
)

type PromoHandler struct {
	Service service.PromoService
}

// NewPromoHandler membuat instance baru dari handler promo
func NewPromoHandler(s service.PromoService) *PromoHandler {
	return &PromoHandler{s}
}

// ApplyPromo menangani request POST untuk menerapkan kode promo ke pesanan
func (h *PromoHandler) ApplyPromo(c *gin.Context) {
	var req domain.PromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.ApplyPromo(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// SetupRoutes mendaftarkan rute API untuk promo service
func SetupRoutes(r *gin.Engine, handler *PromoHandler) {
	r.POST("/promo", handler.ApplyPromo)
}
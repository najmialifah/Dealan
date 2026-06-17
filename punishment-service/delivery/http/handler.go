package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/najmialifah/Dealan/punishment-service/domain"
	"github.com/najmialifah/Dealan/punishment-service/service"
)

type PunishmentHandler struct {
	svc service.PunishmentService
}

// NewPunishmentHandler membuat instance baru dari handler punishment
func NewPunishmentHandler(svc service.PunishmentService) *PunishmentHandler {
	return &PunishmentHandler{svc: svc}
}

// Apply menangani request POST untuk menerapkan sanksi kepada akun
func (h *PunishmentHandler) Apply(c *gin.Context) {
	var req domain.PunishmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.ApplyPunishment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// SetupRoutes mendefinisikan rute API untuk punishment service
func SetupRoutes(r *gin.Engine, handler *PunishmentHandler) {
	r.POST("/punishment/apply", handler.Apply)
}
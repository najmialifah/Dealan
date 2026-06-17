package http

import (
	"matching-service/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MatchingHandler struct {
	svc domain.MatchingService
}

func NewMatchingHandler(svc domain.MatchingService) *MatchingHandler {
	return &MatchingHandler{svc: svc}
}

func (h *MatchingHandler) MatchDriver(ctx *gin.Context) {
	var req domain.MatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid: " + err.Error()})
		return
	}

	driver, err := h.svc.MatchOrder(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == "tidak ada driver yang ditemukan di sekitar lokasi" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan internal: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menemukan driver terdekat",
		"data":    driver,
	})
}

package controller // WAJIB package controller

import (
	"matching-service/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MatchingController struct {
	svc domain.MatchingService
}

func NewMatchingController(svc domain.MatchingService) *MatchingController {
	return &MatchingController{svc: svc}
}

func (c *MatchingController) MatchDriver(ctx *gin.Context) {
	var req domain.MatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid: " + err.Error()})
		return
	}

	driver, err := c.svc.MatchOrder(ctx.Request.Context(), &req)
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
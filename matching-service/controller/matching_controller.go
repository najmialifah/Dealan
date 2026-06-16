package controller

import (
	"matching-service/models"
	"matching-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MatchingController struct {
	svc service.MatchingService
}

func NewMatchingController(svc service.MatchingService) *MatchingController {
	return &MatchingController{svc: svc}
}

// MatchDriver adalah handler HTTP untuk menerima permintaan mencarikan driver untuk order tertentu
func (c *MatchingController) MatchDriver(ctx *gin.Context) {
	var req models.MatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid: " + err.Error()})
		return
	}

	driver, err := c.svc.MatchOrder(ctx.Request.Context(), &req)
	if err != nil {
		// Asumsi jika pesan mengandung "tidak ada", kita kembalikan NotFound
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

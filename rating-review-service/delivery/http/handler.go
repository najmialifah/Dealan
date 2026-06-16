package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shakilaaulia/Dealan/rating-review-service/domain"
	"github.com/shakilaaulia/Dealan/rating-review-service/service"
)

type RatingHandler struct {
	svc service.RatingService
}

// NewRatingHandler membuat instance baru dari handler rating
func NewRatingHandler(svc service.RatingService) *RatingHandler {
	return &RatingHandler{svc: svc}
}

// Submit menangani request POST untuk submit rating dan review baru
func (h *RatingHandler) Submit(c *gin.Context) {
	var req domain.RatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.SubmitReview(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// SetupRoutes mendaftarkan rute-rute API untuk rating service
func SetupRoutes(r *gin.Engine, handler *RatingHandler) {
	r.POST("/rating/submit", handler.Submit)
}
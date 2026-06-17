package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"map-route-service/domain"
)

type MapHandler struct {
	mapService domain.MapService
}

func NewMapHandler(svc domain.MapService) *MapHandler {
	return &MapHandler{mapService: svc}
}

// GetRoute melayani pengambilan rute perjalanan.
// Mendukung request POST dengan body JSON, maupun GET dengan query parameters.
func (h *MapHandler) GetRoute(c *gin.Context) {
	var req domain.RouteRequest

	if c.Request.Method == http.MethodPost {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah: " + err.Error()})
			return
		}
	} else {
		req.Origin = c.Query("origin")
		req.Destination = c.Query("destination")
		if req.Origin == "" || req.Destination == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'origin' dan 'destination' wajib diisi"})
			return
		}
	}

	res, err := h.mapService.GetOrCreateRoute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

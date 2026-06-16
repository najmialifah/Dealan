package routes

import (
	"matching-service/controller"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, ctrl *controller.MatchingController) {
	api := router.Group("/api/v1")
	{
		// Langsung definisikan rute /match di dalam /api/v1
		api.POST("/match", ctrl.MatchDriver)
	}
}
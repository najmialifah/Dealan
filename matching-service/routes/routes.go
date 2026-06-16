package routes

import (
	"matching-service/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, ctrl *controller.MatchingController) {
	api := router.Group("/api/v1")
	{
		matchingRoutes := api.Group("/match")
		{
			// Endpoint untuk meminta matching driver dengan suatu order secara langsung
			matchingRoutes.POST("/", ctrl.MatchDriver)
		}
	}
}

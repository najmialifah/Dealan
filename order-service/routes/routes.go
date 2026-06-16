package routes

import (
	"order-service/controller"

	"github.com/gin-gonic/gin"
)

// SetupRoutes mendefinisikan rute-rute API untuk order-service
func SetupRoutes(router *gin.Engine, orderController *controller.OrderController) {
	api := router.Group("/api/v1")
	{
		orderRoutes := api.Group("/orders")
		{
			// Endpoint untuk membuat order baru
			orderRoutes.POST("/", orderController.CreateOrder)
		}
	}
}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/najmialifah/Dealan/shipment-service/controller"
)

// SetupRoutes mendaftarkan endpoint Gin Gonic untuk Shipment Service.
func SetupRoutes(router *gin.Engine, ctrl *controller.ShipmentController) {
	// Endpoint dasar untuk memantau kesehatan service
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Shipment Service is Up and Running!")
	})

	// Grouping endpoint CRUD sesuai standar RESTful API
	shipmentGroup := router.Group("/shipments")
	{
		shipmentGroup.POST("", ctrl.Create)
		shipmentGroup.GET("", ctrl.List)
		shipmentGroup.GET("/:id", ctrl.GetByID)
		shipmentGroup.GET("/track/:tracking_code", ctrl.GetByTrackingCode)
		shipmentGroup.PUT("/:id", ctrl.Update)
		shipmentGroup.DELETE("/:id", ctrl.Delete)
		shipmentGroup.POST("/:id/proof", ctrl.UploadProof)
	}
}

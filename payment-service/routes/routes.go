package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/najmialifah/Dealan/payment-service/controller"
)

// SetupRoutes mendefinisikan seluruh endpoint HTTP untuk Payment Service menggunakan Gin Gonic.
func SetupRoutes(router *gin.Engine, ctrl *controller.PaymentController) {
	// Endpoint dasar untuk memantau kesehatan service
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Payment Service is Up and Running!")
	})

	// Grouping endpoint sesuai spesifikasi PRD kelompok
	paymentGroup := router.Group("/payments")
	{
		paymentGroup.POST("/create", ctrl.Create)
		paymentGroup.POST("/webhook", ctrl.Webhook)
		paymentGroup.GET("/:transaction_id", ctrl.GetStatus)
		paymentGroup.GET("/driver/:driver_id/wallet", ctrl.GetDriverWallet)
	}
}

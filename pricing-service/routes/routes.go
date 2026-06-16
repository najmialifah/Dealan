package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shakilaaulia/Dealan/pricing-service/controller"
)

// SetupRoutes mendefinisikan seluruh mapping URI endpoint HTTP untuk pricing-service.
func SetupRoutes(r *gin.Engine, ctrl *controller.PricingController) {
	pricing := r.Group("/pricing")
	{
		pricing.POST("/estimate", ctrl.CalculateEstimate)
		pricing.POST("/negotiate", ctrl.NegotiatePrice)
		pricing.GET("/rules/:service_type", ctrl.GetRules)
	}
}

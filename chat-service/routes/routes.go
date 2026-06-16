package routes

import (
	"github.com/gin-gonic/gin"
	"chat-service/controller"
)

// SetupRoutes mendefinisikan pemetaan endpoint HTTP dan WebSocket untuk chat-service.
func SetupRoutes(r *gin.Engine, ctrl *controller.ChatController) {
	// Endpoint koneksi WebSocket real-time
	r.GET("/chat/ws", ctrl.UpgradeWS)

	// REST API untuk mendapatkan riwayat pesan dan inisialisasi room chat
	r.GET("/chat/:order_id/history", ctrl.GetHistory)
	r.POST("/chat/room", ctrl.CreateRoom)
}

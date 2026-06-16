package controller

import (
	"net/http"

	"chat-service/models"
	"chat-service/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Mengizinkan semua origin untuk development. Ganti dengan origin yang tepat pada production.
		return true
	},
}

// ChatController mengelola request masuk baik HTTP biasa maupun koneksi WebSocket.
type ChatController struct {
	chatService service.ChatService
}

// NewChatController menginisiasi ChatController baru.
func NewChatController(svc service.ChatService) *ChatController {
	return &ChatController{chatService: svc}
}

// UpgradeWS meningkatkan protokol koneksi HTTP biasa menjadi WebSocket untuk komunikasi real-time.
func (ctrl *ChatController) UpgradeWS(c *gin.Context) {
	orderID := c.Query("order_id")
	userID := c.Query("user_id")
	role := c.Query("role")

	if orderID == "" || userID == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'order_id', 'user_id', dan 'role' wajib diisi"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meningkatkan koneksi ke WebSocket: " + err.Error()})
		return
	}

	client := &service.Client{
		UserID:   userID,
		Role:     role,
		OrderID:  orderID,
		Conn:     conn,
		Send:     make(chan models.WSMessage, 256),
		Hub:      ctrl.chatService,
	}

	ctrl.chatService.Register(client)

	// Jalankan handler pembacaan dan penulisan data websocket secara paralel dalam thread Goroutine
	go client.WritePump()
	go client.ReadPump()
}

// GetHistory mengembalikan riwayat obrolan lengkap berdasarkan ID pesanan (order_id).
func (ctrl *ChatController) GetHistory(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pesanan wajib diisi"})
		return
	}

	history, err := ctrl.chatService.GetHistory(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat chat: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// CreateRoom membuat ruang obrolan baru untuk pesanan tertentu secara manual.
func (ctrl *ChatController) CreateRoom(c *gin.Context) {
	var req struct {
		OrderID  string `json:"order_id" binding:"required"`
		UserID   string `json:"user_id" binding:"required"`
		DriverID string `json:"driver_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah: " + err.Error()})
		return
	}

	err := ctrl.chatService.CreateRoom(c.Request.Context(), req.OrderID, req.UserID, req.DriverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat room chat: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room chat berhasil dibuat", "order_id": req.OrderID})
}

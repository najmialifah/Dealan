package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/najmialifah/Dealan/notification-service/domain"
	"github.com/najmialifah/Dealan/notification-service/service"
)

type NotificationHandler struct {
	svc service.NotificationService
}

// NewNotificationHandler membuat instance baru dari handler HTTP notifikasi
func NewNotificationHandler(svc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// Send menangani request POST untuk mengirim notifikasi secara manual
func (h *NotificationHandler) Send(c *gin.Context) {
	var req domain.NotificationRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.SendNotification(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SetupRoutes mendefinisikan rute endpoint untuk Notification Service
func SetupRoutes(r *gin.Engine, handler *NotificationHandler) {
	r.POST("/notification/send", handler.Send)
}
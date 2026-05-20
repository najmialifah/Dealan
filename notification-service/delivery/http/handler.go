package http

import (
	"encoding/json"
	"net/http"
	
	"github.com/shakilaaulia/Dealan/notification-service/domain"
	"github.com/shakilaaulia/Dealan/notification-service/service"
)

type NotificationHandler struct {
	svc service.NotificationService // Handler memanggil Service
}

func NewNotificationHandler(svc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) Send(w http.ResponseWriter, r *http.Request) {
	var req domain.NotificationRequest
	
	// 1. Decode data JSON dari request body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Panggil fungsi di layer Service
	resp, err := h.svc.SendNotification(r.Context(), req)
	if err != nil {
		// Jika service error, kirim status 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Kirim balik hasilnya dalam bentuk JSON
	json.NewEncoder(w).Encode(resp)
}
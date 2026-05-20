package http

import (
	"encoding/json"
	"net/http"
	"github.com/shakilaaulia/Dealan/punishment-service/domain"
	"github.com/shakilaaulia/Dealan/punishment-service/service"
)

type PunishmentHandler struct {
	svc service.PunishmentService
}

func NewPunishmentHandler(svc service.PunishmentService) *PunishmentHandler {
	return &PunishmentHandler{svc: svc}
}

func (h *PunishmentHandler) Apply(w http.ResponseWriter, r *http.Request) {
	var req domain.PunishmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.svc.ApplyPunishment(r.Context(), req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
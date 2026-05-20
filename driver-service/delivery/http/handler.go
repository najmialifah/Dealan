package http

import (
	"encoding/json"
	"net/http"

	"github.com/shakilaaulia/Dealan/driver-service/domain"
	"github.com/shakilaaulia/Dealan/driver-service/service"
)

type DriverHandler struct {
	svc service.DriverService
}

func NewDriverHandler(svc service.DriverService) *DriverHandler {
	return &DriverHandler{svc: svc}
}

func (h *DriverHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("X-Driver-ID")
	var req domain.UpdateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.svc.UpdateLocation(r.Context(), id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func (h *DriverHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("X-Driver-ID")
	var req domain.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.svc.UpdateStatus(r.Context(), id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func (h *DriverHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("X-Driver-ID")
	
	profile, err := h.svc.GetProfile(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(profile)
}

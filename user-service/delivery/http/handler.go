package http

import (
	"encoding/json"
	"net/http"

	"github.com/shakilaaulia/Dealan/user-service/domain"
	"github.com/shakilaaulia/Dealan/user-service/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// In reality, id comes from middleware/token
	id := r.Header.Get("X-User-ID")
	
	profile, err := h.svc.GetProfile(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(profile)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("X-User-ID")
	var req domain.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.svc.UpdateProfile(r.Context(), id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetInternalName(w http.ResponseWriter, r *http.Request) {
	// Assuming ID comes from query string or URL path (e.g. /users/internal?id=123)
	id := r.URL.Query().Get("id")
	
	name, err := h.svc.GetInternalName(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]string{"name": name})
}

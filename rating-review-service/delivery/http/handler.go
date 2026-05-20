package http

import (
	"encoding/json"
	"net/http"
	"github.com/shakilaaulia/Dealan/rating-review-service/domain"
	"github.com/shakilaaulia/Dealan/rating-review-service/service"
)

type RatingHandler struct {
	svc service.RatingService
}

func NewRatingHandler(svc service.RatingService) *RatingHandler {
	return &RatingHandler{svc: svc}
}

func (h *RatingHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req domain.RatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.svc.SubmitReview(r.Context(), req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
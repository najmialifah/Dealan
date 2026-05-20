package http

import (
	"encoding/json"
	"net/http"

	"map-route-service/domain"
	"map-route-service/service"
)

type MapHandler struct {
	Service service.MapService
}

func NewMapHandler(s service.MapService) *MapHandler {
	return &MapHandler{s}
}

func (h *MapHandler) GetRoute(w http.ResponseWriter, r *http.Request) {

	req := domain.RouteRequest{
		Origin:      "A",
		Destination: "B",
	}

	res, err := h.Service.GetRoute(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}
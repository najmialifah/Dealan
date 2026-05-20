package http

import (
	"encoding/json"
	"net/http"

	"promo-service/domain"
	"promo-service/service"
)

type PromoHandler struct {
	Service service.PromoService
}

func NewPromoHandler(s service.PromoService) *PromoHandler {
	return &PromoHandler{s}
}

func (h *PromoHandler) ApplyPromo(w http.ResponseWriter, r *http.Request) {

	req := domain.PromoRequest{
		Code: "DISKON",
	}

	res, _ := h.Service.ApplyPromo(r.Context(), req)

	json.NewEncoder(w).Encode(res)
}
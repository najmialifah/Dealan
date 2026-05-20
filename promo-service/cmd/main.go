package main

import (
	"log"
	"net/http"

	deliveryHttp "promo-service/delivery/http"
	"promo-service/service"
)

func main() {

	svc := service.NewPromoService(nil)
	handler := deliveryHttp.NewPromoHandler(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("/promo", handler.ApplyPromo)

	log.Println("Promo Service running on :8086")

	http.ListenAndServe(":8086", mux)
}
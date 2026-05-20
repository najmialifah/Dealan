package main

import (
	"log"
	"net/http"

	deliveryHttp "map-route-service/delivery/http"
	"map-route-service/service"
)

func main() {

	// sementara repo = nil (karena belum implement DB)
	mapSvc := service.NewMapService(nil)

	mapHandler := deliveryHttp.NewMapHandler(mapSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("/route", mapHandler.GetRoute)

	log.Println("Map Route Service running on :8085")

	if err := http.ListenAndServe(":8085", mux); err != nil {
		log.Fatal(err)
	}
}
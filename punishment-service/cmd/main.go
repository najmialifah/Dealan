package main

import (
	"log"
	"net/http"

	deliveryHttp "github.com/shakilaaulia/Dealan/punishment-service/delivery/http"
	"github.com/shakilaaulia/Dealan/punishment-service/service"
)

func main() {
	svc := service.NewPunishmentService(nil)
	handler := deliveryHttp.NewPunishmentHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/punishment/apply", handler.Apply)

	log.Println("Punishment Service listening on :8086")
	if err := http.ListenAndServe(":8086", mux); err != nil {
		log.Fatal(err)
	}
}
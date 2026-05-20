package main

import (
	"log"
	"net/http"

	deliveryHttp "github.com/shakilaaulia/Dealan/rating-review-service/delivery/http"
	"github.com/shakilaaulia/Dealan/rating-review-service/service"
)

func main() {
	// Initialize service (passing nil for repository draft)
	ratingSvc := service.NewRatingService(nil)
	ratingHandler := deliveryHttp.NewRatingHandler(ratingSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/rating/submit", ratingHandler.Submit)

	log.Println("Rating Service listening on :8085")
	if err := http.ListenAndServe(":8085", mux); err != nil {
		log.Fatal(err)
	}
}
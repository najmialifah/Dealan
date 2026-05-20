package main

import (
	"log"
	"net/http"

	deliveryHttp "github.com/shakilaaulia/Dealan/driver-service/delivery/http"
	"github.com/shakilaaulia/Dealan/driver-service/service"
	// "github.com/shakilaaulia/Dealan/driver-service/repository" // Import when concrete repository is implemented
)

func main() {
	// Initialize concrete repository (e.g., Postgres)
	// repo := repository.NewPostgresDriverRepository(...)
	
	// Initialize service (passing nil for draft)
	driverSvc := service.NewDriverService(nil)
	
	// Initialize handler
	driverHandler := deliveryHttp.NewDriverHandler(driverSvc)
	
	// Setup Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/drivers/location", driverHandler.UpdateLocation) // PUT
	mux.HandleFunc("/drivers/status", driverHandler.UpdateStatus)     // PATCH
	mux.HandleFunc("/drivers/profile", driverHandler.GetProfile)      // GET
	
	log.Println("Driver Service listening on :8083")
	if err := http.ListenAndServe(":8083", mux); err != nil {
		log.Fatal(err)
	}
}

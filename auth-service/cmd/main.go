package main

import (
	"log"
	"net/http"

	deliveryHttp "github.com/shakilaaulia/Dealan/auth-service/delivery/http"
	"github.com/shakilaaulia/Dealan/auth-service/service"
	// "github.com/shakilaaulia/Dealan/auth-service/repository" // Import when concrete repository is implemented
)

func main() {
	// Initialize concrete repository (e.g., Postgres)
	// repo := repository.NewPostgresAuthRepository(...)
	
	// Initialize service (passing nil for draft)
	authSvc := service.NewAuthService(nil)
	
	// Initialize handler
	authHandler := deliveryHttp.NewAuthHandler(authSvc)
	
	// Setup Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", authHandler.Register) // POST
	mux.HandleFunc("/auth/login", authHandler.Login)       // POST
	mux.HandleFunc("/auth/validate", authHandler.Validate) // POST
	
	log.Println("Auth Service listening on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}

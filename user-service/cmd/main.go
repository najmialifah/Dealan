package main

import (
	"log"
	"net/http"

	deliveryHttp "github.com/shakilaaulia/Dealan/user-service/delivery/http"
	"github.com/shakilaaulia/Dealan/user-service/service"
	// "github.com/shakilaaulia/Dealan/user-service/repository" // Import when concrete repository is implemented
)

func main() {
	// Initialize concrete repository (e.g., Postgres)
	// repo := repository.NewPostgresUserRepository(...)
	
	// Initialize service (passing nil for draft)
	userSvc := service.NewUserService(nil)
	
	// Initialize handler
	userHandler := deliveryHttp.NewUserHandler(userSvc)
	
	// Setup Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/users/profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetProfile(w, r)
		} else if r.Method == http.MethodPut {
			userHandler.UpdateProfile(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/users/internal", userHandler.GetInternalName) // GET
	
	log.Println("User Service listening on :8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}

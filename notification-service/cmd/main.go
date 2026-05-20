package main

import (
	"log"
	"net/http"

	deliveryHttp "github.com/shakilaaulia/Dealan/notification-service/delivery/http"
	"github.com/shakilaaulia/Dealan/notification-service/service"
)

func main() {
	// 1. Initialize service 
	notificationSvc := service.NewNotificationService(nil)
	
	// 2. Initialize handler (menghubungkan service ke HTTP)
	notificationHandler := deliveryHttp.NewNotificationHandler(notificationSvc)
	
	// 3. Setup Routes (Daftar alamat API)
	mux := http.NewServeMux()
	mux.HandleFunc("/notification/send", notificationHandler.Send) // Endpoint untuk kirim notif
	
	// 4. Jalankan Server di Port 8084
	log.Println("Notification Service listening on :8084")
	if err := http.ListenAndServe(":8084", mux); err != nil {
		log.Fatal(err)
	}
}
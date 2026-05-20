package main

import (
	"log"
	"net/http"

	deliveryHttp "chat-service/delivery/http"
	"chat-service/service"
)

func main() {

	svc := service.NewChatService(nil)
	handler := deliveryHttp.NewChatHandler(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("/chat", handler.Send)

	log.Println("Chat Service running on :8087")

	http.ListenAndServe(":8087", mux)
}
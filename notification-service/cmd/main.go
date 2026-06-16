package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	deliveryHttp "github.com/shakilaaulia/Dealan/notification-service/delivery/http"
	deliveryKafka "github.com/shakilaaulia/Dealan/notification-service/delivery/kafka"
	"github.com/shakilaaulia/Dealan/notification-service/repository"
	"github.com/shakilaaulia/Dealan/notification-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Inisialisasi Koneksi database PostgreSQL via GORM
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		// Menggunakan URL default lokal jika env var kosong
		dbURL = "postgres://dealan:dealan_secret@localhost:5432/dealan_db?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Notification Service] Gagal terhubung ke database PostgreSQL: %v", err)
	}
	log.Println("[Notification Service] Koneksi PostgreSQL (GORM) berhasil terjalin")

	// 2. Inisialisasi Repository dan Service
	repo := repository.NewPostgresRepository(db)
	notificationSvc := service.NewNotificationService(repo)

	// 3. Setup Kafka Consumer untuk mendengarkan event PAYMENT_COMPLETED
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := deliveryKafka.NewKafkaConsumer([]string{kafkaBrokers}, "payment.completed", "notification-group", notificationSvc)
	go consumer.Start(ctx)

	// 4. Setup Gin Router & HTTP Server
	r := gin.Default()
	notificationHandler := deliveryHttp.NewNotificationHandler(notificationSvc)
	deliveryHttp.SetupRoutes(r, notificationHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Menjalankan server HTTP di goroutine terpisah agar tidak memblokir shutdown signal
	go func() {
		log.Printf("[Notification Service] Berjalan di port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Notification Service] Gagal menjalankan HTTP server: %v", err)
		}
	}()

	// Graceful Shutdown Handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[Notification Service] Memulai proses shutdown...")

	cancel() // Mematikan Kafka Consumer reader

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[Notification Service] HTTP Server terpaksa dimatikan: %v", err)
	}

	log.Println("[Notification Service] Sukses dimatikan secara aman.")
}
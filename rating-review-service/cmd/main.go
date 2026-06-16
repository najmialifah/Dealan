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
	deliveryHttp "github.com/shakilaaulia/Dealan/rating-review-service/delivery/http"
	"github.com/shakilaaulia/Dealan/rating-review-service/repository"
	"github.com/shakilaaulia/Dealan/rating-review-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Inisialisasi Koneksi PostgreSQL via GORM
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://dealan:dealan_secret@localhost:5432/dealan_db?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Rating Service] Gagal terhubung ke database: %v", err)
	}
	log.Println("[Rating Service] Koneksi database PostgreSQL berhasil terjalin")

	// 2. Inisialisasi Repository dan Service
	repo := repository.NewPostgresRepository(db)
	ratingSvc := service.NewRatingService(repo)

	// 3. Setup Gin Router & Server HTTP
	r := gin.Default()
	ratingHandler := deliveryHttp.NewRatingHandler(ratingSvc)
	deliveryHttp.SetupRoutes(r, ratingHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Menjalankan server HTTP di goroutine terpisah
	go func() {
		log.Printf("[Rating Service] Berjalan di port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Rating Service] Gagal menjalankan HTTP server: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[Rating Service] Memulai proses shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[Rating Service] HTTP Server terpaksa dimatikan: %v", err)
	}

	log.Println("[Rating Service] Sukses dimatikan secara aman.")
}
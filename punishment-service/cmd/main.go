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
	deliveryHttp "github.com/shakilaaulia/Dealan/punishment-service/delivery/http"
	"github.com/shakilaaulia/Dealan/punishment-service/repository"
	"github.com/shakilaaulia/Dealan/punishment-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Inisialisasi Database PostgreSQL menggunakan GORM
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://dealan:dealan_secret@localhost:5432/dealan_db?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Punishment Service] Gagal terhubung ke database: %v", err)
	}
	log.Println("[Punishment Service] Koneksi database PostgreSQL berhasil terjalin")

	// 2. Inisialisasi Repository dan Service
	repo := repository.NewPostgresRepository(db)
	svc := service.NewPunishmentService(repo)

	// 3. Setup Gin HTTP Server & Rute
	r := gin.Default()
	handler := deliveryHttp.NewPunishmentHandler(svc)
	deliveryHttp.SetupRoutes(r, handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Menjalankan server HTTP di goroutine terpisah
	go func() {
		log.Printf("[Punishment Service] Berjalan di port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Punishment Service] Gagal menjalankan HTTP server: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[Punishment Service] Memulai proses shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[Punishment Service] HTTP Server terpaksa dimatikan: %v", err)
	}

	log.Println("[Punishment Service] Sukses dimatikan secara aman.")
}
package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Tambahan: Import package buat baca file .env

	deliveryHttp "github.com/najmialifah/Dealan/auth-service/delivery/http"
	"github.com/najmialifah/Dealan/auth-service/domain"
	"github.com/najmialifah/Dealan/auth-service/pkg/kafka"
	"github.com/najmialifah/Dealan/auth-service/repository"
	"github.com/najmialifah/Dealan/auth-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Memuat konfigurasi dari file .env
	// Ini wajib dipanggil duluan biar program tahu harus nyari password ke mana
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan atau gagal dibaca, menggunakan variabel environment default")
	}

	// 2. Setup Environment Variables
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		// Menggunakan default local untuk pengembangan (diubah ke 127.0.0.1 biar aman)
		dbURL = "postgres://dealan:dealan_secret@127.0.0.1:5432/dealan_db?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dealan_super_secret_key_12345!"
	}

	kafkaBrokersEnv := os.Getenv("KAFKA_BROKERS")
	var kafkaBrokers []string
	if kafkaBrokersEnv != "" {
		kafkaBrokers = strings.Split(kafkaBrokersEnv, ",")
	} else {
		kafkaBrokers = []string{"localhost:9092"}
	}

	// 3. Inisialisasi PostgreSQL GORM
	log.Println("Menghubungkan ke PostgreSQL di:", dbURL)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// 4. AutoMigrate Skema Database Auth
	log.Println("Menjalankan migrasi database otomatis...")
	err = db.AutoMigrate(&domain.AuthCredential{}, &domain.OTPCode{}, &domain.RefreshToken{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	// 5. Inisialisasi Kafka Producer
	var producer service.EventProducer
	if len(kafkaBrokers) > 0 {
		log.Println("Menginisialisasi Kafka Producer dengan Broker:", kafkaBrokers)
		producer = kafka.NewKafkaProducer(kafkaBrokers)
	}

	// 6. Inisialisasi Repository, Service, dan Handler
	repo := repository.NewPostgresAuthRepository(db)
	authSvc := service.NewAuthService(repo, producer, jwtSecret, 24*time.Hour)
	handler := deliveryHttp.NewAuthHandler(authSvc)

	// 7. Inisialisasi Gin Router
	r := gin.Default()

	// Menambahkan Middleware CORS sederhana
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup API Routes
	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.Login)
	r.POST("/auth/validate", handler.Validate)

	// Menjalankan server di port 3001 sesuai spesifikasi docker-compose
	log.Println("Auth Service berjalan pada port :3001")
	if err := r.Run(":3001"); err != nil {
		log.Fatalf("Gagal menjalankan server HTTP: %v", err)
	}
}
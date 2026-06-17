package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	deliveryHttp "github.com/shakilaaulia/Dealan/user-service/delivery/http"
	"github.com/shakilaaulia/Dealan/user-service/domain"
	"github.com/shakilaaulia/Dealan/user-service/pkg/kafka"
	"github.com/shakilaaulia/Dealan/user-service/repository"
	"github.com/shakilaaulia/Dealan/user-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/dealan?sslmode=disable"
	}

	kafkaBrokersEnv := os.Getenv("KAFKA_BROKERS")
	var kafkaBrokers []string
	if kafkaBrokersEnv != "" {
		kafkaBrokers = strings.Split(kafkaBrokersEnv, ",")
	} else {
		kafkaBrokers = []string{"localhost:9092"}
	}

	// 1. Inisialisasi PostgreSQL GORM
	log.Println("Menghubungkan ke PostgreSQL di:", dbURL)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// 2. AutoMigrate Tabel User
	log.Println("Menjalankan migrasi database otomatis...")
	err = db.AutoMigrate(&domain.User{}, &domain.UserProfile{}, &domain.UserRating{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	// 3. Inisialisasi Repository dan Service
	repo := repository.NewPostgresUserRepository(db)
	userSvc := service.NewUserService(repo)
	handler := deliveryHttp.NewUserHandler(userSvc)

	// 4. Inisialisasi & Jalankan Kafka Consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(kafkaBrokers) > 0 {
		consumer := kafka.NewKafkaConsumer(kafkaBrokers, "user-service-group", userSvc)
		go consumer.Start(ctx)
	}

	// 5. Inisialisasi Gin Router
	r := gin.Default()

	// Middleware CORS sederhana
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-User-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup API Routes
	r.GET("/users/profile", handler.GetProfile)
	r.PUT("/users/profile", handler.UpdateProfile)
	r.GET("/users/internal", handler.GetInternalName)

	// Menjalankan server di port 3002 sesuai dengan arsitektur
	log.Println("User Service berjalan pada port :3002")
	if err := r.Run(":3002"); err != nil {
		log.Fatalf("Gagal menjalankan server HTTP: %v", err)
	}
}

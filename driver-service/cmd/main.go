package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	deliveryHttp "github.com/shakilaaulia/Dealan/driver-service/delivery/http"
	"github.com/shakilaaulia/Dealan/driver-service/domain"
	"github.com/shakilaaulia/Dealan/driver-service/pkg/kafka"
	"github.com/shakilaaulia/Dealan/driver-service/repository"
	"github.com/shakilaaulia/Dealan/driver-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://dealan:dealan_secret@localhost:5432/dealan_db?sslmode=disable"
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

	// 2. AutoMigrate Tabel Driver, Vehicle, Status, Rating (dengan Foreign Key)
	log.Println("Menjalankan migrasi database otomatis...")
	err = db.AutoMigrate(&domain.Driver{}, &domain.Vehicle{}, &domain.DriverStatus{}, &domain.DriverRating{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	// 3. Inisialisasi Repository dan Service
	repo := repository.NewPostgresDriverRepository(db)
	driverSvc := service.NewDriverService(repo)
	handler := deliveryHttp.NewDriverHandler(driverSvc)

	// 4. Inisialisasi & Jalankan Kafka Consumer secara asinkron
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(kafkaBrokers) > 0 {
		consumer := kafka.NewKafkaConsumer(kafkaBrokers, "driver-service-group", driverSvc)
		go consumer.Start(ctx)
	}

	// 5. Inisialisasi Gin Router
	r := gin.Default()

	// Middleware CORS sederhana
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Driver-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup API Routes
	r.PUT("/drivers/location", handler.UpdateLocation)
	r.PATCH("/drivers/status", handler.UpdateStatus)
	r.GET("/drivers/profile", handler.GetProfile)
	r.POST("/drivers/vehicles", handler.AddVehicle)

	// Menjalankan server di port 3003 sesuai dengan arsitektur docker-compose
	log.Println("Driver Service berjalan pada port :3003")
	if err := r.Run(":3003"); err != nil {
		log.Fatalf("Gagal menjalankan server HTTP: %v", err)
	}
}

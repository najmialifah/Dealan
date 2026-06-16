package main

import (
	"log"
	"matching-service/controller"
	"matching-service/repository"
	"matching-service/routes"
	"matching-service/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Inisialisasi Database (PostgreSQL / PostGIS)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=location_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// 2. Inisialisasi Kafka Producer
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	// 3. Dependency Injection
	repo := repository.NewMatchingRepository(db)
	svc := service.NewMatchingService(repo, kafkaWriter, "MATCH_FOUND")
	ctrl := controller.NewMatchingController(svc)

	// 4. Setup Gin
	router := gin.Default()
	routes.SetupRoutes(router, ctrl)

	// 5. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("Matching Service berjalan di port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server gagal dijalankan: %v", err)
	}
}

package main

import (
	"log"
	"order-service/controller"
	"order-service/models"
	"order-service/repository"
	"order-service/routes"
	"order-service/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Inisialisasi Koneksi Database (PostgreSQL dengan GORM)
	// DSN sebaiknya diambil dari environment variable di production
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=order_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// Migrasi otomatis untuk tabel Order
	err = db.AutoMigrate(&models.Order{})
	if err != nil {
		log.Fatalf("Gagal melakukan auto-migrate: %v", err)
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
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, kafkaWriter, "ORDER_CREATED")
	orderController := controller.NewOrderController(orderService)

	// 4. Setup Gin Router
	router := gin.Default()
	routes.SetupRoutes(router, orderController)

	// 5. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Order Service berjalan di port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

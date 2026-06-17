package main

import (
	"log"
	"order-service/delivery/http"
	"order-service/domain"
	"order-service/repository"
	"order-service/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Memuat konfigurasi dari file .env (opsional, tidak error jika file tidak ada)
	godotenv.Load()

	// 2. Inisialisasi Koneksi Database (PostgreSQL dengan GORM)
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=order_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// Migrasi otomatis untuk tabel Order
	err = db.AutoMigrate(&domain.Order{})
	if err != nil {
		log.Fatalf("Gagal melakukan auto-migrate: %v", err)
	}

	// 3. Inisialisasi Kafka Producer
	kafkaBroker := os.Getenv("KAFKA_BROKERS")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	// 4. Dependency Injection
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, kafkaWriter, "ORDER_CREATED")
	orderHandler := http.NewOrderHandler(orderService)

	// 5. Setup Gin Router
	router := gin.Default()

	// Menambahkan Middleware CORS
	router.Use(func(c *gin.Context) {
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

	api := router.Group("/api/v1")
	{
		orderRoutes := api.Group("/orders")
		{
			// Endpoint untuk membuat order baru sesuai skenario Postman "Create Order - Valid"
			orderRoutes.POST("/", orderHandler.CreateOrder)
		}
	}

	// 6. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3004" // Ensure port 3004 is the default as per Postman tests
	}
	log.Printf("Order Service berjalan di port :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
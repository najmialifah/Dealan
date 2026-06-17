package main

import (
	"log"
	"net/http"
	"os"

	deliveryHttp "matching-service/delivery/http"
	"matching-service/repository"
	"matching-service/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Memuat konfigurasi dari file .env (opsional)
	godotenv.Load()

	// 1. SETUP DATABASE (PostgreSQL / Supabase)
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=dealan port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("[Warning] Gagal connect Database:", err)
	} else {
		log.Println("[Info] Berhasil connect ke Database PostgreSQL")
	}

	// 2. SETUP KAFKA PRODUCER
	kafkaBroker := os.Getenv("KAFKA_BROKERS")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    "order.matched",
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	// 3. WIRING (Merakit Clean Architecture)
	repo := repository.NewMatchingRepository(db)
	svc := service.NewMatchingService(repo, kafkaWriter, "order.matched")
	handler := deliveryHttp.NewMatchingHandler(svc)

	// 4. SETUP ROUTER GIN
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

	// Endpoint Health Check
	router.GET("/matching/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "matching-service",
		})
	})

	api := router.Group("/api/v1")
	{
		// Endpoint untuk mencocokkan driver
		api.POST("/match", handler.MatchDriver)
	}

	// 6. JALANKAN SERVER
	port := os.Getenv("PORT")
	if port == "" {
		port = "3005"
	}
	log.Printf("Service berjalan aman. Menunggu request di port :%s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}
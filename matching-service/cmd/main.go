package main

import (
	"log"
	"net/http"
	"os"

	"matching-service/controller"
	//"matching-service/domain"
	"matching-service/repository"
	"matching-service/routes"
	"matching-service/service"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. SETUP DATABASE (PostgreSQL / Supabase)
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
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
	kafkaBrokersEnv := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersEnv == "" {
		kafkaBrokersEnv = "localhost:9092"
	}
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokersEnv),
		Topic:    "order.matched",
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	// 3. WIRING (Merakit Clean Architecture)
	repo := repository.NewMatchingRepository(db)
	svc := service.NewMatchingService(repo, kafkaWriter, "order.matched")
	ctrl := controller.NewMatchingController(svc)

	// 4. SETUP ROUTER GIN
	router := gin.Default()

	// Endpoint Health Check
	router.GET("/matching/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "matching-service",
		})
	})

	// 5. DAFTARKAN ROUTES DARI FOLDER ROUTES
	routes.SetupRoutes(router, ctrl)

	// 6. JALANKAN SERVER
	log.Println("Service berjalan aman. Menunggu request di port 3005...")
	if err := router.Run(":3005"); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}
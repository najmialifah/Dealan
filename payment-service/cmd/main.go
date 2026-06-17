package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	deliveryHttp "github.com/shakilaaulia/Dealan/payment-service/delivery/http"
	"github.com/shakilaaulia/Dealan/payment-service/domain"
	"github.com/shakilaaulia/Dealan/payment-service/repository"
	"github.com/shakilaaulia/Dealan/payment-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// kafkaProducerImpl mengimplementasikan service.KafkaProducer menggunakan segmentio/kafka-go.
type kafkaProducerImpl struct {
	writer *kafka.Writer
}

func (k *kafkaProducerImpl) PublishPaymentCompleted(ctx context.Context, trxID, orderID, status string, driverEarnings float64) error {
	eventPayload := map[string]interface{}{
		"event": "PAYMENT_COMPLETED",
		"payload": map[string]interface{}{
			"transaction_id":  trxID,
			"order_id":        orderID,
			"status":          status,
			"driver_earnings": driverEarnings,
		},
	}

	payloadBytes, err := json.Marshal(eventPayload)
	if err != nil {
		return fmt.Errorf("gagal merubah payload event ke JSON: %w", err)
	}

	err = k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(trxID),
		Value: payloadBytes,
	})
	if err != nil {
		return fmt.Errorf("gagal mengirimkan event ke broker Kafka: %w", err)
	}

	return nil
}

func main() {
	// 1. Baca Konfigurasi Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3006" // Sesuai dengan skenario Postman
	}

	// 2. Koneksi Database PostgreSQL menggunakan GORM
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		// Konfigurasi default (bisa diarahkan ke Supabase / Postgres lokal)
		dbURL = "postgres://postgres:postgres@localhost:5432/dealan?sslmode=disable"
	}
	log.Printf("Menghubungkan ke database PostgreSQL di: %s\n", dbURL)

	var db *gorm.DB
	var err error
	// Retry loop untuk menunggu database siap jika dijalankan di docker-compose
	for i := 1; i <= 5; i++ {
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Percobaan koneksi database ke-%d gagal, menunggu 5 detik...\n", i)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Gagal menghubungkan ke database PostgreSQL setelah beberapa percobaan: %v", err)
	}

	// 3. Auto-Migrate skema database sesuai dengan spesifikasi
	log.Println("Memulai auto-migrasi skema database...")
	err = db.AutoMigrate(
		&domain.Transaction{},
		&domain.DriverWallet{},
		&domain.WalletTransaction{},
		&domain.PaymentLog{},
		&domain.IdempotencyKey{},
	)
	if err != nil {
		log.Fatalf("Gagal melakukan auto-migrasi database: %v", err)
	}
	log.Println("Auto-migrasi database berhasil dilakukan.")

	// 4. Inisialisasi Kafka Broker
	kafkaBrokersEnv := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersEnv == "" {
		kafkaBrokersEnv = "localhost:9092"
	}
	brokers := strings.Split(kafkaBrokersEnv, ",")

	log.Printf("Menghubungkan ke Kafka broker: %v\n", brokers)
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "payment.completed",
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	// 5. Setup Dependency Injection (Clean Architecture)
	paymentRepo := repository.NewPaymentRepository(db)
	kafkaProd := &kafkaProducerImpl{writer: kafkaWriter}
	paymentService := service.NewPaymentService(paymentRepo, kafkaProd)
	paymentHandler := deliveryHttp.NewPaymentHandler(paymentService)

	// 6. Setup Gin Engine & Routing
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Menambahkan Middleware CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Idempotency-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Hubungkan route endpoint
	api := router.Group("/payments")
	{
		api.POST("/create", paymentHandler.Create)
		api.POST("/webhook", paymentHandler.Webhook)
		api.GET("/:transaction_id", paymentHandler.GetStatus)
		api.GET("/driver/:driver_id/wallet", paymentHandler.GetDriverWallet)
	}

	// Endpoint tambahan untuk health check k8s
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// 7. Jalankan Server HTTP
	addr := ":" + port
	log.Printf("Payment Service berhasil dimulai dan mendengarkan di %s\n", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Gagal menjalankan server Gin: %v", err)
	}
}
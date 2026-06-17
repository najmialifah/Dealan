package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shakilaaulia/Dealan/shipment-service/controller"
	"github.com/shakilaaulia/Dealan/shipment-service/domain"
	"github.com/shakilaaulia/Dealan/shipment-service/repository"
	"github.com/shakilaaulia/Dealan/shipment-service/routes"
	"github.com/shakilaaulia/Dealan/shipment-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Baca Konfigurasi Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8094" // Port default sesuai spesifikasi kubernetes
	}

	// 2. Koneksi Database PostgreSQL menggunakan GORM
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/dealan?sslmode=disable"
	}
	log.Printf("Menghubungkan ke database PostgreSQL di: %s\n", dbURL)

	var db *gorm.DB
	var err error
	// Retry loop untuk menunggu database siap
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
	err = db.AutoMigrate(&domain.Shipment{})
	if err != nil {
		log.Fatalf("Gagal melakukan auto-migrasi database: %v", err)
	}
	log.Println("Auto-migrasi database berhasil dilakukan.")

	// 4. Setup Dependency Injection (Clean Architecture)
	shipmentRepo := repository.NewShipmentRepository(db)
	shipmentService := service.NewShipmentService(shipmentRepo)
	shipmentCtrl := controller.NewShipmentController(shipmentService)

	// 5. Setup Gin Engine & Routing
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Hubungkan route endpoint
	routes.SetupRoutes(router, shipmentCtrl)

	// Endpoint tambahan untuk health check k8s
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// 6. Jalankan Server HTTP
	addr := ":" + port
	log.Printf("Shipment Service berhasil dimulai dan mendengarkan di %s\n", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Gagal menjalankan server Gin: %v", err)
	}
}
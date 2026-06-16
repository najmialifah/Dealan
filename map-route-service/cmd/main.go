package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"map-route-service/controller"
	"map-route-service/models"
	"map-route-service/repository"
	"map-route-service/routes"
	"map-route-service/service"
)

func main() {
	log.Println("Memulai Map Route Service...")

	// 1. Konfigurasi Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088" // Menggunakan port 8088 sesuai dengan kode lama
	}

	// 2. Hubungkan ke PostgreSQL Database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=dealan password=dealan_secret dbname=dealan_db port=5432 sslmode=disable"
	}

	var db *gorm.DB
	var err error
	maxRetries := 5

	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err == nil {
			log.Println("Berhasil terhubung ke database PostgreSQL untuk Map Route!")
			break
		}
		log.Printf("Gagal terhubung ke DB Map Route (percobaan %d/%d): %v. Mencoba kembali dalam 3 detik...", i, maxRetries, err)
		time.Sleep(3 * time.Second)
	}

	// 3. Auto Migration untuk model MapRoute
	if db != nil {
		log.Println("Menjalankan migrasi database Map Route...")
		err = db.AutoMigrate(&models.MapRoute{})
		if err != nil {
			log.Fatalf("Gagal melakukan migrasi database: %v", err)
		}
	} else {
		log.Println("[Peringatan] Berjalan tanpa koneksi database aktif untuk Map Route Service.")
	}

	// 4. Inisialisasi Clean Architecture layers
	mapRepo := repository.NewMapRepository(db)
	mapSvc := service.NewMapService(mapRepo)
	mapCtrl := controller.NewMapController(mapSvc)

	// 5. Setup Gin Router
	r := gin.Default()

	// Register routes
	routes.SetupRoutes(r, mapCtrl)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": "map-route-service"})
	})

	log.Printf("Map Route Service berjalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server HTTP: %v", err)
	}
}
package main

import (
	"log"
	"os"

	deliveryHttp "location-service/delivery/http"
	"location-service/repository"
	"location-service/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Memuat konfigurasi dari file .env (opsional)
	godotenv.Load()

	// 1. Inisialisasi Database (PostGIS)
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=location_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// Buat ekstensi PostGIS (jika belum ada)
	db.Exec("CREATE EXTENSION IF NOT EXISTS postgis;")

	// Tabel driver_locations dibuat manual dengan query ini jika AutoMigrate tidak mendukung tipe spatial custom
	// Kita sediakan DDL jika tabel belum ada
	ddl := `
	CREATE TABLE IF NOT EXISTS driver_locations (
		driver_id INTEGER PRIMARY KEY,
		location GEOMETRY(Point, 4326),
		updated_at TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS location_gix ON driver_locations USING GIST (location);
	`
	db.Exec(ddl)

	// 2. Dependency Injection
	repo := repository.NewLocationRepository(db)
	svc := service.NewLocationService(repo)
	handler := deliveryHttp.NewLocationHandler(svc)

	// 3. Setup Gin
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

	api := router.Group("/api/v1/locations")
	{
		api.POST("/update", handler.UpdateLocation)
		api.GET("/nearby", handler.FindNearby)
	}

	// 4. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3008" // Sesuai dengan skenario Postman
	}
	log.Printf("Location Service berjalan di port :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server gagal dijalankan: %v", err)
	}
}
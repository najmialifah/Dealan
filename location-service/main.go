package main

import (
	"location-service/controller"
	"location-service/repository"
	"location-service/routes"
	"location-service/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Inisialisasi Database (PostGIS)
	dsn := os.Getenv("DATABASE_URL")
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
	ctrl := controller.NewLocationController(svc)

	// 3. Setup Gin
	router := gin.Default()
	routes.SetupRoutes(router, ctrl)

	// 4. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Location Service berjalan di port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server gagal dijalankan: %v", err)
	}
}

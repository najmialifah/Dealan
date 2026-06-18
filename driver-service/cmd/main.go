package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	deliveryHttp "github.com/najmialifah/Dealan/driver-service/delivery/http"
	"github.com/najmialifah/Dealan/driver-service/domain"
	"github.com/najmialifah/Dealan/driver-service/pkg/kafka"
	"github.com/najmialifah/Dealan/driver-service/repository"
	"github.com/najmialifah/Dealan/driver-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/dealan?sslmode=disable"
	}

	kafkaBrokersEnv := os.Getenv("KAFKA_BROKERS")
	var kafkaBrokers []string
	if kafkaBrokersEnv != "" {
		kafkaBrokers = strings.Split(kafkaBrokersEnv, ",")
	} else {
		kafkaBrokers = []string{"localhost:9092"}
	}

	// 1. Inisialisasi PostgreSQL GORM
	log.Println("Menghubungkan ke PostgreSQL di:", dbURL)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// 2. AutoMigrate Tabel Driver, Vehicle, Status, Rating (dengan Foreign Key)
	log.Println("Menjalankan migrasi database otomatis...")
	err = db.AutoMigrate(&domain.Driver{}, &domain.Vehicle{}, &domain.DriverStatus{}, &domain.DriverRating{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	// Setup PostGIS dan sinkronisasi kolom lokasi
	log.Println("Menyiapkan ekstensi PostGIS dan kolom lokasi...")
	db.Exec("CREATE EXTENSION IF NOT EXISTS postgis;")
	db.Exec("ALTER TABLE driver_status ADD COLUMN IF NOT EXISTS lokasi GEOGRAPHY(Point, 4326);")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_driver_lokasi ON driver_status USING GIST(lokasi);")

	// Trigger untuk sinkronisasi kolom lokasi otomatis dari lat dan long
	triggerSQL := `
	CREATE OR REPLACE FUNCTION update_driver_status_lokasi()
	RETURNS TRIGGER AS $$
	BEGIN
		IF NEW.lat IS NOT NULL AND NEW.long IS NOT NULL THEN
			NEW.lokasi := ST_SetSRID(ST_MakePoint(NEW.long, NEW.lat), 4326)::geography;
		END IF;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS trigger_update_driver_status_lokasi ON driver_status;
	CREATE TRIGGER trigger_update_driver_status_lokasi
	BEFORE INSERT OR UPDATE OF lat, long ON driver_status
	FOR EACH ROW
	EXECUTE FUNCTION update_driver_status_lokasi();
	`
	if err := db.Exec(triggerSQL).Error; err != nil {
		log.Printf("[Warning] Gagal membuat trigger sinkronisasi lokasi: %v", err)
	} else {
		log.Println("[Info] Trigger sinkronisasi lokasi berhasil didaftarkan")
	}

	// 3. Inisialisasi Repository dan Service
	repo := repository.NewPostgresDriverRepository(db)
	driverSvc := service.NewDriverService(repo)
	handler := deliveryHttp.NewDriverHandler(driverSvc)

	// 4. Inisialisasi & Jalankan Kafka Consumer secara asinkron
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(kafkaBrokers) > 0 {
		consumer := kafka.NewKafkaConsumer(kafkaBrokers, "driver-service-group", driverSvc)
		go consumer.Start(ctx)
	}

	// 5. Inisialisasi Gin Router
	r := gin.Default()

	// Middleware CORS sederhana
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Driver-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup API Routes
	r.PUT("/drivers/location", handler.UpdateLocation)
	r.PATCH("/drivers/status", handler.UpdateStatus)
	r.GET("/drivers/profile", handler.GetProfile)
	r.POST("/drivers/vehicles", handler.AddVehicle)

	// Menjalankan server di port 3003 sesuai dengan arsitektur docker-compose
	log.Println("Driver Service berjalan pada port :3003")
	if err := r.Run(":3003"); err != nil {
		log.Fatalf("Gagal menjalankan server HTTP: %v", err)
	}
}

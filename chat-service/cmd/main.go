package main

import (
	"log"
	"os"
	"time"

	"chat-service/controller"
	"chat-service/models"
	"chat-service/repository"
	"chat-service/routes"
	"chat-service/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("Memulai Chat Service...")

	// 1. Ambil Port dari Environment Variables atau gunakan default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8087" // Menggunakan port 8087 sesuai dengan spesifikasi lama
	}

	// 2. Ambil Koneksi Database PostgreSQL URL
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
			log.Println("Berhasil terhubung ke database PostgreSQL untuk Chat Service!")
			break
		}
		log.Printf("Gagal terhubung ke DB Chat (percobaan %d/%d): %v. Mencoba kembali dalam 3 detik...", i, maxRetries, err)
		time.Sleep(3 * time.Second)
	}

	// 3. Auto Migration untuk models ChatRoom dan ChatMessage
	if db != nil {
		log.Println("Menjalankan migrasi database Chat...")
		err = db.AutoMigrate(&models.ChatRoom{}, &models.ChatMessage{})
		if err != nil {
			log.Fatalf("Gagal melakukan migrasi database chat: %v", err)
		}
	} else {
		log.Println("[Peringatan] Berjalan tanpa koneksi database aktif untuk Chat Service.")
	}

	// 4. Inisialisasi Clean Architecture layers
	chatRepo := repository.NewChatRepository(db)
	chatSvc := service.NewChatService(chatRepo)
	chatCtrl := controller.NewChatController(chatSvc)

	// 5. Jalankan WebSocket Hub di background goroutine
	go chatSvc.Run()

	// 6. Setup Gin HTTP Server
	r := gin.Default()
	routes.SetupRoutes(r, chatCtrl)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": "chat-service"})
	})

	log.Printf("Chat Service berjalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server HTTP: %v", err)
	}
}
package repository

import (
	"context"
	"time"

	"github.com/najmialifah/Dealan/notification-service/domain"
	"gorm.io/gorm"
)

// NotificationLog merepresentasikan entitas database untuk riwayat notifikasi
type NotificationLog struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TargetID       string    `gorm:"type:varchar(100);not null"`
	Title          string    `gorm:"type:varchar(255);not null"`
	Body           string    `gorm:"type:text;not null"`
	ActionLink     string    `gorm:"type:varchar(255)"`
	DeliveryStatus string    `gorm:"type:varchar(50);not null"`
	MessageID      string    `gorm:"type:varchar(100);not null"`
	CreatedAt      time.Time `gorm:"default:now()"`
}

// TableName menentukan nama tabel di database
func (NotificationLog) TableName() string {
	return "notification_logs"
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository membuat instance baru dari repository PostgreSQL
func NewPostgresRepository(db *gorm.DB) NotificationRepository {
	// AutoMigrate tabel saat inisialisasi
	db.AutoMigrate(&NotificationLog{})
	return &postgresRepository{db: db}
}

// SaveLog menyimpan riwayat pengiriman notifikasi ke database
func (r *postgresRepository) SaveLog(ctx context.Context, req domain.NotificationRequest, res domain.NotificationResponse) error {
	log := NotificationLog{
		TargetID:       req.TargetID,
		Title:          req.Title,
		Body:           req.Body,
		ActionLink:     req.ActionLink,
		DeliveryStatus: res.DeliveryStatus,
		MessageID:      res.MessageID,
	}
	return r.db.WithContext(ctx).Create(&log).Error
}

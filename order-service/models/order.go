package models

import (
	"time"

	"gorm.io/gorm"
)

// Order merepresentasikan tabel orders di database
type Order struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	DriverID    *uint          `json:"driver_id"` // Bisa null jika belum ada driver
	Status      string         `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"`
	DetailPaket map[string]any `gorm:"type:jsonb" json:"detail_paket"` // Menggunakan tipe data JSONB di PostgreSQL
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONB adalah custom type untuk menangani kolom JSONB di PostgreSQL secara fleksibel.
type JSONB map[string]interface{}

// Value mengubah map Go menjadi representasi JSON byte untuk disimpan di PostgreSQL.
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan mengurai data JSON byte dari PostgreSQL ke map Go.
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("assertion tipe data ke []byte gagal")
	}
	return json.Unmarshal(bytes, j)
}

// PricingRule menyimpan konfigurasi tarif dasar dan variasi aturan tarif dalam format JSONB.
type PricingRule struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServiceType string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"service_type"` // ride, car, send
	BasePrice   float64   `gorm:"type:decimal(15,2);not null" json:"base_price"`             // tarif dasar (misal untuk jarak awal)
	Active      bool      `gorm:"default:true" json:"active"`
	Config      JSONB     `gorm:"type:jsonb" json:"config"`                                  // konfigurasi tambahan seperti per_km_rate, rush_hour_multiplier, dll.
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PricingNegotiation menyimpan riwayat pengajuan penawaran harga dari pengguna.
type PricingNegotiation struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID        string    `gorm:"type:varchar(50);not null;index" json:"order_id"`
	OriginalPrice  float64   `gorm:"type:decimal(15,2);not null" json:"original_price"`
	RequestedPrice float64   `gorm:"type:decimal(15,2);not null" json:"requested_price"`
	Status         string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, approved, rejected
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// PricingEstimateRequest mendefinisikan parameter untuk kalkulasi harga.
type PricingEstimateRequest struct {
	ServiceType string  `json:"service_type" binding:"required"`
	Distance    float64 `json:"distance" binding:"required"` // Jarak dalam km
	Weight      float64 `json:"weight"`                      // Berat barang dalam kg (khusus GoSend)
	IsRushHour  bool    `json:"is_rush_hour"`                // Menandai jam sibuk
	IsNight     bool    `json:"is_night"`                    // Menandai larat malam
}

// PricingEstimateResponse mengembalikan hasil kalkulasi harga beserta batas negosiasi.
type PricingEstimateResponse struct {
	ServiceType    string  `json:"service_type"`
	Distance       float64 `json:"distance"`
	EstimatedPrice float64 `json:"estimated_price"`
	MinPrice       float64 `json:"min_price"`
	MaxPrice       float64 `json:"max_price"`
}

// NegotiationRequest mendefinisikan data penawaran harga baru.
type NegotiationRequest struct {
	OrderID        string  `json:"order_id" binding:"required"`
	OriginalPrice  float64 `json:"original_price" binding:"required"`
	RequestedPrice float64 `json:"requested_price" binding:"required"`
}

// NegotiationResponse mengembalikan status persetujuan penawaran harga.
type NegotiationResponse struct {
	OrderID        string  `json:"order_id"`
	OriginalPrice  float64 `json:"original_price"`
	RequestedPrice float64 `json:"requested_price"`
	Status         string  `json:"status"` // approved atau rejected
}

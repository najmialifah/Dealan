package domain

import (
	"time"
)

// Driver merepresentasikan model utama pengemudi beserta relasi-relasinya
type Driver struct {
	ID           string        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Nama         string        `gorm:"type:varchar(100);not null" json:"nama"`
	NomorHP      string        `gorm:"type:varchar(20);uniqueIndex;not null" json:"nomor_hp"`
	Email        string        `gorm:"type:varchar(100)" json:"email"`
	FotoProfil   string        `gorm:"type:varchar(500)" json:"foto_profil,omitempty"`
	FotoKTP      string        `gorm:"type:varchar(500)" json:"foto_ktp,omitempty"`
	Status       string        `gorm:"type:varchar(20);default:'active'" json:"status"` // active, suspended, banned
	CreatedAt    time.Time     `gorm:"default:now()" json:"created_at"`
	UpdatedAt    time.Time     `gorm:"default:now()" json:"updated_at"`
	Vehicles     []Vehicle     `gorm:"foreignKey:DriverID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"vehicles,omitempty"`
	DriverStatus DriverStatus  `gorm:"foreignKey:DriverID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"driver_status,omitempty"`
	DriverRating DriverRating  `gorm:"foreignKey:DriverID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"driver_rating,omitempty"`
}

// Vehicle merepresentasikan tabel kendaraan milik driver
type Vehicle struct {
	ID            string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DriverID      string `gorm:"type:uuid;not null;index" json:"driver_id"`
	Jenis         string `gorm:"type:varchar(10);not null" json:"jenis"` // motor, mobil
	Merek         string `gorm:"type:varchar(50)" json:"merek"`
	Model         string `gorm:"type:varchar(50)" json:"model"`
	PlatNomor     string `gorm:"type:varchar(20);uniqueIndex;not null" json:"plat_nomor"`
	Tahun         int    `json:"tahun"`
	FotoKendaraan string `gorm:"type:varchar(500)" json:"foto_kendaraan,omitempty"`
	IsVerified    bool   `gorm:"default:false" json:"is_verified"`
}

// DriverStatus merepresentasikan status online dan lokasi GPS driver
type DriverStatus struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DriverID       string    `gorm:"type:uuid;uniqueIndex;not null" json:"driver_id"`
	IsOnline       bool      `gorm:"default:false" json:"is_online"`
	LayananAktif   string    `gorm:"type:varchar(50);default:''" json:"layanan_aktif"` // Disimpan sebagai string terpisah koma (misal: "ride,send")
	CurrentOrderID string    `gorm:"type:uuid" json:"current_order_id,omitempty"`
	LastSeen       time.Time `gorm:"default:now()" json:"last_seen"`
	Lat            float64   `gorm:"type:decimal(10,8)" json:"lat"`
	Long           float64   `gorm:"type:decimal(11,8)" json:"long"`
}

// DriverRating merepresentasikan agregasi rating driver
type DriverRating struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DriverID    string    `gorm:"type:uuid;uniqueIndex;not null" json:"driver_id"`
	TotalRating float64   `gorm:"type:decimal(3,2);default:5.00" json:"total_rating"`
	TotalReview int       `gorm:"default:0" json:"total_review"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`
}

// UpdateLocationRequest merepresentasikan payload update lokasi koordinat GPS
type UpdateLocationRequest struct {
	Lat  float64 `json:"lat" binding:"required"`
	Long float64 `json:"long" binding:"required"`
}

// UpdateStatusRequest merepresentasikan payload update status online driver
type UpdateStatusRequest struct {
	IsOnline     bool   `json:"is_online"`
	LayananAktif string `json:"layanan_aktif"` // misal: "ride,send"
}

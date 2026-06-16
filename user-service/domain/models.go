package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONMap didefinisikan untuk mendukung kolom bertipe JSONB di PostgreSQL GORM
type JSONMap map[string]interface{}

// Value men-serialize map menjadi byte JSON untuk disimpan di database
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(j)
}

// Scan men-deserialize data byte JSON dari database menjadi JSONMap
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// User merepresentasikan model pengguna utama
type User struct {
	ID          string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Nama        string       `gorm:"type:varchar(100);not null" json:"nama"`
	Email       string       `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	NomorHP     string       `gorm:"type:varchar(20);uniqueIndex;not null" json:"nomor_hp"`
	Alamat      string       `gorm:"type:text" json:"alamat"`
	FotoProfil  string       `gorm:"type:varchar(500)" json:"foto_profil,omitempty"`
	Status      string       `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt   time.Time    `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"default:now()" json:"updated_at"`
	UserProfile UserProfile  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_profile,omitempty"`
	UserRating  UserRating   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_rating,omitempty"`
}

// UserProfile merepresentasikan tabel detail profil tambahan user
type UserProfile struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     string    `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Preferensi JSONMap   `gorm:"type:jsonb;default:'{}'" json:"preferensi"`
	TotalOrder int       `gorm:"default:0" json:"total_order"`
	TotalSpend float64   `gorm:"type:decimal(15,2);default:0.00" json:"total_spend"`
	CreatedAt  time.Time `gorm:"default:now()" json:"created_at"`
}

// UserRating merepresentasikan agregasi rating user
type UserRating struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      string    `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	TotalRating float64   `gorm:"type:decimal(3,2);default:5.00" json:"total_rating"`
	TotalReview int       `gorm:"default:0" json:"total_review"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`
}

// UpdateProfileRequest merepresentasikan request dari REST API untuk update profil
type UpdateProfileRequest struct {
	Nama       string `json:"nama" binding:"required"`
	NomorHP    string `json:"nomor_hp" binding:"required"`
	Alamat     string `json:"alamat"`
	FotoProfil string `json:"foto_profil"`
}

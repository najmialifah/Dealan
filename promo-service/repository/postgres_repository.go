package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Promo merepresentasikan data kupon promo di PostgreSQL GORM
type Promo struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code      string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Discount  float64   `gorm:"type:decimal(15,2);not null"`
	Quota     int       `gorm:"type:integer;not null"`
	Active    bool      `gorm:"type:boolean;default:true"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
}

// TableName menentukan nama tabel promo di database
func (Promo) TableName() string {
	return "promos"
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository membuat instance baru repository promo
func NewPostgresRepository(db *gorm.DB) PromoRepository {
	db.AutoMigrate(&Promo{})
	return &postgresRepository{db: db}
}

// CheckPromo memvalidasi promo dan mengurangi kuota menggunakan Database Lock SELECT ... FOR UPDATE via GORM
func (r *postgresRepository) CheckPromo(ctx context.Context, code string) (float64, error) {
	var discount float64

	// Membuka transaksi database GORM
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var promo Promo

		// Menerapkan SELECT ... FOR UPDATE via Clauses(clause.Locking{Strength: "UPDATE"})
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("code = ? AND active = ?", code, true).
			First(&promo).Error
		if err != nil {
			return err // Promo tidak ditemukan atau tidak aktif
		}

		// Periksa apakah kuota masih tersedia
		if promo.Quota <= 0 {
			return errors.New("kuota promo ini sudah habis")
		}

		// Kurangi kuota promo
		promo.Quota--
		err = tx.Save(&promo).Error
		if err != nil {
			return err
		}

		discount = promo.Discount
		return nil
	})

	if err != nil {
		return 0, err
	}

	return discount, nil
}

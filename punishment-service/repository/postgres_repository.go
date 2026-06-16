package repository

import (
	"context"
	"time"

	"github.com/shakilaaulia/Dealan/punishment-service/domain"
	"gorm.io/gorm"
)

// ViolationLog merepresentasikan log pelanggaran di database
type ViolationLog struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AccountID  string    `gorm:"type:varchar(100);not null"`
	ReasonCode string    `gorm:"type:varchar(50);not null"`
	Duration   int       `gorm:"type:integer;not null"` // dalam jam
	CreatedAt  time.Time `gorm:"default:now()"`
}

// TableName menentukan nama tabel log pelanggaran
func (ViolationLog) TableName() string {
	return "violation_logs"
}

// AccountStatus merepresentasikan status akun di database
type AccountStatus struct {
	AccountID string    `gorm:"type:varchar(100);primaryKey"`
	Status    string    `gorm:"type:varchar(20);not null"` // active, suspended, banned
	UpdatedAt time.Time `gorm:"default:now()"`
}

// TableName menentukan nama tabel status akun
func (AccountStatus) TableName() string {
	return "account_statuses"
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository membuat instance baru dari repository PostgreSQL Punishment
func NewPostgresRepository(db *gorm.DB) PunishmentRepository {
	db.AutoMigrate(&ViolationLog{}, &AccountStatus{})
	return &postgresRepository{db: db}
}

// StoreViolation menyimpan log pelanggaran ke database
func (r *postgresRepository) StoreViolation(ctx context.Context, req domain.PunishmentRequest) (string, error) {
	log := ViolationLog{
		AccountID:  req.AccountID,
		ReasonCode: req.ReasonCode,
		Duration:   req.Duration,
	}

	err := r.db.WithContext(ctx).Create(&log).Error
	if err != nil {
		return "", err
	}
	return log.ID, nil
}

// UpdateAccountStatus mengubah status akun (misalnya suspended atau banned)
func (r *postgresRepository) UpdateAccountStatus(ctx context.Context, id string, status string) error {
	acc := AccountStatus{
		AccountID: id,
		Status:    status,
		UpdatedAt: time.Now(),
	}
	// Gunakan Save untuk melakukan upsert
	return r.db.WithContext(ctx).Save(&acc).Error
}

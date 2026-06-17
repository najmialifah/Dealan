package repository

import (
	"context"
	"time"

	"github.com/najmialifah/Dealan/rating-review-service/domain"
	"gorm.io/gorm"
)

// Review merepresentasikan tabel reviews di PostgreSQL menggunakan GORM
type Review struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID      string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_order_reviewer"`
	ReviewerID   string    `gorm:"type:uuid;not null;uniqueIndex:idx_order_reviewer"`
	ReviewerRole string    `gorm:"type:varchar(10);not null"`
	TargetID     string    `gorm:"type:uuid;not null;index:idx_reviews_target"`
	TargetRole   string    `gorm:"type:varchar(10);not null;index:idx_reviews_target"`
	RatingScore  int       `gorm:"type:smallint;not null"`
	Comment      string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"default:now()"`
}

// TableName menentukan nama tabel di database
func (Review) TableName() string {
	return "reviews"
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository membuat instance baru dari repository PostgreSQL Rating
func NewPostgresRepository(db *gorm.DB) RatingRepository {
	db.AutoMigrate(&Review{})
	return &postgresRepository{db: db}
}

// SaveReview menyimpan review baru ke database
func (r *postgresRepository) SaveReview(ctx context.Context, req domain.RatingRequest) (string, error) {
	review := Review{
		OrderID:      req.OrderID,
		ReviewerID:   req.ReviewerID,
		ReviewerRole: req.ReviewerRole,
		TargetID:     req.TargetID,
		TargetRole:   req.TargetRole,
		RatingScore:  req.RatingScore,
		Comment:      req.Comment,
	}

	err := r.db.WithContext(ctx).Create(&review).Error
	if err != nil {
		return "", err
	}
	return review.ID, nil
}

// GetAverageRating menghitung rata-rata nilai rating secara instan menggunakan fungsi agregasi AVG() SQL di GORM
func (r *postgresRepository) GetAverageRating(ctx context.Context, targetID string) (float64, error) {
	var avgRating float64
	err := r.db.WithContext(ctx).Model(&Review{}).
		Where("target_id = ?", targetID).
		Select("COALESCE(AVG(rating_score), 0)").
		Scan(&avgRating).Error
	return avgRating, err
}

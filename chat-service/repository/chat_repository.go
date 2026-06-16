package repository

import (
	"context"

	"chat-service/models"
	"gorm.io/gorm"
)

// ChatRepository mendefinisikan interface operasi database untuk chat-service.
type ChatRepository interface {
	CreateRoom(ctx context.Context, room *models.ChatRoom) error
	GetRoom(ctx context.Context, orderID string) (*models.ChatRoom, error)
	SaveMessage(ctx context.Context, msg *models.ChatMessage) error
	GetChatHistory(ctx context.Context, orderID string) ([]models.ChatMessage, error)
}

type chatRepositoryImpl struct {
	db *gorm.DB
}

// NewChatRepository membuat instance baru dari ChatRepository.
func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepositoryImpl{db: db}
}

// CreateRoom membuat room obrolan baru antara user dan driver untuk suatu pesanan.
func (r *chatRepositoryImpl) CreateRoom(ctx context.Context, room *models.ChatRoom) error {
	return r.db.WithContext(ctx).Save(room).Error
}

// GetRoom mengambil data room obrolan berdasarkan ID pesanan.
func (r *chatRepositoryImpl) GetRoom(ctx context.Context, orderID string) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// SaveMessage menyimpan pesan chat ke PostgreSQL.
func (r *chatRepositoryImpl) SaveMessage(ctx context.Context, msg *models.ChatMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

// GetChatHistory mengambil semua riwayat pesan untuk suatu pesanan, diurutkan berdasarkan waktu kirim.
func (r *chatRepositoryImpl) GetChatHistory(ctx context.Context, orderID string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Order("sent_at asc").Find(&messages).Error
	return messages, err
}

package service

import (
	"context"
	"testing"

	"chat-service/models"
	"github.com/stretchr/testify/assert"
)

// mockChatRepository adalah mock database untuk pengujian chat-service.
type mockChatRepository struct {
	rooms    map[string]*models.ChatRoom
	messages []models.ChatMessage
	err      error
}

func (m *mockChatRepository) CreateRoom(ctx context.Context, room *models.ChatRoom) error {
	m.rooms[room.OrderID] = room
	return m.err
}

func (m *mockChatRepository) GetRoom(ctx context.Context, orderID string) (*models.ChatRoom, error) {
	if room, exists := m.rooms[orderID]; exists {
		return room, nil
	}
	return nil, m.err
}

func (m *mockChatRepository) SaveMessage(ctx context.Context, msg *models.ChatMessage) error {
	m.messages = append(m.messages, *msg)
	return m.err
}

func (m *mockChatRepository) GetChatHistory(ctx context.Context, orderID string) ([]models.ChatMessage, error) {
	return m.messages, m.err
}

// TestChatService_CreateRoom menguji pembuatan room obrolan baru.
func TestChatService_CreateRoom(t *testing.T) {
	repo := &mockChatRepository{rooms: make(map[string]*models.ChatRoom)}
	svc := NewChatService(repo)

	err := svc.CreateRoom(context.Background(), "ORD-123", "USER-1", "DRIVER-2")
	assert.NoError(t, err)
	assert.Contains(t, repo.rooms, "ORD-123")
	assert.Equal(t, "USER-1", repo.rooms["ORD-123"].UserID)
	assert.Equal(t, "DRIVER-2", repo.rooms["ORD-123"].DriverID)
}

// TestChatService_GetHistory menguji pengambilan riwayat pesan dari database.
func TestChatService_GetHistory(t *testing.T) {
	messages := []models.ChatMessage{
		{OrderID: "ORD-123", SenderID: "USER-1", SenderRole: "user", Message: "Halo", Type: "text"},
		{OrderID: "ORD-123", SenderID: "DRIVER-2", SenderRole: "driver", Message: "Halo juga", Type: "text"},
	}
	repo := &mockChatRepository{
		rooms:    make(map[string]*models.ChatRoom),
		messages: messages,
	}
	svc := NewChatService(repo)

	history, err := svc.GetHistory(context.Background(), "ORD-123")
	assert.NoError(t, err)
	assert.Len(t, history, 2)
	assert.Equal(t, "Halo", history[0].Message)
	assert.Equal(t, "Halo juga", history[1].Message)
}
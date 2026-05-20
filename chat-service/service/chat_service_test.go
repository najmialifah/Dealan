package service

import (
	"context"
	"testing"

	"chat-service/domain"
	"chat-service/mocks"

	"go.uber.org/mock/gomock"
)

func TestChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)

	mockRepo.EXPECT().
		SaveMessage(gomock.Any(), "hello").
		Return("123", nil)

	svc := NewChatService(mockRepo)

	res, _ := svc.Send(context.Background(), domain.ChatRequest{
		Message: "hello",
	})

	if res.MessageID == "" {
		t.Errorf("empty id")
	}
}
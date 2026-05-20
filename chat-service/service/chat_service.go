package service

import (
	"context"
	"time"

	"chat-service/domain"
	"chat-service/repository"
)

type chatService struct {
	repo repository.ChatRepository
}

func NewChatService(r repository.ChatRepository) ChatService {
	return &chatService{r}
}

func (s *chatService) Send(ctx context.Context, req domain.ChatRequest) (*domain.ChatResponse, error) {

	if s.repo == nil {
		return &domain.ChatResponse{
			MessageID: "dummy",
			Status:    true,
		}, nil
	}

	id, _ := s.repo.SaveMessage(ctx, req.Message)

	return &domain.ChatResponse{
		MessageID: id + time.Now().String(),
		Status:    true,
	}, nil
}
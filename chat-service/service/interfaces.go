package service

import (
	"context"
	"chat-service/domain"
)

type ChatService interface {
	Send(ctx context.Context, req domain.ChatRequest) (*domain.ChatResponse, error)
}
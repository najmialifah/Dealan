package repository

import "context"

type ChatRepository interface {
	SaveMessage(ctx context.Context, msg string) (string, error)
}
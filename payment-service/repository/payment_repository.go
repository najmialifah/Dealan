package repository

import (
	"context"

	"github.com/shakilaaulia/Dealan/payment-service/domain"
)

//go:generate mockgen -source=payment_repository.go -destination=../mocks/mock_payment_repository.go -package=mocks

type PaymentRepository interface {
	SaveTransaction(ctx context.Context, req domain.PaymentRequest) (string, error)
	GetTransactionStatus(ctx context.Context, transactionID string) (domain.PaymentResponse, error)
	ProcessGateway(ctx context.Context, transactionID string, nominal float64, metode string) (string, error)
}

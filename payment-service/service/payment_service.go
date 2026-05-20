package service

import (
	"context"

	"github.com/shakilaaulia/Dealan/payment-service/domain"
	"github.com/shakilaaulia/Dealan/payment-service/repository"
)

type PaymentService interface {
	Process(ctx context.Context, req domain.PaymentRequest) (domain.PaymentResponse, error)
	GetStatus(ctx context.Context, transactionID string) (domain.PaymentResponse, error)
}

type paymentServiceImpl struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentServiceImpl{repo: repo}
}

func (s *paymentServiceImpl) Process(ctx context.Context, req domain.PaymentRequest) (domain.PaymentResponse, error) {
	trxID, err := s.repo.SaveTransaction(ctx, req)
	if err != nil {
		return domain.PaymentResponse{}, err
	}

	invoiceURL, err := s.repo.ProcessGateway(ctx, trxID, req.Nominal, req.MetodePembayaran)
	if err != nil {
		return domain.PaymentResponse{}, err
	}

	return domain.PaymentResponse{
		TransactionID: trxID,
		PaymentStatus: "PENDING",
		InvoiceURL:    invoiceURL,
	}, nil
}

func (s *paymentServiceImpl) GetStatus(ctx context.Context, transactionID string) (domain.PaymentResponse, error) {
	return s.repo.GetTransactionStatus(ctx, transactionID)
}

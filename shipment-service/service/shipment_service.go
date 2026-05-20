package service

import (
	"context"

	"github.com/shakilaaulia/Dealan/shipment-service/domain"
	"github.com/shakilaaulia/Dealan/shipment-service/repository"
)

type ShipmentService interface {
	Create(ctx context.Context, req domain.ShipmentRequest) (domain.ShipmentResponse, error)
	UploadProof(ctx context.Context, shipmentID string, proof domain.ProofData) error
}

type shipmentServiceImpl struct {
	repo repository.ShipmentRepository
}

func NewShipmentService(repo repository.ShipmentRepository) ShipmentService {
	return &shipmentServiceImpl{repo: repo}
}

func (s *shipmentServiceImpl) Create(ctx context.Context, req domain.ShipmentRequest) (domain.ShipmentResponse, error) {
	shipmentID, trackingCode, err := s.repo.CreateShipment(ctx, req)
	if err != nil {
		return domain.ShipmentResponse{}, err
	}

	labelURL, err := s.repo.GenerateLabel(ctx, shipmentID)
	if err != nil {
		return domain.ShipmentResponse{}, err
	}

	return domain.ShipmentResponse{
		ShipmentID:    shipmentID,
		KodeTracking:  trackingCode,
		LabelShipping: labelURL,
	}, nil
}

func (s *shipmentServiceImpl) UploadProof(ctx context.Context, shipmentID string, proof domain.ProofData) error {
	return s.repo.SaveProof(ctx, shipmentID, proof)
}

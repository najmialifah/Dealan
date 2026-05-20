package repository

import (
	"context"

	"github.com/shakilaaulia/Dealan/shipment-service/domain"
)

//go:generate mockgen -source=shipment_repository.go -destination=../mocks/mock_shipment_repository.go -package=mocks

type ShipmentRepository interface {
	CreateShipment(ctx context.Context, req domain.ShipmentRequest) (string, string, error) // returns shipmentID, trackingCode, error
	GenerateLabel(ctx context.Context, shipmentID string) (string, error)                   // returns labelURL
	SaveProof(ctx context.Context, shipmentID string, proof domain.ProofData) error
}

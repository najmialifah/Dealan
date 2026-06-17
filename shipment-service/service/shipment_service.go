package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/najmialifah/Dealan/shipment-service/domain"
	"github.com/najmialifah/Dealan/shipment-service/repository"
)

// ShipmentService mendefinisikan interface bisnis untuk operasi CRUD dan pelacakan pengiriman.
type ShipmentService interface {
	CreateShipment(ctx context.Context, req domain.ShipmentRequest) (domain.ShipmentResponse, error)
	GetShipmentByID(ctx context.Context, id string) (*domain.Shipment, error)
	GetShipmentByTrackingCode(ctx context.Context, code string) (*domain.Shipment, error)
	UpdateShipment(ctx context.Context, id string, req domain.ShipmentRequest) (*domain.Shipment, error)
	DeleteShipment(ctx context.Context, id string) error
	ListShipments(ctx context.Context) ([]domain.Shipment, error)
	UploadProof(ctx context.Context, id string, proof domain.ProofData) error
}

type shipmentServiceImpl struct {
	repo repository.ShipmentRepository
}

// NewShipmentService membuat instance baru dari implementasi ShipmentService.
func NewShipmentService(repo repository.ShipmentRepository) ShipmentService {
	return &shipmentServiceImpl{repo: repo}
}

// CreateShipment membuat pengiriman logistik baru dengan kode tracking unik dan menyimpan detail manifes dinamis ke JSONB.
func (s *shipmentServiceImpl) CreateShipment(ctx context.Context, req domain.ShipmentRequest) (domain.ShipmentResponse, error) {
	// 1. Satukan detail barang logistik ke dalam map manifes dinamis
	manifestMap := make(map[string]interface{})
	if req.KategoriBarang != "" {
		manifestMap["kategori"] = req.KategoriBarang
	}
	if req.BeratBarang > 0 {
		manifestMap["berat_kg"] = req.BeratBarang
	}
	if req.NamaPenerima != "" {
		manifestMap["nama_penerima"] = req.NamaPenerima
	}
	if req.NomorPenerima != "" {
		manifestMap["nomor_penerima"] = req.NomorPenerima
	}
	if req.CatatanPickup != "" {
		manifestMap["catatan"] = req.CatatanPickup
	}

	// Gabungkan data manifes dinamis tambahan kustom jika dikirimkan oleh klien
	for k, v := range req.Manifest {
		manifestMap[k] = v
	}

	manifestBytes, err := json.Marshal(manifestMap)
	if err != nil {
		return domain.ShipmentResponse{}, fmt.Errorf("gagal merubah data manifes ke JSON: %w", err)
	}

	// 2. Buat kode tracking acak (format: SHP-YYYYMMDD-XXXXX)
	t := time.Now()
	rand.Seed(t.UnixNano())
	randomDigits := rand.Intn(90000) + 10000 // 5 digit angka acak
	trackingCode := fmt.Sprintf("SHP-%s-%d", t.Format("20060102"), randomDigits)

	shipment := domain.Shipment{
		OrderID:      req.OrderID,
		TrackingCode: trackingCode,
		Status:       "pending",
		Manifest:     domain.JSONB(manifestBytes),
	}

	// 3. Simpan ke database PostgreSQL via GORM
	if err := s.repo.Create(ctx, &shipment); err != nil {
		return domain.ShipmentResponse{}, err
	}

	return domain.ShipmentResponse{
		ShipmentID:    shipment.ID,
		KodeTracking:  shipment.TrackingCode,
		LabelShipping: fmt.Sprintf("https://label.dealan.id/%s", shipment.TrackingCode),
	}, nil
}

func (s *shipmentServiceImpl) GetShipmentByID(ctx context.Context, id string) (*domain.Shipment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *shipmentServiceImpl) GetShipmentByTrackingCode(ctx context.Context, code string) (*domain.Shipment, error) {
	return s.repo.GetByTrackingCode(ctx, code)
}

// UpdateShipment memperbarui manifes pengiriman dinamis atau detail lainnya.
func (s *shipmentServiceImpl) UpdateShipment(ctx context.Context, id string, req domain.ShipmentRequest) (*domain.Shipment, error) {
	shipment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Perbarui manifes dinamis
	manifestMap := make(map[string]interface{})
	// Decode manifes lama terlebih dahulu
	if len(shipment.Manifest) > 0 {
		_ = json.Unmarshal(shipment.Manifest, &manifestMap)
	}

	// Tumpuk/update nilai dengan input baru jika ada
	if req.KategoriBarang != "" {
		manifestMap["kategori"] = req.KategoriBarang
	}
	if req.BeratBarang > 0 {
		manifestMap["berat_kg"] = req.BeratBarang
	}
	if req.NamaPenerima != "" {
		manifestMap["nama_penerima"] = req.NamaPenerima
	}
	if req.NomorPenerima != "" {
		manifestMap["nomor_penerima"] = req.NomorPenerima
	}
	if req.CatatanPickup != "" {
		manifestMap["catatan"] = req.CatatanPickup
	}

	// Tambahkan/update field kustom
	for k, v := range req.Manifest {
		manifestMap[k] = v
	}

	manifestBytes, err := json.Marshal(manifestMap)
	if err != nil {
		return nil, fmt.Errorf("gagal memproses update JSON manifes: %w", err)
	}

	shipment.Manifest = domain.JSONB(manifestBytes)
	shipment.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, shipment); err != nil {
		return nil, err
	}

	return shipment, nil
}

func (s *shipmentServiceImpl) DeleteShipment(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *shipmentServiceImpl) ListShipments(ctx context.Context) ([]domain.Shipment, error) {
	return s.repo.List(ctx)
}

// UploadProof menyimpan URL bukti foto serah terima barang dan merubah status menjadi 'delivered'.
func (s *shipmentServiceImpl) UploadProof(ctx context.Context, id string, proof domain.ProofData) error {
	shipment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	shipment.ProofImageURL = proof.ImageURL
	shipment.Status = "delivered" // Selesai terkirim
	shipment.UpdatedAt = time.Now()

	return s.repo.Update(ctx, shipment)
}

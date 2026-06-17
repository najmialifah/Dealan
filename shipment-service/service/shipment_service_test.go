package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/najmialifah/Dealan/shipment-service/domain"
	"github.com/najmialifah/Dealan/shipment-service/repository"
	"github.com/najmialifah/Dealan/shipment-service/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Menggunakan SQLite in-memory database untuk pengujian unit agar cepat dan terisolasi
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Gagal membuat in-memory SQLite database: %v", err)
	}

	err = db.AutoMigrate(&domain.Shipment{})
	if err != nil {
		t.Fatalf("Gagal melakukan auto-migrasi SQLite: %v", err)
	}

	return db
}

func TestShipmentService_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewShipmentRepository(db)
	srv := service.NewShipmentService(repo)

	ctx := context.Background()

	var createdShipmentID string
	var trackingCode string

	t.Run("✅ CreateShipment: Berhasil membuat pengiriman dengan manifes dinamis", func(t *testing.T) {
		req := domain.ShipmentRequest{
			OrderID:        "ORD-GOSEND-100",
			KategoriBarang: "makanan",
			BeratBarang:    1.5,
			NamaPenerima:   "Budi Santoso",
			NomorPenerima:  "08123456789",
			CatatanPickup:  "Hati-hati makanan hangat",
			Manifest: map[string]interface{}{
				"dimensi_cm": "30x20x15",
				"fragile":    false,
			},
		}

		res, err := srv.CreateShipment(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.ShipmentID)
		assert.NotEmpty(t, res.KodeTracking)
		assert.Contains(t, res.LabelShipping, res.KodeTracking)

		createdShipmentID = res.ShipmentID
		trackingCode = res.KodeTracking

		// Verifikasi data masuk ke DB dan format JSONB
		var shp domain.Shipment
		err = db.Where("id = ?", createdShipmentID).First(&shp).Error
		assert.NoError(t, err)
		assert.Equal(t, "ORD-GOSEND-100", shp.OrderID)
		assert.Equal(t, "pending", shp.Status)

		// Decode manifes JSONB untuk verifikasi data dinamis
		var manifest map[string]interface{}
		err = json.Unmarshal(shp.Manifest, &manifest)
		assert.NoError(t, err)
		assert.Equal(t, "makanan", manifest["kategori"])
		assert.Equal(t, 1.5, manifest["berat_kg"])
		assert.Equal(t, "Budi Santoso", manifest["nama_penerima"])
		assert.Equal(t, "30x20x15", manifest["dimensi_cm"])
		assert.Equal(t, false, manifest["fragile"])
	})

	t.Run("✅ Read: Mendapatkan detail pengiriman berdasarkan ID", func(t *testing.T) {
		shp, err := srv.GetShipmentByID(ctx, createdShipmentID)
		assert.NoError(t, err)
		assert.Equal(t, "ORD-GOSEND-100", shp.OrderID)
		assert.Equal(t, trackingCode, shp.TrackingCode)
	})

	t.Run("✅ Track: Mendapatkan detail pengiriman berdasarkan kode tracking", func(t *testing.T) {
		shp, err := srv.GetShipmentByTrackingCode(ctx, trackingCode)
		assert.NoError(t, err)
		assert.Equal(t, createdShipmentID, shp.ID)
	})

	t.Run("✅ Update: Memperbarui manifes barang secara dinamis", func(t *testing.T) {
		updateReq := domain.ShipmentRequest{
			OrderID:      "ORD-GOSEND-100",
			NamaPenerima: "Siti Rahma", // Ganti nama penerima
			Manifest: map[string]interface{}{
				"fragile":      true, // Ubah fragile menjadi true
				"suhu_celcius": 4,    // Tambahkan field baru di manifest dinamis
			},
		}

		updatedShp, err := srv.UpdateShipment(ctx, createdShipmentID, updateReq)
		assert.NoError(t, err)

		var manifest map[string]interface{}
		err = json.Unmarshal(updatedShp.Manifest, &manifest)
		assert.NoError(t, err)
		assert.Equal(t, "Siti Rahma", manifest["nama_penerima"])
		assert.Equal(t, true, manifest["fragile"])
		assert.Equal(t, float64(4), manifest["suhu_celcius"])
	})

	t.Run("✅ UploadProof: Mengunggah bukti pengiriman dan merubah status ke delivered", func(t *testing.T) {
		proof := domain.ProofData{
			ImageURL: "https://storage.dealan.id/proofs/shp-100.jpg",
		}

		err := srv.UploadProof(ctx, createdShipmentID, proof)
		assert.NoError(t, err)

		// Verifikasi status dan URL bukti di DB
		var shp domain.Shipment
		err = db.Where("id = ?", createdShipmentID).First(&shp).Error
		assert.NoError(t, err)
		assert.Equal(t, "delivered", shp.Status)
		assert.Equal(t, "https://storage.dealan.id/proofs/shp-100.jpg", shp.ProofImageURL)
	})

	t.Run("✅ Delete: Berhasil menghapus pengiriman", func(t *testing.T) {
		err := srv.DeleteShipment(ctx, createdShipmentID)
		assert.NoError(t, err)

		// Pastikan data tidak ditemukan lagi
		_, err = srv.GetShipmentByID(ctx, createdShipmentID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})
}

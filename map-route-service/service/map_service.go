package service

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math/rand"

	"map-route-service/models"
	"map-route-service/repository"
)

// MapService menyediakan kontrak fungsional untuk map-route-service.
type MapService interface {
	GetOrCreateRoute(ctx context.Context, req models.RouteRequest) (*models.RouteResponse, error)
}

type mapServiceImpl struct {
	repo repository.MapRepository
}

// NewMapService membuat instansiasi baru layanan MapService.
func NewMapService(r repository.MapRepository) MapService {
	return &mapServiceImpl{repo: r}
}

// GetOrCreateRoute mengambil rute tersimpan dari database atau membuat rute mock baru jika belum ada.
func (s *mapServiceImpl) GetOrCreateRoute(ctx context.Context, req models.RouteRequest) (*models.RouteResponse, error) {
	// Coba cari rute yang sudah tersimpan di database
	savedRoute, err := s.repo.GetRoute(ctx, req.Origin, req.Destination)
	if err == nil && savedRoute != nil {
		log.Printf("[MapService] Cache hit: Rute ditemukan di database untuk %s -> %s", req.Origin, req.Destination)
		return &models.RouteResponse{
			Origin:      savedRoute.Origin,
			Destination: savedRoute.Destination,
			Polyline:    savedRoute.Polyline,
			Distance:    savedRoute.Distance,
			Duration:    savedRoute.Duration,
		}, nil
	}

	log.Printf("[MapService] Cache miss: Membuat rute baru untuk %s -> %s", req.Origin, req.Destination)

	// Buat hash deterministik berdasarkan origin dan destination untuk membuat mock data yang konsisten
	hasher := sha256.New()
	hasher.Write([]byte(req.Origin + req.Destination))
	hashBytes := hasher.Sum(nil)
	seed := int64(binary.BigEndian.Uint64(hashBytes[:8]))

	// Seed randomizer lokal agar data yang dibuat deterministik bagi rute yang sama
	localRand := rand.New(rand.NewSource(seed))

	// Mock jarak (1.5 km sampai 20 km) dan durasi (kecepatan rata-rata 30 km/jam)
	distance := 1.5 + localRand.Float64()*18.5
	// Bulatkan jarak ke 2 desimal
	distance = float64(int(distance*100)) / 100
	// Durasi dalam detik (misal 2 menit per km)
	durationSeconds := int(distance * 120)

	// Membuat string polyline tiruan (encoded polyline format standar Google Maps)
	mockPolyline := generateMockPolyline(localRand, int(distance))

	// Simpan rute baru ini ke database
	newRoute := &models.MapRoute{
		Origin:      req.Origin,
		Destination: req.Destination,
		Polyline:    mockPolyline,
		Distance:    distance,
		Duration:    durationSeconds,
	}

	if err := s.repo.SaveRoute(ctx, newRoute); err != nil {
		log.Printf("[MapService-Warning] Gagal menyimpan rute jalan baru ke DB: %v", err)
	} else {
		log.Printf("[MapService] Rute baru berhasil disimpan ke database!")
	}

	return &models.RouteResponse{
		Origin:      newRoute.Origin,
		Destination: newRoute.Destination,
		Polyline:    newRoute.Polyline,
		Distance:    newRoute.Distance,
		Duration:    newRoute.Duration,
	}, nil
}

// generateMockPolyline menghasilkan string encoded polyline tiruan untuk simulasi visualisasi peta.
func generateMockPolyline(r *rand.Rand, steps int) string {
	// Karakter representatif encoded polyline
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	polyline := "_p~iF~ps|U" // Titik awal mock

	for i := 0; i < steps; i++ {
		// Buat segmen acak sepanjang 5 karakter
		segment := make([]byte, 5)
		for j := 0; j < 5; j++ {
			segment[j] = chars[r.Intn(len(chars))]
		}
		polyline += "@" + string(segment)
	}

	return polyline
}
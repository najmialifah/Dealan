package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"map-route-service/models"
)

// mockMapRepository adalah stub database untuk unit test map-route-service.
type mockMapRepository struct {
	route   *models.MapRoute
	getErr  error
	saveErr error
}

func (m *mockMapRepository) GetRoute(ctx context.Context, origin, destination string) (*models.MapRoute, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.route, nil
}

func (m *mockMapRepository) SaveRoute(ctx context.Context, route *models.MapRoute) error {
	return m.saveErr
}

// TestGetOrCreateRoute menguji logika pencarian dan pembuatan rute beserta polyline.
func TestGetOrCreateRoute(t *testing.T) {
	t.Run("✅ Cache hit (rute diambil dari database jika ada)", func(t *testing.T) {
		mockRoute := &models.MapRoute{
			Origin:      "Stasiun Gambir",
			Destination: "Monas",
			Polyline:    "_p~iF~ps|U@abcde",
			Distance:    1.8,
			Duration:    360,
		}
		repo := &mockMapRepository{route: mockRoute}
		svc := NewMapService(repo)

		res, err := svc.GetOrCreateRoute(context.Background(), models.RouteRequest{
			Origin:      "Stasiun Gambir",
			Destination: "Monas",
		})

		assert.NoError(t, err)
		assert.Equal(t, "Stasiun Gambir", res.Origin)
		assert.Equal(t, "Monas", res.Destination)
		assert.Equal(t, "_p~iF~ps|U@abcde", res.Polyline)
		assert.Equal(t, 1.8, res.Distance)
	})

	t.Run("✅ Cache miss (membuat rute baru deterministik dan menyimpannya)", func(t *testing.T) {
		repo := &mockMapRepository{getErr: errors.New("record not found")}
		svc := NewMapService(repo)

		res1, err := svc.GetOrCreateRoute(context.Background(), models.RouteRequest{
			Origin:      "Stasiun Gambir",
			Destination: "Monas",
		})

		assert.NoError(t, err)
		assert.Equal(t, "Stasiun Gambir", res1.Origin)
		assert.Equal(t, "Monas", res1.Destination)
		assert.NotEmpty(t, res1.Polyline)

		// Rute yang sama harus mengembalikan data yang deterministik
		res2, err := svc.GetOrCreateRoute(context.Background(), models.RouteRequest{
			Origin:      "Stasiun Gambir",
			Destination: "Monas",
		})
		assert.NoError(t, err)
		assert.Equal(t, res1.Distance, res2.Distance)
		assert.Equal(t, res1.Duration, res2.Duration)
		assert.Equal(t, res1.Polyline, res2.Polyline)
	})
}
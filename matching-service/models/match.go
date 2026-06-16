package models

// MatchRequest struktur request untuk mencari driver
type MatchRequest struct {
	OrderID   uint    `json:"order_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Radius    float64 `json:"radius,omitempty"` // Opsional, default 5000m di controller jika kosong
}

// MatchedDriver struktur response driver yang ditemukan
type MatchedDriver struct {
	DriverID  uint    `json:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance_meters"`
}

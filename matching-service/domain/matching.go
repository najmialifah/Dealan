package domain

type MatchingRequest struct {
	OrderID     string  `json:"order_id"`
	ServiceType string  `json:"service_type"` // ride/car/send
	LatUser     float64 `json:"lat_user"`     // lokasi user
	LongUser    float64 `json:"long_user"`
}

type MatchingResponse struct {
	DriverID      string `json:"driver_id"`      // driver terpilih
	EstimasiWaktu int    `json:"estimasi_waktu"` // waktu jemput
}

// Kontrak untuk Database/Service lain
type MatchingRepository interface {
	FindNearestDriver(req *MatchingRequest) (*MatchingResponse, error)
}
package domain

type PricingRequest struct {
	ServiceType string  `json:"service_type"` // ride, car, send
	Jarak       float64 `json:"jarak"`
	BeratBarang float64 `json:"berat_barang"`
}

type PricingResponse struct {
	Harga float64 `json:"harga"`
}

type NegotiationData struct {
	RequestedPrice float64 `json:"requested_price"`
}

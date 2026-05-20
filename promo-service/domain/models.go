package domain

type PromoRequest struct {
	Code string
}

type PromoResponse struct {
	Discount float64
	IsValid  bool
}
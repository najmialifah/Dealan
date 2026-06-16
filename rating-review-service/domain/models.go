package domain

// RatingRequest merepresentasikan data request untuk memberikan rating
type RatingRequest struct {
	OrderID      string `json:"order_id"`
	ReviewerID   string `json:"reviewer_id"`
	ReviewerRole string `json:"reviewer_role"` // 'user' atau 'driver'
	TargetID     string `json:"target_id"`
	TargetRole   string `json:"target_role"`   // 'user' atau 'driver'
	RatingScore  int    `json:"rating_score"`  // Skor rating antara 1-5
	Comment      string `json:"comment"`
}

// RatingResponse merepresentasikan data response setelah rating disubmit
type RatingResponse struct {
	AverageRating float64 `json:"average_rating"`
	ReviewID      string  `json:"review_id"`
	Status        string  `json:"status"`
}
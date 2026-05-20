package domain

type RatingRequest struct {
	OrderID      string `json:"order_id"`
	ReviewerID   string `json:"reviewer_id"`
	RatingScore  int    `json:"rating_score"` // 1-5
	Comment      string `json:"comment"`
}

type RatingResponse struct {
	AverageRating float64 `json:"average_rating"`
	ReviewID      string  `json:"review_id"`
	Status        string  `json:"status"`
}
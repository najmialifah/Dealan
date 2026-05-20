package domain

type PunishmentRequest struct {
	AccountID  string `json:"account_id"`
	ReasonCode string `json:"reason_code"`
	Duration   int    `json:"duration"` // dalam jam
}

type PunishmentResponse struct {
	PenaltyID        string `json:"penalty_id"`
	NewAccountStatus string `json:"new_account_status"` // Suspended / Banned
	Status           string `json:"status"`
}
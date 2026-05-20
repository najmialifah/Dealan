package domain


// input
type NotificationRequest struct {
	TargetID   string `json:"target_id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	ActionLink string `json:"action_link"`
}


// output
type NotificationResponse struct {
	DeliveryStatus string `json:"delivery_status"` // Sent / Failed
	MessageID      string `json:"message_id"`
}

// Interface untuk Mock
type NotificationProvider interface {
	Push(req NotificationRequest) (NotificationResponse, error)
}
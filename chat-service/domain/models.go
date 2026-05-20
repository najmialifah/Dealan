package domain

type ChatRequest struct {
	Message string
}

type ChatResponse struct {
	MessageID string
	Status    bool
}
package domain

type LocationUpdate struct {
	EntityID  string  `json:"entity_id"` // ID Driver atau User
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Bearing   float64 `json:"bearing"`   // arah hadap
}

// Kontrak untuk Database
type LocationRepository interface {
	SaveLocation(loc *LocationUpdate) error
}
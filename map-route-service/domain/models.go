package domain

type RouteRequest struct {
	Origin      string
	Destination string
	Mode        string
}

type RouteResponse struct {
	Distance float64
	Duration int
}
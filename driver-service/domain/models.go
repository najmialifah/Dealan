package domain

type DriverStatus string

const (
	StatusAvailable DriverStatus = "Available"
	StatusBusy      DriverStatus = "Busy"
	StatusOffline   DriverStatus = "Offline"
)

type Driver struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	VehicleNumber string       `json:"vehicle_number"`
	Lat           float64      `json:"lat"`
	Long          float64      `json:"long"`
	Status        DriverStatus `json:"status"`
	Rating        float64      `json:"rating"`
}

type UpdateLocationRequest struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type UpdateStatusRequest struct {
	Status DriverStatus `json:"status"`
}

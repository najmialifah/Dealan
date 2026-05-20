package domain

type ShipmentRequest struct {
	OrderID        string  `json:"order_id"`
	KategoriBarang string  `json:"kategori_barang"`
	BeratBarang    float64 `json:"berat_barang"`
	NamaPenerima   string  `json:"nama_penerima"`
	NomorPenerima  string  `json:"nomor_penerima"`
	CatatanPickup  string  `json:"catatan_pickup"`
}

type ShipmentResponse struct {
	ShipmentID    string `json:"shipment_id"`
	KodeTracking  string `json:"kode_tracking"`
	LabelShipping string `json:"label_shipping"`
}

type ProofData struct {
	ImageURL string `json:"image_url"`
}

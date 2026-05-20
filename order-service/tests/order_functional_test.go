package tests

import (
	"testing"
)

func TestCreateOrder_RealDB(t *testing.T) {
	// Karena database MySQL/Postgres kita belum dibikin, 
	// kita kasih tau Jenkins untuk nge-skip tes ini sementara waktu.
	t.Skip("SKIP DULU: Menunggu integrasi dengan database beneran selesai.")
	
	// Nanti kode untuk konek database dan nembak API beneran bakal ditulis di sini.
}
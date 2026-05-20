package service

import (
	"errors"
	"order-service/domain"
)

type OrderService struct {
	Repo domain.OrderRepository
}

// Fungsi utama buat bikin pesanan
func (s *OrderService) CreateOrder(order *domain.Order) error {
	// Logika sederhana: Cek apakah tipe layanan kosong?
	if order.ServiceType == "" {
		return errors.New("tipe layanan tidak boleh kosong")
	}
	
	// Kalau aman, simpan ke database
	return s.Repo.Save(order)
}
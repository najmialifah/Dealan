package models

import "time"

// ChatRoom mendefinisikan ruang obrolan per pesanan (order) antara user dan driver.
type ChatRoom struct {
	OrderID   string    `gorm:"type:varchar(50);primaryKey" json:"order_id"`
	UserID    string    `gorm:"type:varchar(50);not null" json:"user_id"`
	DriverID  string    `gorm:"type:varchar(50);not null" json:"driver_id"`
	Status    string    `gorm:"type:varchar(10);default:'active'" json:"status"` // active, closed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChatMessage menyimpan riwayat setiap pesan di database PostgreSQL.
type ChatMessage struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID    string    `gorm:"type:varchar(50);not null;index:idx_order_msg" json:"order_id"`
	SenderID   string    `gorm:"type:varchar(50);not null" json:"sender_id"`
	SenderRole string    `gorm:"type:varchar(10);not null" json:"sender_role"` // user | driver
	Message    string    `gorm:"type:text;not null" json:"message"`
	Type       string    `gorm:"type:varchar(20);default:'text'" json:"type"` // text, image, location_share
	SentAt     time.Time `gorm:"default:now()" json:"sent_at"`
	ReadStatus bool      `gorm:"default:false" json:"read_status"`
}

// WSMessage adalah struktur pesan yang dikirim/diterima melalui koneksi WebSocket.
type WSMessage struct {
	OrderID    string    `json:"order_id"`
	SenderID   string    `json:"sender_id"`
	SenderRole string    `json:"sender_role"` // user | driver
	Message    string    `json:"message"`
	Type       string    `json:"type"` // text, image, location_share
	SentAt     time.Time `json:"sent_at"`
}

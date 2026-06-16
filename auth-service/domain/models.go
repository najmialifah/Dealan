package domain

import "time"

// UserRole mendefinisikan tipe role pengguna
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleDriver UserRole = "driver"
	RoleAdmin  UserRole = "admin"
)

// AuthCredential merepresentasikan tabel auth_credentials untuk autentikasi GORM
type AuthCredential struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID    string    `gorm:"type:uuid;uniqueIndex;not null" json:"account_id"`
	Role         UserRole  `gorm:"type:varchar(10);not null" json:"role"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"` // Tambahan untuk mempermudah login via email
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:now()" json:"updated_at"`
}

// OTPCode merepresentasikan tabel otp_codes untuk kode OTP pendaftaran/login/reset password
type OTPCode struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	NomorHP   string    `gorm:"type:varchar(20);not null;index:idx_otp_nomor_hp" json:"nomor_hp"`
	OTPCode   string    `gorm:"type:varchar(6);not null" json:"otp_code"`
	Purpose   string    `gorm:"type:varchar(20);not null" json:"purpose"` // login, register, reset
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
}

// RefreshToken merepresentasikan tabel refresh_tokens untuk pengelolaan session
type RefreshToken struct {
	ID         string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID  string     `gorm:"type:uuid;not null;index:idx_rt_account_id" json:"account_id"`
	TokenHash  string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"token_hash"`
	DeviceInfo string     `gorm:"type:jsonb" json:"device_info"` // Menyimpan data device dalam format JSON
	ExpiresAt  time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	CreatedAt  time.Time  `gorm:"default:now()" json:"created_at"`
}

// LoginRequest adalah payload request untuk login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest adalah payload request untuk registrasi akun baru
type RegisterRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Role     UserRole `json:"role" binding:"required,oneof=user driver admin"`
	Nama     string   `json:"nama" binding:"required"`
	NomorHP  string   `json:"nomor_hp" binding:"required"`
}

// AuthResponse adalah payload response setelah berhasil login atau registrasi
type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	AccountID    string `json:"account_id"`
	Role         string `json:"role"`
}

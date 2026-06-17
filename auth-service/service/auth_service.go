package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/najmialifah/Dealan/auth-service/domain"
	"github.com/najmialifah/Dealan/auth-service/repository"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo          repository.AuthRepository
	producer      EventProducer
	jwtSecret     []byte
	jwtExpiration time.Duration
}

// Claims adalah struktur JWT Claims custom
type Claims struct {
	AccountID string `json:"account_id"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

// NewAuthService membuat instance baru dari service autentikasi
func NewAuthService(repo repository.AuthRepository, producer EventProducer, jwtSecret string, expiration time.Duration) AuthService {
	return &authService{
		repo:          repo,
		producer:      producer,
		jwtSecret:     []byte(jwtSecret),
		jwtExpiration: expiration,
	}
}

// Register menangani pendaftaran akun baru, hashing bcrypt, serta publish event Kafka
func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Cek apakah email sudah terdaftar
	existing, _ := s.repo.GetCredentialByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// Generate Account ID baru menggunakan UUID
	accountID := uuid.New().String()

	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Buat entitas AuthCredential
	cred := &domain.AuthCredential{
		ID:           uuid.New().String(),
		AccountID:    accountID,
		Role:         req.Role,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}

	// Simpan ke database menggunakan GORM
	err = s.repo.CreateCredential(ctx, cred)
	if err != nil {
		return nil, err
	}

	// Publish event ke Kafka secara asinkron tergantung Role
	if req.Role == domain.RoleUser {
		if s.producer != nil {
			_ = s.producer.PublishUserCreated(ctx, accountID, req.Nama, req.Email, req.NomorHP)
		}
	} else if req.Role == domain.RoleDriver {
		if s.producer != nil {
			_ = s.producer.PublishDriverCreated(ctx, accountID, req.Nama, req.Email, req.NomorHP)
		}
	}

	// Generate token JWT
	token, err := s.generateJWT(accountID, string(req.Role), req.Email)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:     token,
		AccountID: accountID,
		Role:      string(req.Role),
	}, nil
}

// Login menangani autentikasi kredensial email & password
func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	// Cari kredensial berdasarkan email
	cred, err := s.repo.GetCredentialByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("kredensial tidak valid")
	}

	if !cred.IsActive {
		return nil, errors.New("akun dinonaktifkan")
	}

	// Bandingkan password hash dengan plain text password menggunakan bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("kredensial tidak valid")
	}

	// Generate token JWT
	token, err := s.generateJWT(cred.AccountID, string(cred.Role), cred.Email)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:     token,
		AccountID: cred.AccountID,
		Role:      string(cred.Role),
	}, nil
}

// ValidateToken memvalidasi token JWT dan mengembalikan info kredensial
func (s *authService) ValidateToken(ctx context.Context, tokenStr string) (*domain.AuthCredential, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token tidak valid atau kedaluwarsa")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("gagal memproses claims token")
	}

	return &domain.AuthCredential{
		AccountID: claims.AccountID,
		Role:      domain.UserRole(claims.Role),
		Email:     claims.Email,
	}, nil
}

func (s *authService) generateJWT(accountID, role, email string) (string, error) {
	expirationTime := time.Now().Add(s.jwtExpiration)
	claims := &Claims{
		AccountID: accountID,
		Role:      role,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

package repository

import (
	"context"

	"github.com/najmialifah/Dealan/payment-service/domain"
	"gorm.io/gorm"
)

// PaymentRepository mendefinisikan kontrak untuk interaksi dengan PostgreSQL menggunakan GORM.
type PaymentRepository interface {
	// GetDB mengembalikan instance *gorm.DB utama untuk memulai transaksi (tx.Begin) di level service.
	GetDB() *gorm.DB
	CreateTransaction(ctx context.Context, tx *gorm.DB, transaction *domain.Transaction) error
	GetTransactionByTrxID(ctx context.Context, trxID string) (*domain.Transaction, error)
	UpdateTransaction(ctx context.Context, tx *gorm.DB, transaction *domain.Transaction) error
	GetDriverWallet(ctx context.Context, driverID string) (*domain.DriverWallet, error)
	GetDriverWalletWithLock(ctx context.Context, tx *gorm.DB, driverID string) (*domain.DriverWallet, error)
	UpdateDriverWallet(ctx context.Context, tx *gorm.DB, wallet *domain.DriverWallet) error
	CreateWalletTransaction(ctx context.Context, tx *gorm.DB, walletTx *domain.WalletTransaction) error
	CreatePaymentLog(ctx context.Context, tx *gorm.DB, log *domain.PaymentLog) error
	GetIdempotencyKey(ctx context.Context, key string) (*domain.IdempotencyKey, error)
	CreateIdempotencyKey(ctx context.Context, tx *gorm.DB, idem *domain.IdempotencyKey) error
}

type paymentRepositoryImpl struct {
	db *gorm.DB
}

// NewPaymentRepository membuat instance baru dari implementasi PaymentRepository.
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepositoryImpl{db: db}
}

func (r *paymentRepositoryImpl) GetDB() *gorm.DB {
	return r.db
}

func (r *paymentRepositoryImpl) CreateTransaction(ctx context.Context, tx *gorm.DB, transaction *domain.Transaction) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(transaction).Error
}

func (r *paymentRepositoryImpl) GetTransactionByTrxID(ctx context.Context, trxID string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.WithContext(ctx).Where("transaction_id = ?", trxID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *paymentRepositoryImpl) UpdateTransaction(ctx context.Context, tx *gorm.DB, transaction *domain.Transaction) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Save(transaction).Error
}

func (r *paymentRepositoryImpl) GetDriverWallet(ctx context.Context, driverID string) (*domain.DriverWallet, error) {
	var wallet domain.DriverWallet
	err := r.db.WithContext(ctx).Where("driver_id = ?", driverID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *paymentRepositoryImpl) GetDriverWalletWithLock(ctx context.Context, tx *gorm.DB, driverID string) (*domain.DriverWallet, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	var wallet domain.DriverWallet
	// Menggunakan SELECT ... FOR UPDATE untuk mengunci baris dompet driver agar tidak terjadi race condition
	err := db.WithContext(ctx).Set("gorm:query_option", "FOR UPDATE").Where("driver_id = ?", driverID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *paymentRepositoryImpl) UpdateDriverWallet(ctx context.Context, tx *gorm.DB, wallet *domain.DriverWallet) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Save(wallet).Error
}

func (r *paymentRepositoryImpl) CreateWalletTransaction(ctx context.Context, tx *gorm.DB, walletTx *domain.WalletTransaction) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(walletTx).Error
}

func (r *paymentRepositoryImpl) CreatePaymentLog(ctx context.Context, tx *gorm.DB, log *domain.PaymentLog) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(log).Error
}

func (r *paymentRepositoryImpl) GetIdempotencyKey(ctx context.Context, key string) (*domain.IdempotencyKey, error) {
	var idem domain.IdempotencyKey
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&idem).Error
	if err != nil {
		return nil, err
	}
	return &idem, nil
}

func (r *paymentRepositoryImpl) CreateIdempotencyKey(ctx context.Context, tx *gorm.DB, idem *domain.IdempotencyKey) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(idem).Error
}

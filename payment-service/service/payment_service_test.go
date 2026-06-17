package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/najmialifah/Dealan/payment-service/domain"
	"github.com/najmialifah/Dealan/payment-service/repository"
	"github.com/najmialifah/Dealan/payment-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockKafkaProducer adalah mock untuk KafkaProducer.
type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) PublishPaymentCompleted(ctx context.Context, trxID, orderID, status string, driverEarnings float64) error {
	args := m.Called(ctx, trxID, orderID, status, driverEarnings)
	return args.Error(0)
}

func setupTestDB(t *testing.T) *gorm.DB {
	// Menggunakan SQLite in-memory database untuk pengujian unit agar cepat dan terisolasi
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Gagal membuat in-memory SQLite database: %v", err)
	}

	err = db.AutoMigrate(
		&domain.Transaction{},
		&domain.DriverWallet{},
		&domain.WalletTransaction{},
		&domain.PaymentLog{},
		&domain.IdempotencyKey{},
	)
	if err != nil {
		t.Fatalf("Gagal melakukan auto-migrasi SQLite: %v", err)
	}

	return db
}

func TestPaymentService_Process(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPaymentRepository(db)
	mockKafka := new(MockKafkaProducer)
	srv := service.NewPaymentService(repo, mockKafka)

	ctx := context.Background()

	t.Run("✅ Berhasil membuat transaksi pembayaran baru", func(t *testing.T) {
		req := domain.PaymentRequest{
			OrderID:          "ORD-001",
			Nominal:          25000,
			MetodePembayaran: "qris",
			UserID:           "b58d0426-8c46-4c7c-b391-a185bbf90f05",
			DriverID:         "cfc846bf-efb0-4dbf-8181-79ef88be4fbf",
			IdempotencyKey:   "idem-key-1",
		}

		res, err := srv.Process(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.TransactionID)
		assert.Equal(t, "PENDING", res.PaymentStatus)

		// Verifikasi transaksi disimpan di DB
		var trx domain.Transaction
		err = db.Where("transaction_id = ?", res.TransactionID).First(&trx).Error
		assert.NoError(t, err)
		assert.Equal(t, "ORD-001", trx.OrderID)
		assert.Equal(t, 25000.0, trx.GrossAmount)
		assert.Equal(t, 20000.0, trx.DriverEarnings) // 80% dari gross
		assert.Equal(t, 5000.0, trx.PlatformFee)      // 20% dari gross
	})

	t.Run("✅ Idempotensi: Request dengan key sama tidak diproses ulang", func(t *testing.T) {
		req := domain.PaymentRequest{
			OrderID:          "ORD-001",
			Nominal:          25000,
			MetodePembayaran: "qris",
			UserID:           "b58d0426-8c46-4c7c-b391-a185bbf90f05",
			DriverID:         "cfc846bf-efb0-4dbf-8181-79ef88be4fbf",
			IdempotencyKey:   "idem-key-2",
		}

		// Panggilan Pertama
		res1, err := srv.Process(ctx, req)
		assert.NoError(t, err)

		// Panggilan Kedua dengan Key yang sama
		res2, err := srv.Process(ctx, req)
		assert.NoError(t, err)

		// Response harus identik
		assert.Equal(t, res1.TransactionID, res2.TransactionID)
		assert.Equal(t, res1.PaymentStatus, res2.PaymentStatus)

		// Pastikan record di database untuk transaksi ID tersebut hanya ada 1
		var count int64
		db.Model(&domain.Transaction{}).Where("order_id = ?", "ORD-001").Count(&count)
		// 1 dari testcase sebelumnya (ORD-001) + 1 dari testcase ini (tidak duplikat) = 2
		assert.Equal(t, int64(2), count)
	})
}

func TestPaymentService_ProcessWebhook(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPaymentRepository(db)
	mockKafka := new(MockKafkaProducer)
	srv := service.NewPaymentService(repo, mockKafka)

	ctx := context.Background()

	// Inisialisasi awal transaksi status pending
	trx := domain.Transaction{
		TransactionID:    "TRX-WEBHOOK-TEST",
		OrderID:          "ORD-100",
		UserID:           "b58d0426-8c46-4c7c-b391-a185bbf90f05",
		DriverID:         "cfc846bf-efb0-4dbf-8181-79ef88be4fbf",
		GrossAmount:      50000,
		NetAmount:        50000,
		DriverEarnings:   40000, // 80%
		PlatformFee:      10000, // 20%
		MetodePembayaran: "qris",
		Status:           "pending",
	}
	db.Create(&trx)

	t.Run("✅ Berhasil memproses webhook sukses & tambah saldo driver", func(t *testing.T) {
		// Mock Kafka publisher sukses
		mockKafka.On("PublishPaymentCompleted", mock.Anything, "TRX-WEBHOOK-TEST", "ORD-100", "success", 40000.0).Return(nil).Once()

		err := srv.ProcessWebhook(ctx, "TRX-WEBHOOK-TEST", "success")
		assert.NoError(t, err)

		// Cek status transaksi terupdate
		var updatedTrx domain.Transaction
		err = db.Where("transaction_id = ?", "TRX-WEBHOOK-TEST").First(&updatedTrx).Error
		assert.NoError(t, err)
		assert.Equal(t, "success", updatedTrx.Status)
		assert.NotNil(t, updatedTrx.PaidAt)

		// Cek saldo dompet driver bertambah
		var wallet domain.DriverWallet
		err = db.Where("driver_id = ?", trx.DriverID).First(&wallet).Error
		assert.NoError(t, err)
		assert.Equal(t, 40000.0, wallet.Balance)
		assert.Equal(t, 40000.0, wallet.TotalEarned)

		// Cek mutasi wallet
		var walletTx domain.WalletTransaction
		err = db.Where("driver_id = ?", trx.DriverID).First(&walletTx).Error
		assert.NoError(t, err)
		assert.Equal(t, "credit", walletTx.Type)
		assert.Equal(t, 40000.0, walletTx.Amount)
		assert.Equal(t, "TRX-WEBHOOK-TEST", walletTx.RefID)

		mockKafka.AssertExpectations(t)
	})

	t.Run("✅ Webhook sukses bersifat Idempoten (pemrosesan berulang tidak menambah saldo)", func(t *testing.T) {
		// Pemanggilan ulang webhook untuk transaksi yang sudah success
		err := srv.ProcessWebhook(ctx, "TRX-WEBHOOK-TEST", "success")
		assert.NoError(t, err)

		// Saldo dompet driver harus tetap 40000 (tidak menjadi 80000)
		var wallet domain.DriverWallet
		err = db.Where("driver_id = ?", trx.DriverID).First(&wallet).Error
		assert.NoError(t, err)
		assert.Equal(t, 40000.0, wallet.Balance)
	})

	t.Run("❌ Webhook error saat transaksi tidak ditemukan", func(t *testing.T) {
		err := srv.ProcessWebhook(ctx, "TRX-UNKNOWN", "success")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})
}

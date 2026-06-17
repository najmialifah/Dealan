package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/najmialifah/Dealan/payment-service/domain"
	"github.com/najmialifah/Dealan/payment-service/repository"
	"gorm.io/gorm"
)

// KafkaProducer mendefinisikan interface untuk mempublikasikan event ke Kafka.
type KafkaProducer interface {
	PublishPaymentCompleted(ctx context.Context, trxID, orderID, status string, driverEarnings float64) error
}

// PaymentService mendefinisikan logika bisnis untuk memproses pembayaran dan mengelola dompet driver.
type PaymentService interface {
	Process(ctx context.Context, req domain.PaymentRequest) (domain.PaymentResponse, error)
	ProcessWebhook(ctx context.Context, trxID string, status string) error
	GetStatus(ctx context.Context, transactionID string) (domain.PaymentResponse, error)
	GetDriverWallet(ctx context.Context, driverID string) (*domain.DriverWallet, error)
}

type paymentServiceImpl struct {
	repo          repository.PaymentRepository
	kafkaProducer KafkaProducer
}

// NewPaymentService membuat instance baru dari PaymentService.
func NewPaymentService(repo repository.PaymentRepository, kafkaProducer KafkaProducer) PaymentService {
	return &paymentServiceImpl{
		repo:          repo,
		kafkaProducer: kafkaProducer,
	}
}

// Process membuat transaksi pembayaran baru dengan jaminan Idempotensi dan GORM Transaction (ACID).
func (s *paymentServiceImpl) Process(ctx context.Context, req domain.PaymentRequest) (domain.PaymentResponse, error) {
	// 1. Cek Idempotensi untuk menghindari double payment jika key disediakan
	if req.IdempotencyKey != "" {
		idem, err := s.repo.GetIdempotencyKey(ctx, req.IdempotencyKey)
		if err == nil && idem != nil {
			var cachedRes domain.PaymentResponse
			if errJson := json.Unmarshal([]byte(idem.Response), &cachedRes); errJson == nil {
				log.Printf("[Idempotency] Transaksi duplikat dideteksi untuk key: %s, mengembalikan cache response\n", req.IdempotencyKey)
				return cachedRes, nil
			}
		}
	}

	// 2. Buat ID transaksi unik secara dinamis (format: TRX-YYYYMMDD-XXXXX)
	t := time.Now()
	rand.Seed(t.UnixNano())
	randomDigits := rand.Intn(90000) + 10000 // 5 digit angka acak
	trxID := fmt.Sprintf("TRX-%s-%d", t.Format("20060102"), randomDigits)

	// Hitung fee platform (20%) dan pembagian driver (80%)
	platformFee := req.Nominal * 0.20
	driverEarnings := req.Nominal * 0.80

	transaction := domain.Transaction{
		TransactionID:    trxID,
		OrderID:          req.OrderID,
		UserID:           req.UserID,
		DriverID:         req.DriverID,
		GrossAmount:      req.Nominal,
		DiscountAmount:   0,
		NetAmount:        req.Nominal,
		MetodePembayaran: req.MetodePembayaran,
		Status:           "pending",
		PlatformFee:      platformFee,
		DriverEarnings:   driverEarnings,
	}

	// 3. Mulai Transaksi GORM untuk menjamin sifat ACID agar tidak terjadi double payment
	tx := s.repo.GetDB().Begin()
	if tx.Error != nil {
		return domain.PaymentResponse{}, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Simpan transaksi
	if err := s.repo.CreateTransaction(ctx, tx, &transaction); err != nil {
		tx.Rollback()
		return domain.PaymentResponse{}, err
	}

	// Simpan log audit
	payloadBytes, _ := json.Marshal(req)
	logEntry := domain.PaymentLog{
		TransactionID: trxID,
		Event:         "created",
		Payload:       string(payloadBytes),
	}
	if err := s.repo.CreatePaymentLog(ctx, tx, &logEntry); err != nil {
		tx.Rollback()
		return domain.PaymentResponse{}, err
	}

	res := domain.PaymentResponse{
		TransactionID: trxID,
		PaymentStatus: "PENDING",
		InvoiceURL:    fmt.Sprintf("https://checkout.dealan.id/%s", trxID),
	}

	// Simpan key idempotensi bersama response
	if req.IdempotencyKey != "" {
		resBytes, _ := json.Marshal(res)
		idemKey := domain.IdempotencyKey{
			Key:      req.IdempotencyKey,
			Response: string(resBytes),
		}
		if err := s.repo.CreateIdempotencyKey(ctx, tx, &idemKey); err != nil {
			tx.Rollback()
			return domain.PaymentResponse{}, err
		}
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return domain.PaymentResponse{}, err
	}

	return res, nil
}

// ProcessWebhook memproses notifikasi dari Payment Gateway (e.g. Midtrans/Xendit) untuk merubah status transaksi.
func (s *paymentServiceImpl) ProcessWebhook(ctx context.Context, trxID string, status string) error {
	// 1. Mulai Transaksi GORM
	tx := s.repo.GetDB().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 2. Ambil data transaksi dengan status saat ini
	transaction, err := s.repo.GetTransactionByTrxID(ctx, trxID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Jika status sudah diproses (success/failed), lewati proses untuk memastikan sifat idempoten webhook
	if transaction.Status == "success" || transaction.Status == "failed" {
		tx.Commit()
		log.Printf("[Webhook] Transaksi %s sudah diproses sebelumnya dengan status %s\n", trxID, transaction.Status)
		return nil
	}

	// 3. Update status transaksi
	transaction.Status = status
	now := time.Now()
	if status == "success" {
		transaction.PaidAt = &now

		// Update saldo dompet driver jika pembayaran non-tunai (QRIS/Bank) berhasil
		if transaction.MetodePembayaran == "qris" || transaction.MetodePembayaran == "bank" {
			// Dapatkan dompet driver dengan lock (SELECT ... FOR UPDATE) untuk mencegah race condition
			wallet, err := s.repo.GetDriverWalletWithLock(ctx, tx, transaction.DriverID)
			if err != nil {
				// Jika dompet belum ada, buat baru
				if errors.Is(err, gorm.ErrRecordNotFound) {
					wallet = &domain.DriverWallet{
						DriverID:    transaction.DriverID,
						Balance:     0,
						TotalEarned: 0,
					}
					// GORM AutoMigrate menjamin tabel sudah ada
					if errCreate := tx.Create(wallet).Error; errCreate != nil {
						tx.Rollback()
						return errCreate
					}
				} else {
					tx.Rollback()
					return err
				}
			}

			// Tambahkan bagi hasil ke dompet driver (80% dari gross amount)
			wallet.Balance += transaction.DriverEarnings
			wallet.TotalEarned += transaction.DriverEarnings
			wallet.UpdatedAt = now

			if errUpdate := s.repo.UpdateDriverWallet(ctx, tx, wallet); errUpdate != nil {
				tx.Rollback()
				return errUpdate
			}

			// Buat riwayat mutasi dompet
			walletTx := domain.WalletTransaction{
				DriverID:     transaction.DriverID,
				Type:         "credit",
				Amount:       transaction.DriverEarnings,
				BalanceAfter: wallet.Balance,
				RefID:        transaction.TransactionID,
				Note:         fmt.Sprintf("Pembayaran order %s berhasil", transaction.OrderID),
			}
			if errWTx := s.repo.CreateWalletTransaction(ctx, tx, &walletTx); errWTx != nil {
				tx.Rollback()
				return errWTx
			}
		}
	}

	// Simpan update transaksi
	if errSave := s.repo.UpdateTransaction(ctx, tx, transaction); errSave != nil {
		tx.Rollback()
		return errSave
	}

	// Simpan log audit webhook
	logEntry := domain.PaymentLog{
		TransactionID: trxID,
		Event:         "webhook_received",
		Payload:       fmt.Sprintf(`{"status": "%s"}`, status),
	}
	if errLog := s.repo.CreatePaymentLog(ctx, tx, &logEntry); errLog != nil {
		tx.Rollback()
		return errLog
	}

	// Commit transaksi
	if errCommit := tx.Commit().Error; errCommit != nil {
		return errCommit
	}

	// 4. Publish Event "PAYMENT_COMPLETED" ke Kafka jika transaksi sukses
	if status == "success" {
		errKafka := s.kafkaProducer.PublishPaymentCompleted(ctx, transaction.TransactionID, transaction.OrderID, status, transaction.DriverEarnings)
		if errKafka != nil {
			log.Printf("[Kafka] Gagal mempublikasikan event PAYMENT_COMPLETED untuk transaksi %s: %v\n", trxID, errKafka)
			// Catatan: Di microservices, kegagalan Kafka biasanya tidak membatalkan transaksi DB (eventual consistency), cukup dicatat.
		} else {
			log.Printf("[Kafka] Berhasil mempublikasikan event PAYMENT_COMPLETED untuk transaksi %s\n", trxID)
		}
	}

	return nil
}

// GetStatus mengambil status pembayaran berdasarkan transaction_id.
func (s *paymentServiceImpl) GetStatus(ctx context.Context, transactionID string) (domain.PaymentResponse, error) {
	trx, err := s.repo.GetTransactionByTrxID(ctx, transactionID)
	if err != nil {
		return domain.PaymentResponse{}, err
	}
	return domain.PaymentResponse{
		TransactionID: trx.TransactionID,
		PaymentStatus: trx.Status,
		InvoiceURL:    fmt.Sprintf("https://checkout.dealan.id/%s", trx.TransactionID),
	}, nil
}

// GetDriverWallet mengambil data dompet driver terkini.
func (s *paymentServiceImpl) GetDriverWallet(ctx context.Context, driverID string) (*domain.DriverWallet, error) {
	return s.repo.GetDriverWallet(ctx, driverID)
}

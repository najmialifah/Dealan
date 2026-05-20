package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	// Driver postgres nanti di-uncomment kalau sudah pakai database beneran
	// _ "github.com/lib/pq" 
)

func TestNotificationFunctional_Send_Success(t *testing.T) {
	// 1. Setup Database Connection (Gunakan nama DB notification_db_test)
	dbDSN := "postgres://user:pass@localhost:5432/notification_db_test?sslmode=disable"
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		t.Skipf("Skipping functional test: Database not accessible: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2. Cek apakah database benar-benar bisa dihubungi
	if err := db.PingContext(ctx); err != nil {
		t.Skipf("Skipping functional test: Failed to ping database: %v", err)
	}

	t.Log("Notification Functional test structure is correct. Business logic is draft, expecting failure.")

	t.Fail() 
}
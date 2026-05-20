package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	// _ "github.com/lib/pq" // Driver postgres
)

func TestPunishmentFunctional_Apply_Success(t *testing.T) {
	// 1. Setup Database Connection (Gunakan nama DB punishment_db_test)
	dbDSN := "postgres://user:pass@localhost:5432/punishment_db_test?sslmode=disable"
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		t.Skipf("Skipping functional test: Database not accessible: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2. Simulasi Ping seperti di Auth Service temanmu
	if err := db.PingContext(ctx); err != nil {
		t.Skipf("Skipping functional test: Failed to ping database: %v", err)
	}

	t.Log("Punishment Functional test structure is correct. Business logic is draft, expecting failure.")
	
	// Sengaja digagalkan karena masih tahap draft
	t.Fail() 
}
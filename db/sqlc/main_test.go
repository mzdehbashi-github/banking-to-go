package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://postgres:password@localhost:5432/banking?sslmode=disable"
)

// Will be initialized by func `TestMain`
var conn *sql.DB

// var dbs *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	conn, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can not connect to DB", err)
	}

	os.Exit(m.Run())
}

// withTransaction executes the test function, using a database transaction,
// and finally rolls back the transaction
func withTransaction(t *testing.T, fn func(*Queries)) {
	tx, err := conn.Begin()

	if err != nil {
		t.Fatalf("Error beginning transaction: %v/n", err)
	}

	// Rollback changes at the end of the test
	defer tx.Rollback()

	queries := New(tx)

	// Execute the function with the transaction
	fn(queries)
}

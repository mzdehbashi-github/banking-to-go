package db

import (
	"database/sql"
	"gopsql/banking/util"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// Will be initialized by func `TestMain`
var conn *sql.DB

// var dbs *pgxpool.Pool

func TestMain(m *testing.M) {
	config := util.LoadConfig()
	var err error
	conn, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to DB", err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxIdleTime(2 * time.Second)
	conn.SetConnMaxLifetime(5 * time.Second)
	log.Println("conn: ", conn)

	os.Exit(m.Run())
}

// withTransaction executes the test function, using a database transaction,
// and finally rolls back the transaction
func withTransaction(t *testing.T, fn func(Querier)) {
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

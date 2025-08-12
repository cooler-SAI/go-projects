package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

type Account struct {
	ID      int
	Name    string
	Balance float64
}

// setupDatabase initializes the database and table.
func setupDatabase(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS accounts;")
	if err != nil {
		log.Fatalf("Failed to drop table: %v", err)
	}

	createTableSQL := `
	CREATE TABLE accounts (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		balance NUMERIC(10, 2)
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	insertSQL := `INSERT INTO accounts (name, balance) VALUES ('Alice', 1000.00);`
	_, err = db.Exec(insertSQL)
	if err != nil {
		log.Fatalf("Failed to insert initial data: %v", err)
	}
	fmt.Println("Database and initial data set up successfully.")
}

// readWithIsolationLevel demonstrates how a transaction sees data based on its isolation level.
func readWithIsolationLevel(db *sql.DB, isolation sql.IsolationLevel, wg *sync.WaitGroup) {
	defer wg.Done()

	// Define the context and begin a transaction.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Begin a transaction with the specified isolation level.
	// We use the TxOptions struct to configure the transaction's isolation level.
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: isolation})
	if err != nil {
		log.Printf("Failed to begin transaction with isolation level %s: %v", isolation, err)
		return
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Printf("Failed to rollback transaction: %v", err)
		}
	}(tx) // Ensure rollback happens if something goes wrong.
}

func main() {

	fmt.Println("Starting Isolation Level Demonstration...")

}

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

	fmt.Printf("\n--- Reader Transaction started with Isolation Level: %s ---\n", isolation)

	// First read: Read the initial balance.
	var initialBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE name = 'Alice';").Scan(&initialBalance)
	if err != nil {
		log.Printf("Reader failed to read initial balance: %v", err)
		return
	}
	fmt.Printf("Initial read (before writer's commit): Alice's balance is %.2f\n", initialBalance)

	// Wait for a moment to allow the writer transaction to do its work.
	time.Sleep(2 * time.Second)

	// Second read: Read the balance again.
	var secondBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE name = 'Alice';").Scan(&secondBalance)
	if err != nil {
		log.Printf("Reader failed to read second balance: %v", err)
		return
	}
	fmt.Printf("Second read (after writer's commit): Alice's balance is %.2f\n", secondBalance)

	// Check if the two reads are different (a "non-repeatable read" anomaly).
	if initialBalance != secondBalance {
		fmt.Printf("--- Anomaly Detected: Non-repeatable read! Initial balance %.2f is different from second balance %.2f ---\n", initialBalance, secondBalance)
	} else {
		fmt.Printf("--- No Anomaly: The balance remained consistent throughout the transaction. ---\n")
	}

	// Commit the reader transaction.
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit reader transaction: %v", err)
	}
}

// writerTransaction updates the balance and commits.
func writerTransaction(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Start a transaction.
	tx, err := db.BeginTx(ctx, nil) // Use default isolation level.
	if err != nil {
		log.Printf("Failed to begin writer transaction: %v", err)
		return
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Printf("Failed to rollback writer transaction: %v", err)
		}
	}(tx) // Rollback on error.

	// Update Alice's balance.
	_, err = tx.Exec("UPDATE accounts SET balance = balance + 500.00 WHERE name = 'Alice';")
	if err != nil {
		log.Printf("Writer failed to update balance: %v", err)
		return
	}
	fmt.Println("\nWriter Transaction: Alice's balance updated to 1500.00 (but not yet committed).")

	// Wait for a moment to ensure the reader has already done its first read.
	time.Sleep(1 * time.Second)

	// Commit the writer transaction.
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit writer transaction: %v", err)
		return
	}
	fmt.Println("Writer Transaction: Committed the balance change.")
}

func main() {

	fmt.Println("Starting Isolation Level Demonstration...")

	// Connect to the PostgreSQL database.
	connStr := "user=postgres password=example host=localhost port=5432 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Failed to close the database connection: %v", err)
		}
	}(db)
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)

	}
	// Setup DB for a clean test
	setupDatabase(db)

	var wg sync.WaitGroup
	// --- 1. Demonstrate READ COMMITTED (default) ---
	fmt.Println("\n=================================================")
	fmt.Println("RUNNING DEMO with Isolation Level: READ COMMITTED")
	fmt.Println("=================================================")

	wg.Add(2)
	// Launch the reader and writer goroutines.
	go readWithIsolationLevel(db, sql.LevelReadCommitted, &wg)
	go writerTransaction(db, &wg)
	wg.Wait()

	// Setup database again for the next test.
	setupDatabase(db)

	// --- 2. Demonstrate REPEATABLE READ ---
	fmt.Println("\n=================================================")
	fmt.Println("RUNNING DEMO with Isolation Level: REPEATABLE READ")
	fmt.Println("=================================================")

	wg.Add(2)
	go readWithIsolationLevel(db, sql.LevelRepeatableRead, &wg)
	go writerTransaction(db, &wg)
	wg.Wait()

	fmt.Println("\nDemonstration completed.")

}

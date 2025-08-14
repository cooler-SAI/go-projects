package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

// Account - a simple struct for our database table.
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
// It now uses two channels for strict synchronization with the writer goroutine.
func readWithIsolationLevel(db *sql.DB, isolation sql.IsolationLevel, writerReady chan struct{}, writerDone chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: isolation})
	if err != nil {
		log.Printf("Failed to begin transaction with isolation level %s: %v", isolation, err)
		return
	}
	defer func() {
		if err := tx.Rollback(); !errors.Is(err, sql.ErrTxDone) && nil != err {
			log.Printf("readWithIsolationLevel: failed to rollback: %v", err)
		}
	}()

	fmt.Printf("\n--- Reader Transaction started with Isolation Level: %s ---\n", isolation)

	// First read: Read the initial balance.
	var initialBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE name = 'Alice';").Scan(&initialBalance)
	if err != nil {
		log.Printf("Reader failed to read initial balance: %v", err)
		return
	}
	fmt.Printf("Initial read (before writer's commit): Alice's balance is %.2f\n", initialBalance)

	// Signal to the writer that the first read is complete.
	close(writerReady)

	// Wait for the writer to commit its changes.
	<-writerDone

	// Second read: Read the balance again.
	var secondBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE name = 'Alice';").Scan(&secondBalance)
	if err != nil {
		log.Printf("Reader failed to read second balance: %v", err)
		return
	}
	fmt.Printf("Second read (after writer's commit): Alice's balance is %.2f\n", secondBalance)

	// Check if the two reads differ (a "non-repeatable read" anomaly).
	if initialBalance != secondBalance {
		fmt.Printf("--- Anomaly Detected: Non-repeatable read! Initial balance %.2f is different from second balance %.2f ---\n", initialBalance, secondBalance)
	} else {
		fmt.Printf("--- No Anomaly: The balance remained consistent throughout the transaction. ---\n")
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit reader transaction: %v", err)
	}
}

// writerTransaction updates the balance and commits the changes.
// It now uses two channels for strict synchronization with the reader.
func writerTransaction(db *sql.DB, writerReady chan struct{}, writerDone chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin writer transaction: %v", err)
		return
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("writerTransaction: failed to rollback: %v", err)
		}
	}()

	// Wait for the reader to perform its first read.
	<-writerReady

	// Update Alice's balance.
	_, err = tx.Exec("UPDATE accounts SET balance = balance + 500.00 WHERE name = 'Alice';")
	if err != nil {
		log.Printf("Writer failed to update balance: %v", err)
		return
	}
	fmt.Println("\nWriter Transaction: Alice's balance updated to 1500.00 (but not yet committed).")

	// Commit the writer's transaction.
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit writer transaction: %v", err)
		return
	}
	fmt.Println("Writer Transaction: Committed the balance change.")

	// Signal to the reader that the writer's work is done.
	close(writerDone)
}

func main() {
	fmt.Println("Starting Isolation Level Demonstration...")

	connStr := "user=postgres password=example host=localhost port=5432 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("main: failed to close db: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	setupDatabase(db)

	var wg sync.WaitGroup

	// --- 1. Demonstrate READ COMMITTED (default) ---
	fmt.Println("\n=================================================")
	fmt.Println("RUNNING DEMO with Isolation Level: READ COMMITTED")
	fmt.Println("=================================================")

	writerReady1 := make(chan struct{})
	writerDone1 := make(chan struct{})
	wg.Add(2)
	go readWithIsolationLevel(db, sql.LevelReadCommitted, writerReady1, writerDone1, &wg)
	go writerTransaction(db, writerReady1, writerDone1, &wg)
	wg.Wait()

	setupDatabase(db)

	// --- 2. Demonstrate REPEATABLE READ ---
	fmt.Println("\n=================================================")
	fmt.Println("RUNNING DEMO with Isolation Level: REPEATABLE READ")
	fmt.Println("=================================================")

	writerReady2 := make(chan struct{})
	writerDone2 := make(chan struct{})
	wg.Add(2)
	go readWithIsolationLevel(db, sql.LevelRepeatableRead, writerReady2, writerDone2, &wg)
	go writerTransaction(db, writerReady2, writerDone2, &wg)
	wg.Wait()

	fmt.Println("\nDemonstration completed.")
}

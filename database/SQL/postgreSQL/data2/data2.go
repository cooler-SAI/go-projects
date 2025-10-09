package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

type Account struct {
	ID      int
	Name    string
	Balance float64
}

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
	fmt.Println("Database and initial data set up successfully (Alice's balance: 1000.00).")
}

func readWithIsolationLevel(db *sql.DB, isolation sql.IsolationLevel, writerReady chan struct{}, writerDone chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: isolation})
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return
	}

	fmt.Printf("\n--- [Reader] Transaction started with Isolation Level: %s ---\n", isolation)

	// First read
	var initialBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE name = 'Alice';").Scan(&initialBalance)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Printf("Failed to rollback transaction: %v", err)
			return
		}
		log.Printf("Reader failed to read initial balance: %v", err)
		return
	}
	fmt.Printf(" [Reader] Initial read (before writer's commit): Alice's balance is %.2f\n", initialBalance)

	close(writerReady)
	<-writerDone

	// Second read
	var secondBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE name = 'Alice';").Scan(&secondBalance)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Printf("Failed to rollback transaction: %v", err)
			return
		}
		log.Printf("Reader failed to read second balance: %v", err)
		return
	}
	fmt.Printf(" [Reader] Second read (after writer's commit): Alice's balance is %.2f\n", secondBalance)

	if initialBalance != secondBalance {
		fmt.Printf(" >>> ANOMALY DETECTED: Non-repeatable read! Balance changed from %.2f to %.2f\n", initialBalance, secondBalance)
	} else {
		fmt.Printf(" >>> ISOLATION SUCCESSFUL: Balance remained consistent (%.2f) throughout the transaction.\n", initialBalance)
	}

	// Only commit if successful, otherwise rollback is already handled
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit reader transaction: %v", err)
	}
}

func writerTransaction(db *sql.DB, writerReady chan struct{}, writerDone chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin writer transaction: %v", err)
		return
	}

	<-writerReady

	_, err = tx.Exec("UPDATE accounts SET balance = balance + 500.00 WHERE name = 'Alice';")
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Printf("Failed to rollback writer transaction: %v", err)
			return
		}
		log.Printf("Writer failed to update balance: %v", err)
		return
	}
	fmt.Println("\n [Writer] Alice's balance updated to 1500.00 (but not yet committed).")

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit writer transaction: %v", err)
		return
	}
	fmt.Println(" [Writer] Committed the balance change.")

	close(writerDone)
}

func main() {
	fmt.Println("Starting Isolation Level Demonstration...")

	connStr := "user=postgres password=mysecretpassword host=localhost port=5432 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing DB connection: %v", err)
		}
	}(db)

	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	var wg sync.WaitGroup

	// Demo 1: READ COMMITTED
	setupDatabase(db)
	fmt.Println("\n" + "=================================================")
	fmt.Println("DEMO 1: READ COMMITTED Isolation Level")
	fmt.Println("=================================================")

	writerReady1 := make(chan struct{})
	writerDone1 := make(chan struct{})
	wg.Add(2)
	go readWithIsolationLevel(db, sql.LevelReadCommitted, writerReady1, writerDone1, &wg)
	go writerTransaction(db, writerReady1, writerDone1, &wg)
	wg.Wait()

	// Demo 2: REPEATABLE READ
	setupDatabase(db)
	fmt.Println("\n" + "=================================================")
	fmt.Println("DEMO 2: REPEATABLE READ Isolation Level")
	fmt.Println("=================================================")

	writerReady2 := make(chan struct{})
	writerDone2 := make(chan struct{})
	wg.Add(2)
	go readWithIsolationLevel(db, sql.LevelRepeatableRead, writerReady2, writerDone2, &wg)
	go writerTransaction(db, writerReady2, writerDone2, &wg)
	wg.Wait()

	fmt.Println("\nDemonstration completed successfully!")
}

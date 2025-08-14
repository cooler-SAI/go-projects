package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func initDatabase(db *sql.DB) {
	// Drops and recreates the accounts table
	_, err := db.Exec(`
        DROP TABLE IF EXISTS accounts;
        CREATE TABLE accounts (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) UNIQUE NOT NULL,
            balance DECIMAL(10, 2) NOT NULL
        );
    `)
	if err != nil {
		log.Fatal(err)
	}
}

func createAccount(db *sql.DB, name string, balance float64) {
	// Inserts a new account with the given name and balance
	_, err := db.Exec("INSERT INTO accounts (name, balance) VALUES ($1, $2)", name, balance)
	if err != nil {
		log.Fatal(err)
	}
}

func printBalances(db *sql.DB) {
	// Retrieves and prints all account balances
	rows, err := db.Query("SELECT name, balance FROM accounts")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v\n", err)
		}
	}(rows)

	fmt.Println("Reading balances...")
	for rows.Next() {
		var name string
		var balance float64
		if err := rows.Scan(&name, &balance); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s: %.2f\n", name, balance)
	}
}

func transferMoney(db *sql.DB, from, to string, amount float64) error {
	// Starts a database transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Rolls back if there's an error
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			return
		}
	}()

	// Checks sender's balance (with row locking)
	var fromBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE name = $1 FOR UPDATE", from).Scan(&fromBalance)
	if err != nil {
		return fmt.Errorf("failed to check sender balance: %v", err)
	}

	if fromBalance < amount {
		return fmt.Errorf("insufficient funds: %s has only %.2f", from, fromBalance)
	}

	// Deducts from sender
	_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE name = $2", amount, from)
	if err != nil {
		return fmt.Errorf("failed to deduct funds: %v", err)
	}

	// Adds to recipient
	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE name = $2", amount, to)
	if err != nil {
		return fmt.Errorf("failed to add funds: %v", err)
	}

	// Commits the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func main() {
	// Database connection string
	connStr := "user=postgres password=example host=localhost port=5432 dbname=postgres sslmode=disable"

	// Opens a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing DB connection: %v\n", err)
		}
	}(db)

	// Initializes the database
	initDatabase(db)
	fmt.Println("Database initialized.")

	// Creates accounts
	createAccount(db, "Alice", 1000.00)
	createAccount(db, "Bob", 500.00)
	fmt.Println("Accounts created.")

	// Prints initial balances
	printBalances(db)

	// Transfers money
	fmt.Println("Transferring money...")
	err = transferMoney(db, "Alice", "Bob", 200.00)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Money transferred.")

	// Prints final balances
	printBalances(db)

	fmt.Println("Done.")
	fmt.Println("Don't forget to stop Docker-container with: docker stop my-postgres")

}

package main

import (
	"database/sql"
	"fmt"

	"github.com/cooler-SAI/go-Tools/zerolog"
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
		zerolog.Log.Fatal().Err(err).Msg("Failed to initialize database")
	}
}

func createAccount(db *sql.DB, name string, balance float64) {
	// Inserts a new account with the given name and balance
	_, err := db.Exec("INSERT INTO accounts (name, balance) VALUES ($1, $2)", name, balance)
	if err != nil {
		zerolog.Log.Fatal().Err(err).Msg("Failed to create account")
	}
}

func printBalances(db *sql.DB) {
	// Retrieves and prints all account balances
	rows, err := db.Query("SELECT name, balance FROM accounts")
	if err != nil {
		zerolog.Log.Fatal().Err(err).Msg("Failed to query balances")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zerolog.Log.Error().Err(err).Msg("Error closing rows")
		}
	}(rows)

	fmt.Println("Reading balances...")
	for rows.Next() {
		var name string
		var balance float64
		if err := rows.Scan(&name, &balance); err != nil {
			zerolog.Log.Fatal().Err(err).Msg("Failed to scan row")
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
	// Initializes the logger
	zerolog.Init()

	zerolog.Log.Info().Msg("Set up zerolog in production settings")

	fmt.Println("Starting transaction demonstration in Go with PostgreSQL...")

	// Database connection string
	connStr := "user=postgres password=example host=localhost port=5432 dbname=postgres sslmode=disable"

	// Opens a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		zerolog.Log.Fatal().Err(err).Msg("Failed to open database connection")
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			zerolog.Log.Error().Err(err).Msg("Error closing DB connection")
		}
	}(db)

	// Initializes the database
	initDatabase(db)
	zerolog.Log.Info().Msg("Database initialized.")

	// Creates accounts
	createAccount(db, "Alice", 1000.00)
	createAccount(db, "Bob", 500.00)
	zerolog.Log.Info().Msg("Accounts created.")

	// Prints initial balances
	printBalances(db)

	// Transfers money
	zerolog.Log.Info().Msg("Transferring money...")
	err = transferMoney(db, "Alice", "Bob", 200.00)
	if err != nil {
		zerolog.Log.Fatal().Err(err).Msg("Failed to transfer money")
	}
	zerolog.Log.Info().Msg("Money transferred.")

	// Prints final balances
	printBalances(db)

	zerolog.Log.Info().Msg("Done.")
	zerolog.Log.Info().Msg("Don't forget to stop Docker-container with: docker stop my-postgres")
}

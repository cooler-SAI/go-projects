package main

import (
	"database/sql" // Standard Go package for database operations
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver import
	"log"
)

// Account - struct to represent a bank account.
type Account struct {
	ID      int
	Name    string
	Balance float64
}

// transferFunds encapsulates the transfer logic within a transaction.
// tx: The transaction object. All operations inside this function will be part of this transaction.
// fromAccountID: ID of the account to debit.
// toAccountID: ID of the account to credit.
// amount: Amount to transfer.
func transferFunds(tx *sql.Tx, fromAccountID, toAccountID int, amount float64) error {
	// Debit funds
	_, err := tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2;", amount, fromAccountID)
	if err != nil {
		return fmt.Errorf("error debiting funds from account %d: %w", fromAccountID, err)
	}
	fmt.Printf("  Debited %.2f from account %d\n", amount, fromAccountID)

	// Simulate an error (uncomment to see transaction rollback)
	// if fromAccountID == 1 && amount == 50.00 { // For example, if transfer from Alice with 50, simulate an error
	// 	return fmt.Errorf("simulated error after debit from account %d", fromAccountID)
	// }

	// Credit funds
	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2;", amount, toAccountID)
	if err != nil {
		return fmt.Errorf("error crediting funds to account %d: %w", toAccountID, err)
	}
	fmt.Printf("  Credited %.2f to account %d\n", amount, toAccountID)
	return nil
}

func main() {
	fmt.Println("Starting transaction demonstration in Go with PostgreSQL...")

	connStr := "user=postgres password=mysecretpassword host=localhost port=5432 dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing DB connection")
		}
	}(db)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	// --- 1. Create 'accounts' table ---
	fmt.Println("\nCreating 'accounts' table...")
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		balance NUMERIC(10, 2) NOT NULL DEFAULT 0.00
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating 'accounts' table: %v", err)
	}
	fmt.Println("Table 'accounts' created successfully or already exists.")

	// --- 2. Clear and insert initial data ---
	fmt.Println("\nClearing and inserting initial data into 'accounts'...")
	_, err = db.Exec("DELETE FROM accounts;") // Clear table for a clean experiment
	if err != nil {
		log.Fatalf("Error clearing table: %v", err)
	}

	insertAccountSQL := `INSERT INTO accounts(name, balance) VALUES ($1, $2) RETURNING id;`
	var accountID int

	err = db.QueryRow(insertAccountSQL, "Alice", 1000.00).Scan(&accountID)
	if err != nil {
		log.Fatalf("Error inserting account for Alice: %v", err)
	}
	fmt.Printf("Created account for Alice with ID: %d, Balance: 1000.00\n", accountID)

	err = db.QueryRow(insertAccountSQL, "Bob", 500.00).Scan(&accountID)
	if err != nil {
		log.Fatalf("Error inserting account for Bob: %v", err)
	}
	fmt.Printf("Created account for Bob with ID: %d, Balance: 500.00\n", accountID)

	// --- 4. Execute transaction (successful scenario) ---
	fmt.Println("\n--- Scenario 1: Successful transfer (Alice -> Bob, 100.00) ---")
	tx, err := db.Begin() // Start a new transaction
	if err != nil {
		log.Fatalf("Error starting transaction: %v", err)
	}

	err = transferFunds(tx, 1, 2, 100.00) // Transfer 100 from Alice (ID 1) to Bob (ID 2)
	if err != nil {
		fmt.Printf("Transfer error: %v. Rolling back transaction.\n", err)
		err := tx.Rollback()
		if err != nil {
			log.Fatalf("Error rolling back transaction: %v", err)
			return
		} // Rollback transaction on error
	} else {
		err := tx.Commit()
		if err != nil {
			log.Fatalf("Error committing transaction: %v", err)
			return
		} // Commit transaction if all is successful
		fmt.Println("Transfer completed and committed successfully.")
	}

	// --- 5. Execute transaction (error and rollback scenario) ---
	// Uncomment the "Simulate an error" block in transferFunds function to see rollback.
	fmt.Println("\n--- Scenario 2: Transfer with error (Alice -> Bob, 50.00) ---")
	tx2, err := db.Begin() // Start another new transaction
	if err != nil {
		log.Fatalf("Error starting second transaction: %v", err)
	}

	err = transferFunds(tx2, 1, 2, 50.00) // Transfer 50 from Alice (ID 1) to Bob (ID 2)
	if err != nil {
		fmt.Printf("Transfer error: %v. Rolling back transaction.\n", err)
		err := tx2.Rollback()
		if err != nil {
			fmt.Println("Error rolling back transaction")
			return
		} // Rollback transaction on error
	} else {
		err := tx2.Commit()
		if err != nil {
			return
		} // Commit transaction if all is successful
		fmt.Println("Transfer completed and committed successfully.")
	}

	// --- 6. Check final balances ---
	fmt.Println("\nFinal account balances:")
	rows, err := db.Query("SELECT id, name, balance FROM accounts ORDER BY id;")
	if err != nil {
		log.Fatalf("Error querying final balances: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error closing rows: ", err)
		}
	}(rows)

	for rows.Next() {
		var acc Account
		err := rows.Scan(&acc.ID, &acc.Name, &acc.Balance)
		if err != nil {
			log.Fatalf("Error scanning account: %v", err)
		}
		fmt.Printf("  ID: %d, Name: %s, Balance: %.2f\n", acc.ID, acc.Name, acc.Balance)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("Error after iterating through accounts: %v", err)
	}

	fmt.Println("\nTransaction demonstration completed.")
	fmt.Println("Don't forget to stop Docker-container with: docker stop my-postgres")
}

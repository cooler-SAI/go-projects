package main

import (
	"database/sql" // Provides a generic SQL interface.
	"fmt"          // For formatted I/O.
	"log"          // For logging errors.

	_ "github.com/lib/pq" // The PostgreSQL driver. The blank identifier `_` is used because we only need its side effects (registering the driver).
)

// Account - struct to represent a bank account.
type Account struct {
	ID      int
	Name    string
	Balance float64
}

// transferFunds encapsulates the core logic of transferring money between two accounts.
// It's designed to be executed within a database transaction to ensure atomicity.
// If any operation within this function fails, the entire transaction can be rolled back,
// leaving the database in its original state.
//
// Parameters:
//
//	tx: The active SQL transaction object.
//	fromAccountID: The ID of the account from which funds will be debited.
//	toAccountID: The ID of the account to which funds will be credited.
//	amount: The amount of money to transfer.
func transferFunds(tx *sql.Tx, fromAccountID, toAccountID int, amount float64) error {
	// Step 1: Debit the amount from the sender's account.
	_, err := tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2;", amount, fromAccountID)
	if err != nil {
		// If the debit fails, wrap the error with more context and return it.
		return fmt.Errorf("error debiting funds from account %d: %w", fromAccountID, err)
	}
	fmt.Printf("  Debited %.2f from account %d\n", amount, fromAccountID)

	// --- SIMULATED ERROR ---
	// Uncomment the following lines to test the transaction rollback mechanism.
	// This simulates a failure occurring after the debit but before the credit.
	// if amount == 50.00 {
	// 	return fmt.Errorf("simulated network error after debit")
	// }

	// Step 2: Credit the amount to the receiver's account.
	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2;", amount, toAccountID)
	if err != nil {
		// If the credit fails, wrap the error and return it.
		return fmt.Errorf("error crediting funds to account %d: %w", toAccountID, err)
	}
	fmt.Printf("  Credited %.2f to account %d\n", amount, toAccountID)
	return nil
}

func main() {
	fmt.Println("Starting transaction demonstration in Go with PostgreSQL...")

	// Database connection string. Replace with your actual database credentials.
	connStr := "user=postgres password=example host=localhost port=5432 dbname=postgres sslmode=disable"

	// sql.Open() initializes a connection pool. It does not create a connection itself.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}
	// defer db.Close() ensures that the connection pool is closed before the main function exits.
	// It's wrapped in a function to handle any potential error from Close().
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing DB connection: %v\n", err)
		}
	}(db)

	// db.Ping() verifies that a connection to the database is still alive,
	// establishing a connection if necessary.
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	// --- 1. Create 'accounts' table ---
	fmt.Println("\nCreating 'accounts' table...")
	// "CREATE TABLE IF NOT EXISTS" is idempotent; it won't cause an error if the table already exists.
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
	// Clear the table to ensure a fresh start for each run of the script.
	_, err = db.Exec("DELETE FROM accounts;")
	if err != nil {
		log.Fatalf("Error clearing table: %v", err)
	}

	// Insert new accounts and capture their generated IDs.
	// Using RETURNING id is efficient as it avoids a separate SELECT query.
	insertAccountSQL := `INSERT INTO accounts(name, balance) VALUES ($1, $2) RETURNING id;`
	var aliceID, bobID int

	// Insert Alice's account and scan the returned ID into the aliceID variable.
	err = db.QueryRow(insertAccountSQL, "Alice", 1000.00).Scan(&aliceID)
	if err != nil {
		log.Fatalf("Error inserting account for Alice: %v", err)
	}
	fmt.Printf("Created account for Alice with ID: %d, Balance: 1000.00\n", aliceID)

	// Insert Bob's account and scan the returned ID into the bobID variable.
	err = db.QueryRow(insertAccountSQL, "Bob", 500.00).Scan(&bobID)
	if err != nil {
		log.Fatalf("Error inserting account for Bob: %v", err)
	}
	fmt.Printf("Created account for Bob with ID: %d, Balance: 500.00\n", bobID)

	// --- 3. Execute transaction (successful scenario) ---
	fmt.Println("\n--- Scenario 1: Successful transfer (Alice -> Bob, 100.00) ---")
	// db.Begin() starts a new database transaction.
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error starting transaction: %v", err)
	}

	// Perform the transfer using the dynamically retrieved account IDs.
	err = transferFunds(tx, aliceID, bobID, 100.00)
	if err != nil {
		// If transferFunds returns an error, something went wrong.
		fmt.Printf("Transfer error: %v. Rolling back transaction.\n", err)
		// We must roll back the transaction to undo any partial changes (like the initial debit).
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Fatalf("Error rolling back transaction: %v", err)
		}
	} else {
		// If transferFunds completes without error, we commit the transaction.
		// This makes all changes within the transaction permanent.
		if cmtErr := tx.Commit(); cmtErr != nil {
			log.Fatalf("Error committing transaction: %v", err)
		}
		fmt.Println("Transfer completed and committed successfully.")
	}

	// --- 4. Execute transaction (error and rollback scenario) ---
	// To test this scenario, uncomment the "SIMULATED ERROR" block in the transferFunds function.
	fmt.Println("\n--- Scenario 2: Transfer with error (Alice -> Bob, 50.00) ---")
	tx2, err := db.Begin() // Start a second transaction.
	if err != nil {
		log.Fatalf("Error starting second transaction: %v", err)
	}

	err = transferFunds(tx2, aliceID, bobID, 50.00)
	if err != nil {
		fmt.Printf("Transfer error: %v. Rolling back transaction.\n", err)
		// The rollback is crucial here. If the simulated error occurs after the debit,
		// the rollback will reverse that debit, ensuring Alice's balance is not incorrectly reduced.
		if rbErr := tx2.Rollback(); rbErr != nil {
			log.Fatalf("Error rolling back second transaction: %v", rbErr)
		}
	} else {
		if cmtErr := tx2.Commit(); cmtErr != nil {
			log.Fatalf("Error committing second transaction: %v", cmtErr)
		}
		fmt.Println("Transfer completed and committed successfully.")
	}

	// --- 6. Check final balances ---
	fmt.Println("\nFinal account balances:")
	rows, err := db.Query("SELECT id, name, balance FROM accounts ORDER BY id;")
	if err != nil {
		log.Fatalf("Error querying final balances: %v", err)
	}
	// It's crucial to close the rows iterator to release the database connection.
	// Deferring it ensures it runs even if there's a panic during row scanning.
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v\n", err)
		}
	}(rows)

	// Iterate over the result set.
	for rows.Next() {
		var acc Account
		err := rows.Scan(&acc.ID, &acc.Name, &acc.Balance)
		if err != nil {
			log.Fatalf("Error scanning account: %v", err)
		}
		fmt.Printf("  ID: %d, Name: %s, Balance: %.2f\n", acc.ID, acc.Name, acc.Balance)
	}
	// After the loop, check for any errors that occurred during iteration.
	if err = rows.Err(); err != nil {
		log.Fatalf("Error after iterating through accounts: %v", err)
	}

	fmt.Println("\nTransaction demonstration completed.")
	fmt.Println("Don't forget to stop Docker-container with: docker stop my-postgres")
}

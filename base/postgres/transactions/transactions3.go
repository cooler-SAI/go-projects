package main

import (
	"database/sql"
	"fmt"
	"log"
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

func readWithIsolationLevel() {

}

func main() {

	fmt.Println("Starting Isolation Level Demonstration...")

}

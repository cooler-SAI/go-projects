package main

import (
	"database/sql"
	"fmt"
	"log"
)

// Account - struct to represent a bank account.
type Account struct {
	ID      int
	Name    string
	Balance float64
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
			log.Fatalf("Error closing DB connection: %v", err)
		}
	}(db)
}

package main

import (
	"database/sql"
	"fmt"
	"log"

	// Import PostgreSQL driver
	_ "github.com/lib/pq"
)

// DSN (Data Source Name) for connecting to PostgreSQL Docker container
// Uses: user=postgres, password=mysecretpassword, host=localhost, port=5432
const DSN = "user=postgres password=mysecretpassword host=localhost port=5432 sslmode=disable"

// Item represents the data structure we will store in the database
type Item struct {
	ID    int
	Name  string
	Price float64
}

func main() {
	// 1. Connect to the database
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		// panic is not typically used here, but for connection error demonstration
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing the database connection: %v", err)
		}
	}(db) // Ensure connection is closed

	// Verify the connection is established
	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}
	fmt.Println("Successfully connected to PostgreSQL (Docker).")

	// 2. Create table
	err = createTable(db)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table 'items' created successfully.")

	// 3. Insert data
	itemsToInsert := []Item{
		{Name: "Go Gopher Plush", Price: 19.99},
		{Name: "Postgres Sticker", Price: 2.50},
		{Name: "Go Book", Price: 49.99},
	}
	for _, item := range itemsToInsert {
		if err := insertItem(db, item); err != nil {
			log.Printf("Failed to insert item %s: %v", item.Name, err)
		}
	}
	fmt.Println("Data inserted successfully.")

	// 4. Read data
	if err := readItems(db); err != nil {
		log.Fatalf("Error reading items: %v", err)
	}
}

// createTable drops the old table and creates a new 'items' table
func createTable(db *sql.DB) error {
	// First drop if exists
	_, err := db.Exec("DROP TABLE IF EXISTS items")
	if err != nil {
		return err
	}

	// Create new table
	createTableSQL := `CREATE TABLE items (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(10, 2) NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	return err
}

// insertItem inserts a single record into the 'items' table
func insertItem(db *sql.DB, item Item) error {
	// 'RETURNING id' allows us to get the database-generated ID
	query := "INSERT INTO items (name, price) VALUES ($1, $2) RETURNING id"
	var id int
	err := db.QueryRow(query, item.Name, item.Price).Scan(&id)

	if err != nil {
		return err
	}
	// Update ID in the structure for future use (if needed)
	item.ID = id
	return nil
}

// readItems reads and prints all records from the 'items' table
func readItems(db *sql.DB) error {
	fmt.Println("\n--- Items in the Database ---")

	rows, err := db.Query("SELECT id, name, price FROM items ORDER BY id")
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("Error closing rows: %v\n", err)
		}
	}(rows)

	for rows.Next() {
		var item Item
		// Scan results from the row into Item structure fields
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}
		fmt.Printf("ID: %d, Name: %s, Price: %.2f\n", item.ID, item.Name, item.Price)
	}

	// Check for errors after the loop completes
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error during rows iteration: %w", err)
	}

	return nil
}

// NOTE: To run this code you will need to install the external PostgreSQL driver:
// go get github.com/lib/pq

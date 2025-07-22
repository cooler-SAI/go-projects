package main

import (
	"database/sql" // Standard Go package for database operations
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver import
	"log"
)

// Phone - struct to represent phone data in the 'phones' table.
type Phone struct {
	ID    int
	Brand string
	Model string
	Year  int
}

func main() {
	fmt.Println("Start using PostgreSQL with Go")

	// Database connection string.
	// Ensure your 'my-postgres' Docker container is running!
	connStr := "user=postgres password=mysecretpassword host=localhost port=5432 " +
		"dbname=postgres sslmode=disable"

	// Open database connection.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Can't open connection with DB: %v. Please add PostgreSQL "+
			"support in import here: _ github.com/lib/pq", err)
		return
	}
	// Close the database connection when main function exits.
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Close DB: %v", err)
		}
	}(db)

	// Ping the database to verify connection is established.
	err = db.Ping()
	if err != nil {
		log.Fatalf("Can't open DB: %v", err)
		return
	}
	fmt.Println("Successfully connected to PostgreSQL database")

	// --- 1. Create Table 'phones' ---
	fmt.Println("\nCreating table 'phones'....")
	createTableSQL := `
CREATE TABLE IF NOT EXISTS phones(
    id SERIAL PRIMARY KEY,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INT NOT NULL
);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Can't create table: %v", err)
		return
	}
	fmt.Println("Table 'phones' created successfully")

	// --- 2. Insert Data ---
	fmt.Println("\nInserting data into 'phones' table...")
	insertSQL := `INSERT INTO phones(brand, model, year) VALUES ($1, $2, $3) RETURNING id, brand, model, year;`

	var (
		insertedID    int
		insertedBrand string
		insertedModel string
		insertedYear  int
	)

	// Insert Nokia
	err = db.QueryRow(insertSQL, "Nokia", "3310", 2000).Scan(
		&insertedID, &insertedBrand, &insertedModel, &insertedYear,
	)
	if err != nil {
		log.Fatalf("Error inserting phone: %v", err)
	}
	fmt.Printf("Inserted: %s with ID: %d, Model: %s, Year: %d\n",
		insertedBrand, insertedID, insertedModel, insertedYear,
	)

	// Insert another phone for better index demonstration
	err = db.QueryRow(insertSQL, "Samsung", "Galaxy S21", 2021).Scan(
		&insertedID, &insertedBrand, &insertedModel, &insertedYear,
	)
	if err != nil {
		log.Fatalf("Error inserting phone: %v", err)
	}
	fmt.Printf("Inserted: %s with ID: %d, Model: %s, Year: %d\n",
		insertedBrand, insertedID, insertedModel, insertedYear,
	)

	err = db.QueryRow(insertSQL, "Apple", "iPhone 13", 2021).Scan(
		&insertedID, &insertedBrand, &insertedModel, &insertedYear,
	)
	if err != nil {
		log.Fatalf("Error inserting phone: %v", err)
	}
	fmt.Printf("Inserted: %s with ID: %d, Model: %s, Year: %d\n",
		insertedBrand, insertedID, insertedModel, insertedYear,
	)

	// --- 3. Create Index on 'year' column ---
	fmt.Println("\nCreating index on 'year' column...")
	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_phones_year ON phones (year);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		log.Fatalf("Error creating index: %v", err)
		return
	}
	fmt.Println("Index 'idx_phones_year' created successfully or already exists.")

	// --- 4. Reading all data (original query) ---
	fmt.Println("\nReading all data from 'phones' table (full list)...")
	rows, err := db.Query("SELECT id, brand, model, year FROM phones ORDER BY id")
	if err != nil {
		log.Fatalf("Error querying data: %v", err)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalf("Error closing rows: %v", err)
		}
	}(rows)

	fmt.Println("\nAll phones in database:")
	for rows.Next() {
		var p Phone
		err := rows.Scan(&p.ID, &p.Brand, &p.Model, &p.Year)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		fmt.Printf("ID: %d, Brand: %s, Model: %s, Year: %d\n",
			p.ID, p.Brand, p.Model, p.Year)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error after row iteration: %v", err)
	}

	// --- 5. Reading data using the 'year' index (example query) ---
	fmt.Println("\nReading phones manufactured in 2021 (query potentially using index)...")
	var phones2021 []Phone
	query2021SQL := `SELECT id, brand, model, year FROM phones WHERE year = $1 ORDER BY brand;`
	rows2021, err := db.Query(query2021SQL, 2021)
	if err != nil {
		log.Fatalf("Error querying phones for 2021: %v", err)
		return
	}
	defer func(rows2021 *sql.Rows) {
		err := rows2021.Close()
		if err != nil {
			log.Fatalf("Error closing rows for 2021: %v", err)
		}
	}(rows2021)

	fmt.Println("Phones from 2021:")
	for rows2021.Next() {
		var p Phone
		err := rows2021.Scan(&p.ID, &p.Brand, &p.Model, &p.Year)
		if err != nil {
			log.Fatalf("Error scanning phone row for 2021: %v", err)
		}
		phones2021 = append(phones2021, p)
		fmt.Printf("  ID: %d, Brand: %s, Model: %s, Year: %d\n", p.ID, p.Brand, p.Model, p.Year)
	}
	if err = rows2021.Err(); err != nil {
		log.Fatalf("Error after row iteration for 2021: %v", err)
	}
	if len(phones2021) == 0 {
		fmt.Println("  No phones found for 2021.")
	}

	fmt.Println("\nDemonstration of PostgreSQL with Go completed.")
	fmt.Println("Don't forget to stop Docker-container with: docker stop my-postgres")
}

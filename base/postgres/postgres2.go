package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Car struct {
	ID    int
	Brand string
	Model string
	Year  int
}

func main() {
	fmt.Println("Start using PostgresSQL with Go")

	connStr := "user=postgres password=mysecretpassword host=localhost port=5432 " +
		"dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Cant open connection with DB: %v. \nPlease add posgreSQL "+
			"support in import here: _ github.com/lib/pq", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Close DB: %v", err)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Cant open DB: %v", err)
		return
	}
	fmt.Println("Successfully connected to PostgreSQL database")

	/*// -- Delete DB
	deleteTableSQL := `DROP TABLE IF EXISTS cars;`
	_, err = db.Exec(deleteTableSQL)
	if err != nil {
		log.Fatalf("Cant delete table: %v", err)
		return
	}
	fmt.Println("Table 'cars' deleted successfully")*/

	// -- Create Table
	createTableSQL := `
CREATE TABLE IF NOT EXISTS cars(
    id SERIAL PRIMARY KEY,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INT NOT NULL
)`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Cant create table: %v", err)
		return
	}
	fmt.Println("Table 'cars' created successfully")

	// -- Insert Data
	insertSQL := `INSERT INTO cars(brand, model, year) VALUES ($1, $2, $3)
RETURNING id, brand, model, year;`

	var (
		insertedID    int
		insertedBrand string
		insertedModel string
		insertedYear  int
	)

	// Insert Toyota
	err = db.QueryRow(insertSQL, "Toyota", "Camry", 2022).Scan(
		&insertedID,
		&insertedBrand,
		&insertedModel,
		&insertedYear,
	)
	if err != nil {
		log.Fatalf("Error inserting car: %v", err)
	}

	fmt.Printf("Inserted: %s with ID: %d, Model: %s, Year: %d\n",
		insertedBrand,
		insertedID,
		insertedModel,
		insertedYear,
	)

	// Insert Honda
	err = db.QueryRow(insertSQL, "Honda", "Accord", 2021).Scan(
		&insertedID,
		&insertedBrand,
		&insertedModel,
		&insertedYear,
	)
	if err != nil {
		log.Fatalf("Error inserting car: %v", err)
	}

	fmt.Printf("Inserted: %s with ID: %d, Model: %s, Year: %d\n",
		insertedBrand,
		insertedID,
		insertedModel,
		insertedYear,
	)

	// Insert BMW
	err = db.QueryRow(insertSQL, "BMW", "X5", 2023).Scan(
		&insertedID,
		&insertedBrand,
		&insertedModel,
		&insertedYear,
	)
	if err != nil {
		log.Fatalf("Error inserting car: %v", err)
	}

	fmt.Printf("Inserted: %s with ID: %d, Model: %s, Year: %d\n",
		insertedBrand,
		insertedID,
		insertedModel,
		insertedYear,
	)

	// Insert Tesla
	err = db.QueryRow(insertSQL, "Tesla", "Model 3", 2022).Scan(
		&insertedID,
		&insertedBrand,
		&insertedModel,
		&insertedYear,
	)
	if err != nil {
		log.Fatalf("Error inserting car: %v", err)
	}

	fmt.Printf("Inserted: %s with ID: %d, Model: %s, Year: %d\n",
		insertedBrand,
		insertedID,
		insertedModel,
		insertedYear,
	)

	// --- 3. Reading all data ---
	rows, err := db.Query("SELECT id, brand, model, year FROM cars ORDER BY id")
	if err != nil {
		log.Fatalf("Error querying data: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalf("Error closing rows: %v", err)
		}
	}(rows)

	fmt.Println("\nAll cars in database:")
	for rows.Next() {
		var c Car
		err := rows.Scan(&c.ID, &c.Brand, &c.Model, &c.Year)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		fmt.Printf("ID: %d, Brand: %s, Model: %s, Year: %d\n",
			c.ID, c.Brand, c.Model, c.Year)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error after row iteration: %v", err)
	}

	fmt.Println("\nDemonstration of PostgreSQL with Go completed.")
	fmt.Println("Don't forget to stop Docker-container with: docker stop my-postgres")
}

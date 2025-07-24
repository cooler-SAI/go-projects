package main

import (
	"database/sql" // Standard Go package for database operations
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver import
	"log"
)

type Laptop struct {
	ID    int
	Model string
	Year  int
}

func main() {
	fmt.Println("Start using PostgreSQL with Go")

	connStr := "user=postgres password=mysecretpassword host=localhost port=5432 " +
		"dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Can't open connection with DB: %v. Please add PostgreSQL "+
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
		log.Fatalf("Can't open DB: %v", err)
		return
	}
	fmt.Println("Successfully connected to PostgreSQL database")

	createTableSQL := `
CREATE TABLE IF NOT EXISTS laptops(
    id SERIAL PRIMARY KEY,
    model VARCHAR(100)
        NOT NULL,
    year INT NOT NULL
);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Can't create table: %v", err)
		return
	}
	fmt.Println("Table 'laptops' created successfully")

	fmt.Println("Create index in 'year' table:")
	createIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_laptops_year 
	ON laptops(year);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		log.Fatalf("Can't create index: %v", err)
	}
	fmt.Println("Index on 'year' created successfully")

	fmt.Println("\nInserting data into 'phones' table...")

	var insertedLaptop Laptop
	err = db.QueryRow(`
		INSERT INTO laptops(model, year) 
		VALUES ($1, $2) 
		RETURNING id, model, year`,
		"Dell XPS", 2022,
	).Scan(&insertedLaptop.ID, &insertedLaptop.Model, &insertedLaptop.Year)
	if err != nil {
		log.Fatalf("Error inserting laptop: %v", err)
	}

	fmt.Printf("Inserted laptop: %+v\n", insertedLaptop)

}

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Person struct {
	ID   int
	Name string
	Age  int
}

func main() {
	fmt.Println("Start demonstration work Go with PostgreSQL.....")

	connStr := "user=postgres password=mysecretpassword host=localhost port=5432 " +
		"dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Cant open connection with DB: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Close DB: %v", err)
		}
	}(db)

	fmt.Println("Trying connect to DB....")

	err = db.Ping()
	if err != nil {
		log.Fatalf("Cant open DB: %v", err)
	}
	fmt.Println("Successfully connected to PostgreSQL database")

	// --- 1. Create table ---
	fmt.Println("Creating table 'persons'....")
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS persons (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		age INT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Cant create table: %v", err)
	}
	fmt.Println("Table 'persons' created successfully")

	// --- 2. Inserting data ---
	fmt.Println("\nInserting data into 'persons' table...")
	insertSQL := `INSERT INTO persons(name, age) VALUES ($1, $2) RETURNING id;`

	var insertedID int

	err = db.QueryRow(insertSQL, "Alexey", 30).Scan(&insertedID)
	if err != nil {
		log.Fatalf("Error inserting Alexey: %v", err)
	}
	fmt.Printf("Inserted Alexey with ID: %d\n", insertedID)

	err = db.QueryRow(insertSQL, "Maria", 25).Scan(&insertedID)
	if err != nil {
		log.Fatalf("Error inserting Maria: %v", err)
	}
	fmt.Printf("Inserted Maria with ID: %d\n", insertedID)

	// --- 3. Reading data ---
	fmt.Println("\nReading all data from 'persons' table...")
	rows, err := db.Query("SELECT id, name, age FROM persons;")
	if err != nil {
		log.Fatalf("Error querying data: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalf("Error closing rows: %v", err)
		}
	}(rows)

	var persons []Person
	for rows.Next() {
		var p Person
		err := rows.Scan(&p.ID, &p.Name, &p.Age)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		persons = append(persons, p)
	}
	// Check for errors that may have occurred during iteration
	if err = rows.Err(); err != nil {
		log.Fatalf("Error after row iteration: %v", err)
	}

	fmt.Println("Data from persons table:")
	for _, p := range persons {
		fmt.Printf("  ID: %d, Name: %s, Age: %d\n", p.ID, p.Name, p.Age)
	}

	// --- 4. Updating data ---
	fmt.Println("\nUpdating data (Alexey's age)...")
	updateSQL := `UPDATE persons SET age = $1 WHERE name = $2;`
	res, err := db.Exec(updateSQL, 31, "Alexey")
	if err != nil {
		log.Fatalf("Error updating Alexey: %v", err)
	}
	rowsAffected, _ := res.RowsAffected()
	fmt.Printf("Rows updated: %d\n", rowsAffected)

	// --- 5. Deleting data (optional) ---
	// Commented by default to keep data for multiple runs
	// Uncomment if you want to test deletion
	/*
	   fmt.Println("\nDeleting data (Maria)...")
	   deleteSQL := `DELETE FROM persons WHERE name = $1;`
	   res, err = db.Exec(deleteSQL, "Maria")
	   if err != nil {
	       log.Fatalf("Error deleting Maria: %v", err)
	   }
	   rowsAffected, _ = res.RowsAffected()
	   fmt.Printf("Rows deleted: %d\n", rowsAffected)
	*/

	fmt.Println("\nPostgreSQL with Go demonstration completed.")
	fmt.Println("You can stop the Docker container with: docker stop my-postgres")
	fmt.Println("And remove it: docker rm my-postgres")
}

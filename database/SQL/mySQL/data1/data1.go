package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Connect directly to the database
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/simple_db")
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Error closing the database:", err)
		}
	}(db)

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}
	fmt.Println("✅ Connected to MySQL!")

	// Create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			age INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Table creation error:", err)
	}
	fmt.Println("✅ Table 'users' created")

	// Add test data
	_, err = db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", "John", 25)
	if err != nil {
		log.Fatal("Data insertion error:", err)
	}
	fmt.Println("✅ Test data added")

	// Read data from the table
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		log.Fatal("Data reading error:", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal("Error closing rows:", err)
		}
	}(rows)

	fmt.Println("\n📋 Table data:")
	for rows.Next() {
		var id int
		var name string
		var age int
		err = rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatal("Scan error:", err)
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}

	fmt.Println("\n🎉 All done! Table created and populated.")
}

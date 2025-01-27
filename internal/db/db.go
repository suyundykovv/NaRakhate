package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// InitDB initializes and returns the database connection
func InitDB() (*sql.DB, error) {
	// Get the environment variables
	dbHost := os.Getenv("POSTGRES_HOST")         // Should be 'postgres' based on docker-compose.yml
	dbPort := "5432"                             // Default Postgres port
	dbUser := os.Getenv("POSTGRES_USER")         // 'admin' as per docker-compose.yml
	dbPassword := os.Getenv("POSTGRES_PASSWORD") // 'adminpassword' as per docker-compose.yml
	dbName := os.Getenv("POSTGRES_DB")           // 'aitubet' as per docker-compose.yml

	// Construct the connection string
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	log.Println("Connected to the database!")
	return db, nil
}

package main

import (
	"Aitu-Bet/config"
	"Aitu-Bet/flags"
	"Aitu-Bet/internal/db"
	"Aitu-Bet/internal/servers"
	"Aitu-Bet/logging"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	flags.Setup()

	log.Println("Attempting to connect to the database...")

	if err := logging.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	// Retry loop to connect to the database
	var dbConn *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		dbConn, err = db.InitDB()
		if err == nil {
			log.Println("Successfully connected to the database!")
			break
		}
		log.Printf("Error connecting to the database (attempt %d): %v\n", i+1, err)
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to the database after 5 attempts: %v\n", err)
		os.Exit(1)
	}

	// Pass the DB connection to the server
	server := servers.NewServer(dbConn)
	server.Start(config.Port)
}

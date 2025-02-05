package main

import (
	"Aitu-Bet/config"
	"Aitu-Bet/flags"
	"Aitu-Bet/internal/api"
	"Aitu-Bet/internal/db"
	"Aitu-Bet/internal/services"
	"Aitu-Bet/logging"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	flags.Setup()

	log.Println("Attempting to connect to the database...")

	if err := logging.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	// Retry loop for database connection
	var dbConn *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		dbConn, err = db.InitDB()
		if err == nil {
			log.Println("Successfully connected to the database!")
			break
		}
		log.Printf("Error connecting to the database (attempt %d): %v\n", i+1, err)
		time.Sleep(5 * time.Second) // Wait before retrying
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to the database after 5 attempts: %v\n", err)
		os.Exit(1)
	}

	// Создаем сервис Leaderboard
	leaderboardService := &services.LeaderboardService{
		LeaderboardDAL: &db.LeaderboardDAL{DB: dbConn},
	}

	// Создаем роутер Gin
	router := gin.Default()

	// Регистрируем маршруты
	api.RegisterRoutes(router, leaderboardService)

	// Запускаем сервер
	port := fmt.Sprintf(":%s", config.Port)
	log.Printf("Server is running on port %s", port)
	log.Fatal(router.Run(port))
}

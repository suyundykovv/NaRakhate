package db

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL драйвер
	"log"
)

// Глобальная переменная для хранения подключения к базе данных
var DB *sql.DB

// Инициализация базы данных
func InitDB(connStr string) error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return err
	}

	// Проверяем соединение с базой данных
	err = DB.Ping()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return err
	}

	log.Println("Successfully connected to the database!")
	return nil
}

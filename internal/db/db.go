package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // Импорт драйвера для PostgreSQL
)

var DB *sql.DB // Глобальное подключение к БД

// InitDB инициализирует подключение к БД
func InitDB() (*sql.DB, error) {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := "5432"
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	// Формируем строку подключения
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort,
	)

	// Открываем подключение
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("❌ Error opening database: %v", err)
		return nil, err
	}

	// Проверяем подключение к БД с ретраями (до 5 раз)
	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("✅ Successfully connected to the database!")
			break
		}
		log.Printf("⚠️ Error pinging database (attempt %d): %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Println("❌ Failed to connect to database after 5 attempts.")
		return nil, err
	}

	// Настройка соединений с БД
	db.SetMaxOpenConns(25)                 // Максимальное количество открытых соединений
	db.SetMaxIdleConns(10)                 // Максимальное число неактивных соединений
	db.SetConnMaxLifetime(5 * time.Minute) // Время жизни соединения

	DB = db // Сохраняем подключение в глобальную переменную
	return db, nil
}

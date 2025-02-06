package main

import (
	"Aitu-Bet/config"
	"Aitu-Bet/internal/db"
	"Aitu-Bet/internal/servers"
	"log"
)

func main() {
	// Получаем строку подключения из конфигурации
	connStr := config.GetDBConnectionString()

	// Инициализируем подключение к базе данных
	err := db.InitDB(connStr) // Исправлено, теперь InitDB возвращает только ошибку
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Инициализируем и запускаем сервер
	server := servers.NewServer(db.DB)
	server.Start("8080")
}

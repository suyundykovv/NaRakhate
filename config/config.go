package config

import (
	"fmt"
	"os"
)

var (
	// Порт для сервера
	Port string
	// Путь к директории для хранения данных
	StorageDir string
	// Путь к файлу меню
	MenuFile string
	// Путь к файлу для логирования
	LogFile string
	// Список ограниченных директорий
	RestrictedDirs = []string{"flags", "handlers", "models", "servers", "storage", "utils", "../"}

	// Данные для подключения к базе данных
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
)

// Инициализация конфигурации
func init() {
	// Загружаем порт для сервера из переменных окружения
	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080" // если переменная не установлена, используем порт 8080
	}

	// Получаем данные для подключения к базе данных
	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort = os.Getenv("POSTGRES_PORT")
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	PostgresDB = os.Getenv("POSTGRES_DB")
}

// Формируем строку подключения к базе данных
func GetDBConnectionString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		PostgresUser, PostgresPassword, PostgresDB, PostgresHost, PostgresPort)
}

// Функция для вывода справки по использованию приложения
func PrintUsage() {
	fmt.Println("Data Management System")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("    data [--port <N>] [--dir <S>] ")
	fmt.Println("    data --help")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help     Show this screen.")
	fmt.Println("  --port N   Port number")
	fmt.Println("  --dir S    Path to the directory")
}

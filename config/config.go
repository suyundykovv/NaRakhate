package config

import (
	"fmt"
	"os"
)

var (
	Port           string
	StorageDir     string
	MenuFile       string
	LogFile        string
	RestrictedDirs = []string{"flags", "handlers", "models", "servers", "storage", "utils", "../"}
)

func init() {
	// Загружаем порт из переменной окружения, если не установлен — используем 8080
	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}
}

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

package config

import "fmt"

var (
	Port           string
	StorageDir     string
	MenuFile       string
	LogFile        string
	RestrictedDirs = []string{"flags", "handlers", "models", "servers", "storage", "utils", "../"}
)

const (
	ApiBaseURL = "https://v3.football.api-sports.io/fixtures"
	// Your API key for the API Sports service
	ApiKey = "3d09bb9f390f3766a0a694f6dc764296"
)

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

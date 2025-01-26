package main

import (
	"Aitu-Bet/config"
	"Aitu-Bet/flags"
	"Aitu-Bet/internal/servers"
	"Aitu-Bet/logging"
	"fmt"
	"os"
)

func main() {
	flags.Setup()
	if err := logging.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	server := servers.NewServer()
	server.Start(config.Port)
}

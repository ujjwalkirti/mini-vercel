package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"reverse-proxy/internal/app"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize and run the application
	application := app.New()
	if err := application.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/app"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	} else {
		log.Println("Using .env file")
	}

	application := app.New()
	application.Run()
}

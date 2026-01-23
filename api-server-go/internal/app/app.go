package app

import (
	"log"
	"net/http"
	"os"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/db"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/router"
)

type App struct {
	server *http.Server
}

func New() *App {
	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	r := router.New(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	return &App{server: server}
}

func (a *App) Run() {
	log.Println("Server starting on", a.server.Addr)
	log.Fatal(a.server.ListenAndServe())
}

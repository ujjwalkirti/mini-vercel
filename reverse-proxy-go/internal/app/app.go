package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"reverse-proxy/internal/config"
	"reverse-proxy/internal/db"
	"reverse-proxy/internal/router"
)

type App struct {
	config *config.Config
	db     *sql.DB
	server *http.Server
}

func New() *App {
	return &App{
		config: config.Load(),
	}
}

func (a *App) Run() error {
	// Connect to database
	database, err := db.Connect(a.config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	a.db = database
	defer a.db.Close()

	log.Println("Database connection established")

	// Initialize router
	r := router.New(a.db, a.config.R2PublicURL)

	// Create HTTP server
	a.server = &http.Server{
		Addr:         ":" + a.config.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Reverse proxy server starting on port %s", a.config.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited gracefully")
	return nil
}

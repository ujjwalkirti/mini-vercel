package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/client"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/db"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/kafka/consumer"
	repository "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/deployment"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/router"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/deployment"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/logs"
)

type App struct {
	server      *http.Server
	cancelKafka context.CancelFunc
}

func New() *App {
	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	if err := config.InitSupabase(); err != nil {
		log.Fatal("Failed to initialize Supabase:", err)
	}

	// Initialize ClickHouse repository for logs
	logRepo, err := client.NewLogRepository(client.ClickHouseRepository)
	if err != nil {
		log.Fatal("Failed to create log repository:", err)
	} else {
		log.Println("Log repository created successfully")
	}

	// Initialize services
	logSvc := logs.New(logRepo)
	deploymentRepo := repository.New(database)
	deploymentSvc := deployment.NewDeploymentService(deploymentRepo)

	// Initialize Kafka processor
	processor := consumer.NewProcessor(deploymentSvc, logSvc)

	// Start Kafka consumer in background
	kafkaConfig := config.LoadKafkaConfig()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		log.Println("Starting Kafka consumer...")
		consumer.StartConsumer(ctx, kafkaConfig, processor)
	}()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down Kafka consumer...")
		cancel()
	}()

	r := router.New(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	return &App{
		server:      server,
		cancelKafka: cancel,
	}
}

func (a *App) Run() {
	log.Println("Server starting on", a.server.Addr)
	log.Fatal(a.server.ListenAndServe())
}

func (a *App) Shutdown() {
	if a.cancelKafka != nil {
		log.Println("Cancelling Kafka consumer context...")
		a.cancelKafka()
	}
}

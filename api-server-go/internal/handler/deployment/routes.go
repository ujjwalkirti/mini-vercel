package deployment

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-chi/chi/v5"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/auth"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/client"
	appConfig "github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/middleware"
	repository "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/deployment"
	projectRepository "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/project"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/ecs"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/logs"
)

func Routes(db *sql.DB, jwks *auth.JWKSCache) chi.Router {
	r := chi.NewRouter()

	// Apply auth middleware to all deployment routes
	r.Use(middleware.AuthMiddleware(jwks))

	repository := repository.New(db)
	projectRepo := projectRepository.New(db)

	// Initialize ECS service
	var ecsService *ecs.Service
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Printf("Warning: Failed to load AWS config: %v", err)
	} else {
		// Parse subnets (comma-separated)
		ecsConfig := appConfig.GetECSConfig()
		subnets := []string{}
		if ecsConfig.Subnets != "" {
			subnets = strings.Split(ecsConfig.Subnets, ",")
			for i := range subnets {
				subnets[i] = strings.TrimSpace(subnets[i])
			}
		}

		ecsService = ecs.New(
			awsCfg,
			ecsConfig.Cluster,
			ecsConfig.TaskDefinition,
			subnets,
			ecsConfig.SecurityGroup,
			ecsConfig.AssignPublicIP,
			ecsConfig.LaunchType,
			ecsConfig.ImageName,
			ecsConfig.Count,
		)

		log.Printf("ECS service initialized with config.")
	}

	// Initialize logs service
	var logsService *logs.Service
	clickhouseCfg := appConfig.GetClickHouseConfig()
	if clickhouseCfg.Host != "" {
		logRepo, err := client.NewLogRepository(client.ClickHouseRepository)
		if err != nil {
			log.Printf("Warning: Failed to initialize ClickHouse client: %v", err)
		} else {
			logsService = logs.New(logRepo)
		}
	} else {
		log.Printf("Warning: ClickHouse not configured - logs endpoint will return empty logs")
	}

	h := NewHandler(repository, projectRepo, ecsService, logsService)

	// GET /projects/:projectId/deployments - Get all deployments for a project
	r.Get("/projects/{projectId}", h.GetDeploymentsByProject)

	// GET /deployments/:id - Get specific deployment
	r.Get("/deployments/{id}", h.GetDeployment)

	// POST /deploy - Create new deployment
	r.Post("/deploy", h.CreateDeployment)

	// GET /deployments/:id/logs - Get deployment logs
	r.Get("/deployments/{id}/logs", h.GetDeploymentLogs)

	return r
}

package deployment

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-chi/chi/v5"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/client"
	appConfig "github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/middleware"
	repository "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/deployment"
	projectRepository "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/project"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/ecs"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/logs"
)

func Routes(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	// Apply auth middleware to all deployment routes
	r.Use(middleware.AuthMiddleware)

	repository := repository.New(db)
	projectRepo := projectRepository.New(db)

	// Initialize ECS service
	var ecsService *ecs.Service
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Printf("Warning: Failed to load AWS config: %v", err)
	} else {
		// Parse subnets (comma-separated)
		subnets := []string{}
		if appConfig.ECS_SUBNETS != "" {
			subnets = strings.Split(appConfig.ECS_SUBNETS, ",")
			for i := range subnets {
				subnets[i] = strings.TrimSpace(subnets[i])
			}
		}

		ecsService = ecs.New(
			awsCfg,
			appConfig.ECS_CLUSTER,
			appConfig.ECS_TASK_DEFINITION,
			subnets,
			appConfig.ECS_SECURITY_GROUP,
			appConfig.ECS_ASSIGN_PUBLIC_IP,
			appConfig.ECS_LAUNCH_TYPE,
			appConfig.ECS_IMAGE_NAME,
			appConfig.ECS_COUNT,
		)
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

package consumer

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/IBM/sarama"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/buildlog"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/deployment"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/logs"
	"go.uber.org/zap"
)

type Processor struct {
	deploymentSvc *deployment.DeploymentService
	logSvc        *logs.Service
}

func NewProcessor(
	deploymentSvc *deployment.DeploymentService,
	logSvc *logs.Service,
) *Processor {
	return &Processor{
		deploymentSvc: deploymentSvc,
		logSvc:        logSvc,
	}
}

func (p *Processor) Process(msg *sarama.ConsumerMessage) error {
	if msg.Key == nil || msg.Value == nil {
		return nil // skip silently
	}

	var event buildlog.Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	ctx := context.Background()
	logLower := strings.ToLower(event.Log)

	// ---- status transitions ----
	if logLower == "info: starting build pipeline..." {
		if err := p.deploymentSvc.MarkInProgress(ctx, event.DeploymentID); err != nil {
			return err
		}
	}

	if logLower == "info: pipeline completed successfully." {
		if err := p.deploymentSvc.MarkReady(ctx, event.DeploymentID); err != nil {
			return err
		}
	}

	if strings.HasPrefix(logLower, "error:") &&
		strings.Contains(logLower, "pipeline failed") {
		if err := p.deploymentSvc.MarkFailed(ctx, event.DeploymentID); err != nil {
			return err
		}
	}

	// ---- always insert log ----
	if err := p.logSvc.Insert(ctx, event.DeploymentID, event.Log); err != nil {
		// Log the error but don't fail message processing
		// This prevents message reprocessing and allows other logs to continue
		logger := config.GetLogger()
		logger.Error("Failed to insert log",
			zap.String("deployment_id", event.DeploymentID),
			zap.String("log_preview", truncateString(event.Log, 100)),
			zap.Error(err))
	}

	return nil
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

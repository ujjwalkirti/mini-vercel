package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
	"go.uber.org/zap"
)

func StartConsumer(ctx context.Context, env config.KafkaConfig, processor *Processor) {
	logger := config.GetLogger()

	cfg, err := config.NewSaramaConfig(env)
	if err != nil {
		logger.Fatal("Failed to create Sarama config", zap.Error(err))
	}

	group, err := sarama.NewConsumerGroup(
		env.Brokers,
		"mini-vercel-build-logs-go",
		cfg,
	)
	if err != nil {
		logger.Fatal("Failed to create consumer group", zap.Error(err))
	}

	handler := NewHandler(processor)

	for {
		if err := group.Consume(ctx, []string{"mini-vercel-build-logs"}, handler); err != nil {
			logger.Error("Consumer group error", zap.Error(err))
		}
		if ctx.Err() != nil {
			return
		}
	}
}

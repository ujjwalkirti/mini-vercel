package consumer

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
)

func StartConsumer(ctx context.Context, env config.KafkaConfig, processor *Processor) {
	cfg, err := config.NewSaramaConfig(env)
	if err != nil {
		log.Fatal(err)
	}

	group, err := sarama.NewConsumerGroup(
		env.Brokers,
		"mini-vercel-build-logs-go",
		cfg,
	)
	if err != nil {
		log.Fatal(err)
	}

	handler := NewHandler(processor)

	for {
		if err := group.Consume(ctx, []string{"mini-vercel-build-logs"}, handler); err != nil {
			log.Println(err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

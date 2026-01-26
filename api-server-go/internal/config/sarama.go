package config

import (
	"time"

	"github.com/IBM/sarama"
	kafka_tls "github.com/ujjwalkirti/mini-vercel-api-server/internal/kafka/tls"
)

func NewSaramaConfig(env KafkaConfig) (*sarama.Config, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_6_0_0

	cfg.ClientID = env.ClientID

	cfg.Net.DialTimeout = 30 * time.Second
	cfg.Net.ReadTimeout = 30 * time.Second
	cfg.Net.WriteTimeout = 30 * time.Second

	cfg.Net.SASL.Enable = true
	cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	cfg.Net.SASL.User = env.Username
	cfg.Net.SASL.Password = env.Password

	// Consumer group configuration
	cfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Return.Errors = true

	tlsCfg, err := kafka_tls.NewTLSConfig()
	if err != nil {
		return nil, err
	}

	cfg.Net.TLS.Enable = true
	cfg.Net.TLS.Config = tlsCfg

	return cfg, nil
}

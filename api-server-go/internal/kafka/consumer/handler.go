package consumer

import (
	"github.com/IBM/sarama"
)

type Handler struct {
	pool      *WorkerPool
	processor *Processor
}

func NewHandler(processor *Processor) *Handler {
	return &Handler{
		pool:      NewWorkerPool(50),
		processor: processor,
	}
}

// Called when a new consumer session starts (rebalance)
func (h *Handler) Setup(s sarama.ConsumerGroupSession) error {
	// no-op for now
	return nil
}

// Called when a consumer session ends (rebalance / shutdown)
func (h *Handler) Cleanup(s sarama.ConsumerGroupSession) error {
	// no-op for now
	return nil
}

// Called once per partition
func (h *Handler) ConsumeClaim(
	s sarama.ConsumerGroupSession,
	c sarama.ConsumerGroupClaim,
) error {
	for msg := range c.Messages() {
		h.pool.Submit(func() {
			if err := h.processor.Process(msg); err == nil {
				s.MarkMessage(msg, "")
			}
		})
	}
	return nil
}

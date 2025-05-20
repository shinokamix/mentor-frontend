package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"rating/internal/domain/models"
	grpccleint "rating/internal/transport/grpc/client"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer     sarama.ConsumerGroup
	handler      *consumerHandler
	mentorClient *grpccleint.MentorClient
	log          *slog.Logger
}

type consumerHandler struct {
	ready     chan bool
	processor func(ctx context.Context, msg *models.ReviewEvent) error
	log       *slog.Logger
}

func NewConsumer(brokers []string, groupID string, mentorClient *grpccleint.MentorClient, logger *slog.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Создаем группу консьюмеров
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	// Инициализируем наш обработчик
	handler := &consumerHandler{
		ready: make(chan bool),
		log:   logger,
	}

	return &Consumer{
		consumer:     consumerGroup,
		handler:      handler,
		mentorClient: mentorClient,
		log:          logger,
	}, nil
}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	h.log.Info("cleanup consumer handler")
	return nil
}

func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		var event models.ReviewEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			h.log.Error("falied to unmarshal json", "error", err)
			continue
		}

		if err := h.processor(session.Context(), &event); err != nil {
			h.log.Error("processor error", "error", err, "event", event)
			continue
		}

		session.MarkMessage(msg, "")
	}

	return nil
}

func (c *Consumer) Run(ctx context.Context, topic string) {
	c.handler.processor = func(ctx context.Context, msg *models.ReviewEvent) error {
		return c.mentorClient.MethodMentorRating(ctx, msg.Action, msg.Email, msg.Score)
	}

	go func() {
		for {
			if err := c.consumer.Consume(ctx, []string{topic}, c.handler); err != nil {
				c.log.Error("failed to consume messages", "error", err)
			}

			// если контекст отменён - выходим
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-c.handler.ready
	c.log.Info("consumer started for topic", "topic", topic)
}

func (c *Consumer) Close() error {
	// Закрываем ConsumerGroup, когда хотим завершиться
	c.log.Info("closing consumer group")
	return c.consumer.Close()
}

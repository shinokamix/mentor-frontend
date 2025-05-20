package kafka

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"review/internal/domain/model"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
	logger   *slog.Logger
}

func NewProducer(brokers []string, topic string, logger *slog.Logger) (*Producer, error) {
	config := sarama.NewConfig()

	// Гарантия доставки
	config.Producer.RequiredAcks = sarama.WaitForAll // Ждем подтверждения от всех реплик
	config.Producer.Retry.Max = 3                    // 3 попытки при ошибках
	config.Producer.Idempotent = true                // режим идемпотентности
	config.Net.MaxOpenRequests = 1                   // Обязательно для идемпотентности
	config.Producer.Return.Successes = true          // Получаем подтверждения

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}, nil
}

func (p *Producer) SendReviewEvent(review *model.ReviewEvent) error {
	jsonData, err := json.Marshal(review)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", review.ID)),

		Value: sarama.ByteEncoder(jsonData),
	}

	partition, offset, err := p.producer.SendMessage(msg)

	if err != nil {
		p.logger.Error("fali sand messgae", "partiotion", partition, "offset", offset)
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

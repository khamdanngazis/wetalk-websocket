package kafka

import (
	"chat-websocket/internal/domain"
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaService struct {
	Writer *kafka.Writer
}

func NewKafkaService(brokers []string) domain.Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		Async:        false,
		RequiredAcks: kafka.RequireAll,
	}
	return &KafkaService{Writer: writer}
}

func (k *KafkaService) SendMessage(topic string, key, message []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: message,
	}

	err := k.Writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to send Kafka message: %v", err)
		return err
	}
	return nil
}

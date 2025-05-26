package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// RunMigrations creates default Kafka topics
func RunMigrations(brokers []string) error {
	defaultTopics := []string{
		"notification_batch",
		// Add more default topics here as needed
	}

	return CreateTopics(brokers, defaultTopics)
}

// KafkaProducer is an interface for publishing messages to Kafka
type KafkaProducer interface {
	Publish(ctx context.Context, topic string, data []byte) error
}

type kafkaProducer struct {
	producer *kafka.Producer
}

// NewKafkaProducer creates a new instance of KafkaProducer
func NewKafkaProducer(brokers []string) (KafkaProducer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers[0],
		"client.id":         "go-kafka-producer",
		"acks":              "all",
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	return &kafkaProducer{
		producer: producer,
	}, nil
}

// Publish publishes a message to a Kafka topic
func (k *kafkaProducer) Publish(ctx context.Context, topic string, data []byte) error {
	deliveryChan := make(chan kafka.Event)

	err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to produce message: %v", err)
	}

	return nil
}

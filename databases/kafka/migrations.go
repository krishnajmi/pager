package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// CreateTopics creates the specified Kafka topics if they don't exist
func CreateTopics(brokers []string, topics []string) error {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers[0],
	}

	adminClient, err := kafka.NewAdminClient(config)
	if err != nil {
		return fmt.Errorf("failed to create admin client: %v", err)
	}
	defer adminClient.Close()

	topicSpecs := make([]kafka.TopicSpecification, len(topics))
	for i, topic := range topics {
		topicSpecs[i] = kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	results, err := adminClient.CreateTopics(context.Background(), topicSpecs)
	if err != nil {
		return fmt.Errorf("failed to create topics: %v", err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Printf("Warning: failed to create topic %s: %v", result.Topic, result.Error.String())
		} else {
			log.Printf("Successfully created topic: %s", result.Topic)
		}
	}

	return nil
}

package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/kp/pager/common"
	"github.com/kp/pager/communicator"
)

// StartBatchConsumer starts a Kafka consumer for notification_batch topic
func StartBatchConsumer(brokers []string) {
	config := &kafka.ConfigMap{
		"bootstrap.servers":  brokers[0],
		"group.id":           "go-kafka-consumer",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "true",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	defer consumer.Close()

	topic := "notification_batch"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)

	// Graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := consumer.ReadMessage(100)
			if err == nil {
				if err := ProcessBatchMessages(context.Background(), msg); err != nil {
					fmt.Printf("Failed to process message: %v\n", err)
				}
			} else {
				// Only print real errors, not timeouts
				if err.(kafka.Error).Code() != kafka.ErrTimedOut {
					fmt.Printf("Error consuming message: %v\n", err)
				}
			}
		}
	}
}

func ProcessBatchMessages(ctx context.Context, message *kafka.Message) error {
	// Parse message
	var qMessage communicator.QMessage
	if err := json.Unmarshal(message.Value, &qMessage); err != nil {
		return fmt.Errorf("failed to parse message: %v", err)
	}

	notification := qMessage.GenericModel
	errChan := make(chan error, len(qMessage.Audiences))
	var wg sync.WaitGroup

	for _, audience := range qMessage.Audiences {
		wg.Add(1)
		go func(aud common.AudienceType) {
			defer wg.Done()
			notificationService := communicator.NewCommunicatornNotificationSevice(notification, aud.Email, aud.Context)
			commService := communicator.NewCommunicationService(ctx, notificationService)

			if err := commService.Run(ctx); err != nil {
				errChan <- fmt.Errorf("failed to process audience %v: %w", aud, err)
			}
		}(audience)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect any errors
	for err := range errChan {
		if err != nil {
			slog.Error("batch processing error",
				slog.String("error", err.Error()),
			)
		}
	}

	return nil
}

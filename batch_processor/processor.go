package batchprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/kp/pager/common"
	"github.com/kp/pager/communicator"
	"github.com/kp/pager/databases/kafka"
	"github.com/rs/xid"

	log "github.com/sirupsen/logrus"
)

func NewBatchProcessor(ctx context.Context, audiences []common.AudienceType, model communicator.NotificationType, topicName string, kafkaProducer kafka.KafkaProducer) BatchProcessor {
	return &BatchChannelBased{
		Model:         model,
		TopicName:     topicName,
		Audiences:     audiences,
		KafkaProducer: kafkaProducer,
	}
}

// TODO: pageId or batchID for logging
func (batch *BatchChannelBased) Process(ctx context.Context) (err error) {
	audiences := batch.Audiences
	batchSize := 5
	var wg sync.WaitGroup
	concurrentGoroutines := 20
	semaphore := make(chan struct{}, concurrentGoroutines) // Adjusted to use concurrentGoroutines
	for i := 0; i < len(audiences); i += batchSize {
		var batchAudience []common.AudienceType
		if i+batchSize < len(audiences) {
			batchAudience = audiences[i : i+batchSize]
		} else {
			batchAudience = audiences[i:]
		}
		semaphore <- struct{}{}
		wg.Add(1)
		go func(startIndex int, batchAudience []common.AudienceType) {
			defer func() {
				<-semaphore
				wg.Done()
			}()
			if startIndex+batchSize < len(audiences) {
				err = batch.sendBatchToQueue(ctx, audiences[startIndex:startIndex+batchSize])
			} else {
				err = batch.sendBatchToQueue(ctx, audiences[startIndex:])
			}
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Errorln("CreateBatchAndSendToQueue:ErrorProcessingBatch")
			}
		}(i, batchAudience)
	}
	wg.Wait()
	return
}

func (batch *BatchChannelBased) sendBatchToQueue(c context.Context, audiences []common.AudienceType) (err error) {
	batchID := xid.New().String()
	kafkaMessage := communicator.QMessage{
		BatchID:      batchID,
		Audiences:    audiences,
		GenericModel: batch.Model,
	}

	messageBytes, err := json.Marshal(kafkaMessage)
	if err != nil {
		log.WithFields(log.Fields{
			"error":         err,
			"kafka_message": kafkaMessage,
		}).Errorln("sendBatchToQueueMarshalFailed")
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = batch.KafkaProducer.Publish(c, batch.TopicName, messageBytes)
	if err != nil {
		log.WithFields(log.Fields{
			"error":        err,
			"publish_data": string(messageBytes),
		}).Errorln("BatchFailedToPublish")
	}
	return
}

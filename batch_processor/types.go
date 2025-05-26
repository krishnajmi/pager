package batchprocessor

import (
	"context"

	"github.com/kp/pager/common"
	"github.com/kp/pager/communicator"
	"github.com/kp/pager/databases/kafka"
)

type BatchProcessor interface {
	Process(ctx context.Context) error
}

type BatchChannelBased struct {
	TopicName     string
	Model         communicator.NotificationType
	Audiences     []common.AudienceType
	KafkaProducer kafka.KafkaProducer
}

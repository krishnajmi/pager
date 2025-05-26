package notification

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	batchprocessor "github.com/kp/pager/batch_processor"
	"github.com/kp/pager/communicator"
	"github.com/kp/pager/databases/kafka"
)

func generateUniqueID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func NewNotificationService(ctx context.Context, notificationRequest NotificationRequestType, sessionService NotificationSessionService, kafkaProducer kafka.KafkaProducer) NotificationService {
	return &Notification{
		TemplateID:                 notificationRequest.TemplateID,
		Audiences:                  notificationRequest.Audiences,
		NotificationSessionService: sessionService,
		KafkaProducer:              kafkaProducer,
	}
}

func (c *Notification) SendNotification(ctx context.Context) (*Notification, error) {
	// Create notification session with unique request_id
	session := NotificationSession{
		RequestID:     generateUniqueID(),
		Status:        NotifcationSessionStatusCreated,
		TotalAudience: len(c.Audiences),
		TotalSent:     len(c.Audiences),
		TemplateID:    c.TemplateID,
		TotalSuccess:  0,
		BatchProcessor: &batchprocessor.BatchChannelBased{
			TopicName: "notification_queue",
			Model:     communicator.NotificationType{},
			Audiences: c.Audiences,
		},
	}

	// Save notification session to database for tracking and auditing purposes
	sessionID, err := c.NotificationSessionService.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	notificationType := communicator.NotificationType{
		TemplateID: c.TemplateID,
		SessionID:  sessionID,
		RequestId:  session.RequestID,
	}

	session.BatchProcessor = batchprocessor.NewBatchProcessor(
		ctx,
		c.Audiences,
		notificationType,
		"notification_batch",
		c.KafkaProducer,
	)
	// Process notification batch asynchronously through the batch processor
	if err := session.BatchProcessor.Process(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Notification) FetchNotificationSessions(ctx context.Context, campaignID string) ([]NotificationSession, error) {
	return nil, nil
}

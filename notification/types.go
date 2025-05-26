package notification

import (
	"context"
	"time"

	batchprocessor "github.com/kp/pager/batch_processor"
	"github.com/kp/pager/common"
	"github.com/kp/pager/databases/kafka"
)

type NotificationRequestType struct {
	Notification
	UserName   string `json:"user_name"`
	TemplateID int64  `json:"template_id"`
}

type Notification struct {
	ID                         string                `json:"id"`
	TemplateID                 int64                 `json:"template_id"`
	Audiences                  []common.AudienceType `json:"audiences"`
	CreatedAt                  time.Time             `json:"created_at"`
	UpdatedAt                  time.Time             `json:"updated_at"`
	NotificationSessionService NotificationSessionService
	KafkaProducer              kafka.KafkaProducer
}

type NotificationService interface {
	SendNotification(ctx context.Context) (*Notification, error)
}

type NotificationSessionService interface {
	Create(ctx context.Context, session NotificationSession) (int64, error)
}

type NotificationSession struct {
	ID            string    `json:"id"`
	RequestID     string    `json:"request_id"`
	Status        string    `json:"status"`
	TemplateID    int64     `json:"template_id"`
	TotalAudience int       `json:"total_audience"`
	TotalSent     int       `json:"total_sent"`
	TotalSuccess  int       `json:"total_success"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	batchprocessor.BatchProcessor
}

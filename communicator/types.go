package communicator

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/kp/pager/common"
)

type NotificationType struct {
	To         string            `json:"to"`
	TemplateID int64             `json:"template_id"`
	RequestId  string            `json:"request_id"`
	Subject    string            `json:"subject"`
	Body       string            `json:"body"`
	Context    map[string]string `json:"context"`
	SessionID  int64             `json:"session_id"`
	LogID      int64             `json:"log_id"`
}

type CommunicationHandler interface {
	Run(ctx context.Context) error
}

type CommunicatorNotificationHandler interface {
	Save(ctx context.Context) error
	Validate(ctx context.Context) error
	Prepare(ctx context.Context) (interface{}, error)
	Send(ctx context.Context, payload interface{}) error
}

type communicator struct {
	db                  *gorm.DB
	Notification        NotificationType
	NotificationHanlder CommunicatorNotificationHandler
}

func NewCommunicatornNotificationSevice(notification NotificationType, to string, context map[string]string) CommunicatorNotificationHandler {
	return &NotificationType{
		To:         to,
		TemplateID: notification.TemplateID,
		RequestId:  notification.RequestId,
		Context:    notification.Context,
		SessionID:  notification.SessionID,
	}
}

func NewCommunicationService(ctx context.Context, commNotificationHandler CommunicatorNotificationHandler) CommunicationHandler {
	return &communicator{NotificationHanlder: commNotificationHandler}
}

type QMessage struct {
	BatchID      string                `json:"batch_id"`
	GenericModel NotificationType      `json:"model"`
	Audiences    []common.AudienceType `json:"audiences"`
}

type NotificationPayload struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Name    string `json:"name"`
}

package notification

import (
	"context"
	"errors"

	"github.com/jinzhu/gorm"
	models "github.com/kp/pager/notification/models"
)

type notificationSessionService struct {
	db *gorm.DB
}

func NewNotificationSessionService(db *gorm.DB) NotificationSessionService {
	return &notificationSessionService{
		db: db,
	}
}

func (s *notificationSessionService) Create(ctx context.Context, session NotificationSession) (int64, error) {
	if session.RequestID == "" {
		return 0, errors.New("request_id cannot be empty")
	}

	entry, err := models.NewNotificationSessionEntry(
		ctx,
		s.db,
		session.Status,
		session.RequestID,
		session.TotalAudience,
		session.TemplateID,
	)
	return entry.ID, err
}

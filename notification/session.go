package notification

import (
	"context"
	"errors"
	"strconv"

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

func (s *notificationSessionService) Update(ctx context.Context, session NotificationSession) error {
	if session.ID == "" {
		return errors.New("id cannot be empty")
	}

	id, err := strconv.ParseInt(session.ID, 10, 64)
	if err != nil {
		return errors.New("invalid ID format")
	}

	model := models.NotificationSession{
		ID:            id,
		RequestID:     session.RequestID,
		Status:        session.Status,
		TemplateID:    session.TemplateID,
		TotalAudience: session.TotalAudience,
		UpdatedAt:     session.UpdatedAt,
	}

	updates := map[string]interface{}{
		"status":         model.Status,
		"template_id":    model.TemplateID,
		"total_audience": model.TotalAudience,
		"updated_at":     model.UpdatedAt,
	}

	return s.db.Model(&model).Updates(updates).Error
}

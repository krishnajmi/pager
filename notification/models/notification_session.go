package notification

import (
	"context"
	"time"

	"github.com/kp/pager/databases/sql"
)

const NotificationSessionTableName = "notification_session"

type NotificationSession struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	TemplateID    int64     `gorm:"column:template_id"`
	RequestID     string    `gorm:"column:request_id"`
	TotalAudience int       `gorm:"column:total_audience"`
	Status        string    `gorm:"column:status"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (NotificationSession) TableName() string {
	return NotificationSessionTableName
}

func NewNotificationSessionEntry(ctx context.Context, tx interface{}, status, requestID string, totalAudience int, templateID int64) (*NotificationSession, error) {
	database := sql.GetOrmQuearyable(ctx, tx)
	entry := NotificationSession{
		TotalAudience: totalAudience,
		TemplateID:    templateID,
		RequestID:     requestID,
		Status:        status,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err := database.Create(&entry).Error
	return &entry, err
}

func GetNotificationSessionByID(ctx context.Context, tx interface{}, id int64) (*NotificationSession, error) {
	db := sql.GetOrmQuearyable(ctx, tx)
	entry := NotificationSession{}
	err := db.First(&entry, id).Error
	return &entry, err
}

func GetAllNotificationSessions(ctx context.Context, tx interface{}, limit, offset int) ([]NotificationSession, error) {
	var sessions []NotificationSession
	db := sql.GetOrmQuearyable(ctx, tx)
	query := db.Model(&NotificationSession{})

	if limit > 0 || offset > 0 {
		if limit > 0 {
			query = query.Limit(limit)
		}
		if offset > 0 {
			query = query.Offset(offset)
		}
	}

	err := query.Find(&sessions).Error
	return sessions, err
}

func (session NotificationSession) Save(ctx context.Context, tx interface{}) error {
	db := sql.GetOrmQuearyable(ctx, tx)
	return db.Save(session).Error
}

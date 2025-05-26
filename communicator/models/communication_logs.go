package models

import (
	"context"
	"time"

	"github.com/kp/pager/databases/sql"
)

const CommunicationLogsTableName = "communication_logs"

type CommunicationLogs struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Email      string    `gorm:"column:email"`
	TemplateID int64     `gorm:"column:template_id"`
	RequestID  string    `gorm:"column:request_id;index"`
	Status     string    `gorm:"column:status"`
	Payload    string    `gorm:"column:payload;type:text"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (CommunicationLogs) TableName() string {
	return CommunicationLogsTableName
}

func NewCommunicationLogEntry(ctx context.Context, tx interface{}, email string, templateID int64, requestID string) (*CommunicationLogs, error) {
	database := sql.GetOrmQuearyable(ctx, tx)
	entry := CommunicationLogs{
		Email:      email,
		TemplateID: templateID,
		RequestID:  requestID,
		Status:     "created",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := database.Create(&entry).Error
	return &entry, err
}

func GetCommunicationLogByID(ctx context.Context, tx interface{}, id int64) (*CommunicationLogs, error) {
	db := sql.GetOrmQuearyable(ctx, tx)
	entry := CommunicationLogs{}
	err := db.First(&entry, id).Error
	return &entry, err
}

func GetAllCommunicationLogsByRequestID(ctx context.Context, tx interface{}, requestID string) ([]CommunicationLogs, error) {
	var logs []CommunicationLogs
	db := sql.GetOrmQuearyable(ctx, tx)
	err := db.Where("request_id = ?", requestID).Find(&logs).Error
	return logs, err
}

func GetAllCommunicationLogs(ctx context.Context, tx interface{}, limit, offset int) ([]CommunicationLogs, error) {
	var logs []CommunicationLogs
	db := sql.GetOrmQuearyable(ctx, tx)
	query := db.Model(&CommunicationLogs{})

	if limit > 0 || offset > 0 {
		if limit > 0 {
			query = query.Limit(limit)
		}
		if offset > 0 {
			query = query.Offset(offset)
		}
	}

	err := query.Find(&logs).Error
	return logs, err
}

func (log CommunicationLogs) Save(ctx context.Context, tx interface{}) error {
	db := sql.GetOrmQuearyable(ctx, tx)
	return db.Save(log).Error
}

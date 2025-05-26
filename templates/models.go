package templates

import (
	"context"
	"time"

	"github.com/kp/pager/databases/sql"
)

const TemplateTableName = "notification_template"

type NotificationTemplate struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	Name        string    `gorm:"column:name;unique"`
	Subject     string    `gorm:"column:subject"`
	Content     string    `gorm:"column:content;type:text"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (NotificationTemplate) TableName() string {
	return TemplateTableName
}

func NewTemplateEntry(ctx context.Context, tx interface{}, name, subject, content string) (*NotificationTemplate, error) {
	database := sql.GetOrmQuearyable(ctx, tx)
	entry := NotificationTemplate{
		Name:      name,
		Subject:   subject,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := database.Create(&entry).Error
	return &entry, err
}

func GetTemplateByID(ctx context.Context, tx interface{}, id int64) (*NotificationTemplate, error) {
	db := sql.GetOrmQuearyable(ctx, tx)
	entry := NotificationTemplate{}
	err := db.First(&entry, id).Error
	return &entry, err
}

func GetAllTemplates(ctx context.Context, tx interface{}, limit, offset int) ([]NotificationTemplate, error) {
	var templates []NotificationTemplate
	db := sql.GetOrmQuearyable(ctx, tx)
	query := db.Model(&NotificationTemplate{})

	if limit > 0 || offset > 0 {
		if limit > 0 {
			query = query.Limit(limit)
		}
		if offset > 0 {
			query = query.Offset(offset)
		}
	}

	err := query.Find(&templates).Error
	return templates, err
}

func (template NotificationTemplate) Save(ctx context.Context, tx interface{}) error {
	db := sql.GetOrmQuearyable(ctx, tx)
	return db.Save(template).Error
}

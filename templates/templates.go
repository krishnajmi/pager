package templates

import (
	"context"
	"log/slog"
	"time"

	"github.com/jinzhu/gorm"
)

type Template struct {
	ID        int64     `gorm:"primaryKey"`
	Name      string    `gorm:"unique;not null"`
	Subject   string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type templateManager struct {
	db *gorm.DB
}

func (s *templateManager) CreateTemplate(ctx context.Context, name, subject, content string) (*Template, error) {
	template, err := NewTemplateEntry(ctx, nil, name, subject, content)
	if err != nil {
		slog.Error("createTemplate:unableToCreateTemplate", slog.Any("error", err))
		return nil, err
	}
	return &Template{
		ID:        template.ID,
		Name:      template.Name,
		Subject:   template.Subject,
		Content:   template.Content,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}, nil
}

func (s *templateManager) UpdateTemplate(ctx context.Context, id int64, name, subject, content string) (*Template, error) {
	template := &NotificationTemplate{
		ID:        int64(id),
		Name:      name,
		Subject:   subject,
		Content:   content,
		UpdatedAt: time.Now(),
	}
	err := template.Save(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Template{
		ID:        template.ID,
		Name:      template.Name,
		Subject:   template.Subject,
		Content:   template.Content,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}, nil
}

func (s *templateManager) GetTemplate(ctx context.Context, id int64) (*Template, error) {
	template, err := GetTemplateByID(ctx, nil, int64(id))
	if err != nil {
		return nil, err
	}
	return &Template{
		ID:        template.ID,
		Name:      template.Name,
		Subject:   template.Subject,
		Content:   template.Content,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}, nil
}

func (s *templateManager) GetAllTemplates(ctx context.Context) ([]Template, error) {
	notificationTemplates, err := GetAllTemplates(ctx, nil, 0, 0)
	if err != nil {
		return nil, err
	}
	var templates []Template
	for _, t := range notificationTemplates {
		templates = append(templates, Template{
			ID:        t.ID,
			Name:      t.Name,
			Subject:   t.Subject,
			Content:   t.Content,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}
	return templates, nil
}

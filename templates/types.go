package templates

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/kp/pager/common"
)

type TemplateService interface {
	CreateTemplate(ctx context.Context, name, subject, content string) (*Template, error)
	UpdateTemplate(ctx context.Context, id int64, name, subject, content string) (*Template, error)
	GetTemplate(ctx context.Context, id int64) (*Template, error)
	GetAllTemplates(ctx context.Context) ([]Template, error)
}

func NewTemplateService(db *gorm.DB) TemplateService {
	return &templateManager{db: db}
}

type TemplateRequest struct {
	Name    string `json:"name" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Content string `json:"content" binding:"required"`
}

var TemplateResponse struct {
	common.Response
}

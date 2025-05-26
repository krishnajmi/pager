package communicator

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kp/pager/communicator/models"
	"github.com/kp/pager/databases/sql"
	template "github.com/kp/pager/templates"
)

func (n *NotificationType) Save(ctx context.Context) error {
	entry, err := models.NewCommunicationLogEntry(ctx, nil,
		n.To, n.TemplateID, n.RequestId)
	if err != nil {
		return fmt.Errorf("failed to save communication log: %v", err)
	}

	n.LogID = entry.ID
	return nil
}

func (n *NotificationType) Validate(ctx context.Context) error {
	if strings.EqualFold(n.To, "") {
		return fmt.Errorf("recipient (To) field cannot be empty")
	}
	return nil
}

func (n *NotificationType) Prepare(ctx context.Context) (interface{}, error) {
	var payload NotificationPayload
	templateService := template.NewTemplateService(sql.PagerOrm)
	templateData, err := templateService.GetTemplate(ctx, n.TemplateID)
	if err != nil {
		slog.Error("prepare:failedToGetTemplate",
			slog.Int64("template_id", n.TemplateID),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed to get template: %v", err)
	}

	payload.Body = templateData.Content
	payload.Subject = templateData.Subject
	payload.Name = templateData.Name
	return payload, err
}

func (n *NotificationType) Send(ctx context.Context, payload interface{}) error {
	// Fetch existing log entry
	tx := sql.PagerOrm.Begin()
	defer tx.Rollback()
	entry, err := models.GetCommunicationLogByID(ctx, tx, n.LogID)
	if err != nil {
		return fmt.Errorf("failed to fetch communication log: %v", err)
	}

	// Update status and payload
	entry.Status = "sent"
	entry.Payload = fmt.Sprintf("%v", payload)

	// Save updated entry
	if err := entry.Save(ctx, tx); err != nil {
		slog.Error("send:failedToUpdateCommunicationLog",
			slog.Int64("log_id", n.LogID),
			slog.Any("error", err))
		return fmt.Errorf("failed to update communication log: %v", err)
	}

	tx.Commit()
	return err
}

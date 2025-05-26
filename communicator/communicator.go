package communicator

import (
	"context"
	"fmt"
)

func (c *communicator) Run(ctx context.Context) error {
	// save the notification
	if err := c.NotificationHanlder.Save(ctx); err != nil {
		return fmt.Errorf("save failed: %v", err)
	}

	// Validate the notification
	if err := c.NotificationHanlder.Validate(ctx); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	// Prepare the notification
	payload, err := c.NotificationHanlder.Prepare(ctx)
	if err != nil {
		return fmt.Errorf("preparation failed: %v", err)
	}

	// Send the notification
	if err := c.NotificationHanlder.Send(ctx, payload); err != nil {
		return fmt.Errorf("sending failed: %v", err)
	}

	return nil
}

package notification

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kp/pager/databases/kafka"
	"github.com/kp/pager/databases/sql"
)

type NotificationController struct {
	NotificationService NotificationService
	kafkaProducer       kafka.KafkaProducer
}

func NewNotificationController(kafkaProducer kafka.KafkaProducer) *NotificationController {
	return &NotificationController{
		kafkaProducer: kafkaProducer,
	}
}

func (c *NotificationController) SendNotification(ctx *gin.Context) {
	var notificationRequest NotificationRequestType
	if err := ctx.ShouldBindJSON(&notificationRequest); err != nil {
		slog.Error("sendNotificationView:unableToBindJSON", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notificationRequest.UserName = ctx.GetString("username")
	notificationSessionService := NewNotificationSessionService(sql.PagerOrm)
	notificationService := NewNotificationService(ctx, notificationRequest, notificationSessionService, c.kafkaProducer)
	notificationData, err := notificationService.SendNotification(ctx)
	if err != nil {
		slog.Error("sendNotificationView:unableToSendNotification",
			slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": false,
			"msg":    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"error":  nil,
		"status": true,
		"msg":    "Notification sent successfully",
		"data":   notificationData,
	})
}

package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kp/pager/databases/kafka"
	"github.com/kp/pager/notification"
)

func NotificationRouterGroup(servicePrefix string, middlewares ...gin.HandlerFunc) RouterGroup {
	return RouterGroup{
		Prefix:      servicePrefix,
		Routes:      notificationRoutes(servicePrefix),
		Middlewares: middlewares}
}

func notificationRoutes(prefix string) []Route {
	// Initialize controllers

	brokers := []string{"localhost:9092"} // Replace with your Kafka broker addresses
	kafkaProducer, err := kafka.NewKafkaProducer(brokers)
	if err != nil {
		// Handle the error appropriately
		panic(err)
	}
	notificationCtrl := notification.NewNotificationController(kafkaProducer)
	return []Route{
		newRoute(http.MethodPost, "/trigger/", notificationCtrl.SendNotification, prefix),
		newRoute(http.MethodGet, "/sessions/:id/", notificationCtrl.GetCampaignSessions, prefix),
	}
}

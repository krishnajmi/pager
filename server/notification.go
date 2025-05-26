package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kp/pager/databases/kafka"
	login "github.com/kp/pager/login"
	"github.com/kp/pager/notification"
)

func NotificationRouterGroup(servicePrefix string, brokers []string, middlewares ...gin.HandlerFunc) RouterGroup {
	return RouterGroup{
		Prefix:      servicePrefix,
		Routes:      notificationRoutes(servicePrefix, brokers),
		Middlewares: middlewares}
}

func notificationRoutes(prefix string, brokers []string) []Route {
	// Initialize controllers

	kafkaProducer, err := kafka.NewKafkaProducer(brokers)
	if err != nil {
		panic(err)
	}
	notificationCtrl := notification.NewNotificationController(kafkaProducer)
	return []Route{
		newRoute(http.MethodPost, "/trigger/", notificationCtrl.SendNotification, prefix, login.PagerAdminAccess),
	}
}

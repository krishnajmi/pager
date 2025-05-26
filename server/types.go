package server

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

// each rest end point with associated permission key
type Route struct {
	Method        string
	Path          string
	Handler       gin.HandlerFunc
	PermissionKey string
}

// Router options to assign routes
type RouterOpts func(*gin.RouterGroup)

// Server options to
type ServerOpts func(*gin.Engine)

// default middlewares for service
var defaultMiddlewares []gin.HandlerFunc

// stores permission keys for each API
var permissionKeys sync.Map

// default server timeout, in seconds
var defaultServerTimeout time.Duration = 60 * time.Second

// Specific router/API group details
type RouterGroup struct {
	Prefix         string
	Routes         []Route
	Middlewares    []gin.HandlerFunc
	HealthEndPoint gin.HandlerFunc
}

func setContextValues() gin.HandlerFunc {
	return func(c *gin.Context) {
		//set request id
		reqID := c.Request.Header.Get("X-Request-Id")
		rsessionID := c.Request.Header.Get("X-RSESSIONID")
		c.Set("request_id", reqID)
		c.Set("rsessionid", rsessionID)

		//set logger for this request
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With(
			slog.String("request_id", reqID),
			slog.String("rsessionid", rsessionID),
			slog.String("ip_address", getClientIP(c)),
		)
		c.Set("logger", logger)
		c.Next()
	}
}

func getClientIP(c *gin.Context) string {

	clientIP := c.Request.Header.Get("X-Forwarded-For")
	if len(clientIP) == 0 {
		clientIP = c.Request.Header.Get("X-Real-IP")
	}
	if len(clientIP) == 0 {
		clientIP = c.Request.RemoteAddr
	}
	if strings.Contains(clientIP, ",") {
		clientIP = strings.Split(clientIP, ",")[0]
	}

	return clientIP
}

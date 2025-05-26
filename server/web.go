package server

import (
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

// intializes common settings for each service
func init() {
	defaultMiddlewares = []gin.HandlerFunc{setContextValues(), gin.CustomRecoveryWithWriter(os.Stdout, func(c *gin.Context, recovered interface{}) {
		traceMsg := strings.ReplaceAll(string(debug.Stack()), "\n", " ")
		c.Value("logger").(*slog.Logger).Error("error occurred",
			slog.String("request_id", c.Value("request_id").(string)),
			slog.String("rsession_id", c.Value("rsession_id").(string)),
			slog.String("url", c.Request.RequestURI),
			slog.Int("response_status", c.Writer.Status()),
			slog.String("ip_address", getClientIP(c)),
			slog.String("method", c.Request.Method),
			slog.String("stack_trace", traceMsg),
			slog.Int("response_code", http.StatusInternalServerError))
		c.AbortWithStatus(http.StatusInternalServerError)
	})}
}

func CreateRoutes(g ...RouterGroup) ServerOpts {
	return func(e *gin.Engine) {
		if defaultServerTimeout > 0 {
			defaultMiddlewares = append(defaultMiddlewares, setServerTimeOut(defaultServerTimeout))
		}
		for _, group := range g {
			if len(group.Middlewares) < 1 {
				group.Middlewares = []gin.HandlerFunc{}
			}
			middlewares := append(defaultMiddlewares, group.Middlewares...)
			routerGroup := e.Group(group.Prefix, middlewares...)
			if group.HealthEndPoint != nil {
				routerGroup.Handle(http.MethodGet, "/health/", group.HealthEndPoint)
			}
			routeGenerator := WithRoutes(group.Routes)
			routeGenerator(routerGroup)
		}
	}
}

func timeoutResponse(c *gin.Context) {
	c.String(http.StatusRequestTimeout, "server timeout")
}

func setServerTimeOut(tout time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(tout),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse),
	)
}

func WithRoutes(routes []Route) RouterOpts {
	return func(r *gin.RouterGroup) {
		for _, route := range routes {
			r.Handle(route.Method, route.Path, route.Handler)
		}
	}
}

func WithTimeOut(duration time.Duration) ServerOpts {
	return func(e *gin.Engine) {
		defaultServerTimeout = duration
	}
}

func InitServer(commonMiddlewares []gin.HandlerFunc, opts ...ServerOpts) *gin.Engine {
	appEnv := os.Getenv("APP_ENV")
	if strings.EqualFold(appEnv, "production") || strings.EqualFold(appEnv, "prod") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	if commonMiddlewares == nil {
		commonMiddlewares = []gin.HandlerFunc{}
	}
	r.GET("/health/", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	defaultMiddlewares = append(commonMiddlewares, defaultMiddlewares...)
	for _, opt := range opts {
		opt(r)
	}
	//set gin mode
	return r
}

func newRoute(method, path string, handler gin.HandlerFunc, servicePrefix string, permissions ...string) Route {
	route := Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}

	// Cache permissions for this route (only for static paths)
	if len(permissions) > 0 {
		fullPath := servicePrefix + path
		CachePermissions(method, fullPath, permissions)
	}
	return route
}

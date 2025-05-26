package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kp/pager/server"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testHTTPServer *http.Server
)

type mockServer struct {
	mock.Mock
}

func (m *mockServer) ListenAndServe() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServer) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Mock DB implementation
type mockDB struct {
	gorm.DB
}

func (m *mockDB) AutoMigrate(values ...interface{}) error {
	return nil
}

func (m *mockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	return &gorm.DB{}
}

// Mock controller implementations
type mockTemplateController struct {
	db *mockDB
}

func (m *mockTemplateController) GetTemplates(c *gin.Context) {
	c.JSON(200, gin.H{"data": []string{}})
}

type mockNotificationController struct {
	channels []string
}

func (m *mockNotificationController) GetNotifications(c *gin.Context) {
	c.JSON(200, gin.H{"data": []string{}})
}

type mockAuthController struct {
	db *mockDB
}

func (m *mockAuthController) Login(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}

func TestServerRoutes(t *testing.T) {
	// Setup test router with mock dependencies
	router := gin.New()
	middlewares := []gin.HandlerFunc{}

	db := &mockDB{}
	channels := []string{"email"}

	server.InitServer(middlewares,
		server.WithTimeOut(1*time.Second), // Shorter timeout for tests
		server.CreateRoutes(
			server.TemplateRouterGroup("/templates", &db.DB, middlewares...),
			server.NotificationRouterGroup("/notifications", channels, middlewares...),
			server.AuthRouterGroup("/auth", &db.DB, middlewares...),
		),
	)

	tests := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/templates", 200},
		{"GET", "/notifications", 200},
		{"POST", "/auth/login", 200},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.status, w.Code)
		})
	}
}

func TestApiCmd_Run(t *testing.T) {
	t.Run("normal startup and shutdown", func(t *testing.T) {
		// Use random available port
		ln, err := net.Listen("tcp", ":0")
		require.NoError(t, err)
		port := ln.Addr().(*net.TCPAddr).Port
		ln.Close()

		// Setup test signal channel
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT)
		defer signal.Stop(sigChan)

		// Run command in goroutine
		go func() {
			os.Args = []string{"", "--port", strconv.Itoa(port)}
			apiCmd.Run(&cobra.Command{}, []string{})
		}()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Verify server is running
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Send shutdown signal
		sigChan <- syscall.SIGINT

		// Wait for shutdown
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("startup error", func(t *testing.T) {
		// Force port conflict
		ln, err := net.Listen("tcp", ":0")
		require.NoError(t, err)
		defer ln.Close()
		port := ln.Addr().(*net.TCPAddr).Port

		// Should fail to start on occupied port
		os.Args = []string{"", "--port", strconv.Itoa(port)}
		assert.Panics(t, func() {
			apiCmd.Run(&cobra.Command{}, []string{})
		})
	})
}

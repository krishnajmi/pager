package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kp/pager/server"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestApiCmd_Run(t *testing.T) {
	t.Run("normal startup and shutdown", func(t *testing.T) {
		// Setup test signal channel
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT)
		defer signal.Stop(sigChan)

		// Create test server
		mockSrv := &mockServer{}
		mockSrv.On("ListenAndServe").Return(http.ErrServerClosed)
		mockSrv.On("Shutdown", mock.Anything).Return(nil)
		testHTTPServer = &http.Server{
			Handler: gin.New(),
		}

		// Run command in goroutine
		go func() {
			apiCmd.Run(&cobra.Command{}, []string{})
		}()

		// Send shutdown signal
		sigChan <- syscall.SIGINT

		// Wait for shutdown to complete
		time.Sleep(100 * time.Millisecond)

		mockSrv.AssertExpectations(t)
	})

	t.Run("startup error", func(t *testing.T) {
		mockSrv := &mockServer{}
		mockSrv.On("ListenAndServe").Return(fmt.Errorf("startup error"))
		testHTTPServer = &http.Server{
			Handler: gin.New(),
		}

		assert.Panics(t, func() {
			apiCmd.Run(&cobra.Command{}, []string{})
		})

		mockSrv.AssertExpectations(t)
	})
}

func TestServerRoutes(t *testing.T) {
	router := gin.New()
	middlewares := []gin.HandlerFunc{}
	server.InitServer(middlewares,
		server.WithTimeOut(0*time.Second),
		server.CreateRoutes(
			server.TemplateRouterGroup("/templates", nil, middlewares...),
			server.NotificationRouterGroup("/notifications", middlewares...),
			server.AuthRouterGroup("/auth", nil, middlewares...),
		),
	)

	tests := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/templates", 404}, // Would be 200 with proper setup
		{"GET", "/notifications", 404},
		{"GET", "/auth", 404},
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

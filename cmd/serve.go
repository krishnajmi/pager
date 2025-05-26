package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kp/pager/databases/sql"
	"github.com/kp/pager/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "apis",
	Short: "Start the Pager api server",
	Long:  `This is starting point for apis`,
	Run: func(cmd *cobra.Command, args []string) {
		servicePrefix := "/pager/v1"
		templatePrefix := servicePrefix + "/template"
		notificationPrefix := servicePrefix + "/notification"
		loginPrefix := servicePrefix + "/user"
		middlewares := []gin.HandlerFunc{
			server.AuthPermissionMiddleware(sql.PagerOrm),
			server.RecoveryMiddleware(),
		}
		brokers := []string{"localhost:9092"} // Replace with your Kafka broker addresses
		router := server.InitServer(middlewares, server.WithTimeOut(0*time.Second),
			server.CreateRoutes(
				server.TemplateRouterGroup(templatePrefix, sql.PagerOrm, middlewares...),
				server.NotificationRouterGroup(notificationPrefix, brokers, middlewares...),
				server.AuthRouterGroup(loginPrefix, sql.PagerOrm, middlewares...),
			),
		)

		server := http.Server{
			Addr:    ":" + "8000",
			Handler: router,
		}
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(err)
			}
		}()
		fmt.Println("httpServerListeningAndServingOn ", server.Addr)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("ServerShutdown:", err)
		}
		log.Println("ExitingServer...")
	},
}

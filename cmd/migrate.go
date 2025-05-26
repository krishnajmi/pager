package cmd

import (
	"log/slog"

	comm_models "github.com/kp/pager/communicator/models"
	"github.com/kp/pager/databases/sql"
	login_models "github.com/kp/pager/login/models"
	notification_models "github.com/kp/pager/notification/models"
	"github.com/kp/pager/templates"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run all database migrations for the application`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("Running database migrations...")
		dbMigrate()
		slog.Info("Migrations completed successfully")
	},
}

func dbMigrate() {
	// Notification system tables
	sql.PagerOrm.AutoMigrate(&notification_models.NotificationSession{})
	// Template system tables
	sql.PagerOrm.AutoMigrate(&templates.NotificationTemplate{})
	// Communication system tables
	sql.PagerOrm.AutoMigrate(&comm_models.CommunicationLogs{})
	// Auth system tables
	sql.PagerOrm.AutoMigrate(&login_models.User{})
	sql.PagerOrm.AutoMigrate(&login_models.Permission{})
	sql.PagerOrm.AutoMigrate(&login_models.UserPermission{})
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

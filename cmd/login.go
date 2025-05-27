package cmd

import (
	"log/slog"
	"os"

	"github.com/kp/pager/common"
	"github.com/kp/pager/databases/sql"
	login_models "github.com/kp/pager/login/models"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Long:  `Register a new user with username and password`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		encryptedPass := common.Encryptbase64(password)

		userType, _ := cmd.Flags().GetString("usertype")
		user := &login_models.User{
			Username: username,
			Password: encryptedPass,
			UserType: userType,
		}

		result := sql.PagerOrm.Create(user)
		if result.Error != nil {
			slog.Error("Failed to register user",
				"error", result.Error,
				"username", username,
			)
			os.Exit(1)
		}

		slog.Info("Successfully registered user",
			"username", username,
			"userID", user.ID,
		)
	},
}

func init() {
	registerCmd.Flags().StringP("username", "u", "", "Username for the new account")
	registerCmd.Flags().StringP("password", "p", "", "Password for the new account")
	registerCmd.Flags().StringP("usertype", "t", "user", "User type (admin/user)")
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(migrateCmd)
}

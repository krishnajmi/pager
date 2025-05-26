package cmd

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/kp/pager/consumers"
	aws_db "github.com/kp/pager/databases/aws"
	"github.com/kp/pager/databases/kafka"
	"github.com/spf13/cobra"
)

var (
	appConfig  *AppConfig
	envVarList = []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_REGION",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_USER",
		"DB_PASSWORD",
		"REDIS_HOST",
		"REDIS_PORT",
		"REDIS_PASSWORD",
		"REDIS_DB",
		"KAFKA_BROKERS",
		"KAFKA_TOPIC",
		"KAFKA_USERNAME",
		"KAFKA_PASSWORD",
	}
	rootCmd = &cobra.Command{
		Use:   "pager-cli",
		Short: "All applications for pager",
		Long: `pager collection command
				etc`,
	}
)

func Run() error {
	return rootCmd.Execute()
}

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	cobra.OnInitialize(setupDeps)
}

func setupDeps() {
	ctx := context.Background()
	appConfig = getAppConfig(ctx)
	slog.Info("appConfig", "config", appConfig)

	// Initialize database
	_, _, err := appConfig.DatabaseConfigType.InitDatabase()
	if err != nil {
		slog.Error("errorConnectingToDatabase",
			slog.String("error", err.Error()),
			slog.String("database", appConfig.DatabaseConfigType.Database),
		)
		os.Exit(1)
	}

	// Initialize Kafka
	brokers := strings.Split(appConfig.KafkaConfig.Brokers, ",")
	err = kafka.RunMigrations(brokers) // Add this migration in CLI
	if err != nil {
		slog.Error("errorInitializingKafka",
			slog.String("error", err.Error()),
			slog.String("brokers", appConfig.KafkaConfig.Brokers),
		)
		os.Exit(1)
	}
	// Start batch consumer
	go consumers.StartBatchConsumer(brokers)
	dbMigrate()
}

func getAppConfig(ctx context.Context) *AppConfig {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	var appSecretConfig = &AppConfig{}
	cred := credentials.NewStaticCredentials(accessKey, secretKey, region)
	awsConfig := &aws.Config{
		Region:      aws.String(region),
		Credentials: cred,
	}
	if err := ReadAppConfig(envVarList, appSecretConfig, awsConfig); err != nil {
		panic(err)
	}

	return appSecretConfig
}

func ReadAppConfig(envList []string, destination any, configs ...*aws.Config) error {
	var conf map[string]string = make(map[string]string)
	if !IsEnvLocal() {
		secrets, err := aws_db.GetSecret(configs...)
		if err != nil {
			slog.Error("ErrReadingAWSSecrets", "err_msg", err)
			return err
		}
		json.Unmarshal([]byte(secrets), &conf)
	}

	for _, key := range envList {
		if _, ok := conf[key]; !ok {
			conf[key] = os.Getenv(key)
		}
	}

	mapBytes, _ := json.Marshal(&conf)
	if err := json.Unmarshal(mapBytes, &destination); err != nil {
		slog.Error("ErrReadingAWSSecrets", "err_msg", err)
		return err
	}

	return nil
}

func IsEnvLocal() bool {
	if os.Getenv("APP_ENV") == "" {
		return true
	}
	if strings.EqualFold(os.Getenv("APP_ENV"), "local") {
		return true
	}
	return false
}

func IsEnvProd() bool {
	if strings.EqualFold(os.Getenv("APP_ENV"), "production") {
		return true
	}
	if strings.EqualFold(os.Getenv("APP_ENV"), "prod") {
		return true
	}
	return false
}

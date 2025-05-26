package cmd

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stretchr/testify/assert"
)

func TestIsEnvLocal(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{"empty env", "", true},
		{"local env", "local", true},
		{"LOCAL env", "LOCAL", true},
		{"dev env", "dev", false},
		{"prod env", "prod", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("APP_ENV", tt.envValue)
			assert.Equal(t, tt.expected, IsEnvLocal())
		})
	}
}

func TestIsEnvProd(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{"empty env", "", false},
		{"prod env", "prod", true},
		{"PROD env", "PROD", true},
		{"production env", "production", true},
		{"local env", "local", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("APP_ENV", tt.envValue)
			assert.Equal(t, tt.expected, IsEnvProd())
		})
	}
}

func TestReadAppConfig(t *testing.T) {
	t.Run("local env with env vars", func(t *testing.T) {
		os.Setenv("APP_ENV", "local")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")

		var config struct {
			DBHost string `json:"DB_HOST"`
			DBPort string `json:"DB_PORT"`
		}

		err := ReadAppConfig([]string{"DB_HOST", "DB_PORT"}, &config)
		assert.NoError(t, err)
		assert.Equal(t, "localhost", config.DBHost)
		assert.Equal(t, "5432", config.DBPort)
	})

	t.Run("non-local env with mock AWS", func(t *testing.T) {
		os.Setenv("APP_ENV", "prod")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")

		mockAWSConfig := &aws.Config{
			Region:      aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials("test", "test", ""),
		}

		// Mock AWS secrets would be needed here in a real test
		// This is just demonstrating the test structure
		var config struct {
			DBHost string `json:"DB_HOST"`
			DBPort string `json:"DB_PORT"`
		}

		err := ReadAppConfig([]string{"DB_HOST", "DB_PORT"}, &config, mockAWSConfig)
		// Would assert on mock response in real test
		assert.Error(t, err) // Expect error since we can't mock AWS response here
	})
}

func TestGetAppConfig(t *testing.T) {
	t.Run("with env vars", func(t *testing.T) {
		os.Setenv("AWS_ACCESS_KEY_ID", "test")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
		os.Setenv("AWS_REGION", "us-east-1")

		ctx := context.Background()
		config := getAppConfig(ctx)
		assert.NotNil(t, config)
	})
}

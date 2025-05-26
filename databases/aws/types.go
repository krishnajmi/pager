package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type AwsServicesProviderer interface {
	IntializeAwsProvider(ctx context.Context)
	GetSecretKeyValue(ctx context.Context, secretKey string) (*secretsmanager.GetSecretValueOutput, error)
	GetSecretManager(ctx context.Context) *secretsmanager.SecretsManager
}

type AwsServiceProvider struct {
	secretManager *secretsmanager.SecretsManager
	Region        string
}

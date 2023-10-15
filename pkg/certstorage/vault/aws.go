package vault

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/vault/api"
)

type AwsDynamicCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
	Expiry          time.Time
}

func mapVaultAwsCredentialResponse(secret *api.Secret) (aws.Credentials, error) {
	if secret == nil || secret.Data == nil {
		return aws.Credentials{}, errors.New("empty secret / payload")
	}

	accessKey, ok := secret.Data["access_key"].(string)
	if !ok {
		return aws.Credentials{}, errors.New("empty 'access_key'")
	}

	secretKey, ok := secret.Data["secret_key"].(string)
	if !ok {
		return aws.Credentials{}, errors.New("empty 'secret_key'")
	}

	return aws.Credentials{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		CanExpire:       true,
		Expires:         time.Now().Add(time.Duration(secret.LeaseDuration) * time.Second),
		Source:          "vault",
	}, nil
}

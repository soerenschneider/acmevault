package vault

import (
	"errors"
	"time"

	"github.com/hashicorp/vault/api"
)

type AwsDynamicCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
	Expiry          time.Time
}

func mapVaultAwsCredentialResponse(secret *api.Secret) (*AwsDynamicCredentials, error) {
	if secret == nil || secret.Data == nil {
		return nil, errors.New("empty secret / payload")
	}

	accessKey, ok := secret.Data["access_key"].(string)
	if !ok {
		return nil, errors.New("empty 'access_key'")
	}

	secretKey, ok := secret.Data["secret_key"].(string)
	if !ok {
		return nil, errors.New("empty 'secret_key'")
	}

	return &AwsDynamicCredentials{
		AccessKeyId:     accessKey,
		SecretAccessKey: secretKey,
		Expiry:          time.Now().Add(time.Duration(secret.LeaseDuration) * time.Second),
	}, nil
}

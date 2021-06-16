package vault

import (
	"errors"
	"fmt"
	"github.com/hashicorp/vault/api"
	"time"
)

const awsRole = "acmevault"

type AwsDynamicCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
	Expiry          time.Time
}

func mapVaultAwsCredentialResponse(secret api.Secret) (*AwsDynamicCredentials, error) {
	if secret.Data == nil {
		return nil, errors.New("no 'Data' payload available")
	}

	accessKey := fmt.Sprintf("%s", secret.Data["access_key"])
	secretKey := fmt.Sprintf("%s", secret.Data["secret_key"])
	if len(accessKey) == 0 || len(secretKey) == 0 {
		return nil, errors.New("missing access- and/or secret-key")
	}

	return &AwsDynamicCredentials{
		AccessKeyId:     accessKey,
		SecretAccessKey: secretKey,
		Expiry:          time.Now().Add(time.Duration(secret.LeaseDuration) * time.Second),
	}, nil
}
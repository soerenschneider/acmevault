package acme

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/go-acme/lego/v4/challenge"
	legoRoute53 "github.com/go-acme/lego/v4/providers/dns/route53"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/pkg/certstorage/vault"
)

const AwsIamPropagationImpediment = 20 * time.Second

type DynamicCredentialsProvider struct {
	vault  *vault.VaultBackend
	expiry time.Time
}

func NewDynamicCredentialsProvider(vault *vault.VaultBackend) (aws.CredentialsProvider, error) {
	if nil == vault {
		return nil, errors.New("no vault backend provided")
	}

	bla := &DynamicCredentialsProvider{vault: vault}
	return aws.NewCredentialsCache(bla), nil
}

func (m *DynamicCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	log.Info().Msg("Trying to read AWS credentials from Vault")
	dynamicCredentials, err := m.vault.ReadAwsCredentials()
	if err != nil {
		return aws.Credentials{}, fmt.Errorf("could not login at vault: %v", err)
	}

	m.expiry = dynamicCredentials.Expiry
	cred := ConvertCredentials(*dynamicCredentials)

	log.Info().Msgf("Received AWS credentials with access id %s, waiting for %v for eventual consistency", cred.AccessKeyID, AwsIamPropagationImpediment)

	// The credentials we receive are usually not effective at AWS, yet, so we need to wait for a bit until
	// the changes on AWS are propagated
	time.Sleep(AwsIamPropagationImpediment)
	return cred, nil
}

func ConvertCredentials(dynamicCredentials vault.AwsDynamicCredentials) aws.Credentials {
	return aws.Credentials{
		AccessKeyID:     dynamicCredentials.AccessKeyId,
		SecretAccessKey: dynamicCredentials.SecretAccessKey,
		CanExpire:       true,
		Expires:         dynamicCredentials.Expiry,
		Source:          "vault",
	}
}

func (m *DynamicCredentialsProvider) IsExpired() bool {
	return time.Now().After(m.expiry)
}

func BuildRoute53DnsProvider(credProvider ...aws.CredentialsProvider) (challenge.Provider, error) {
	awsConf, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	if nil == credProvider || len(credProvider) == 0 {
		log.Info().Msg("Trying to use static credentials to build route53 session")
	} else {
		log.Info().Msg("Passing dynamic credentials provider to build route53 session")
		awsConf.Credentials = credProvider[0]
	}

	client := route53.NewFromConfig(awsConf)
	legoConf := legoRoute53.NewDefaultConfig()
	legoConf.Client = client
	return legoRoute53.NewDNSProviderConfig(legoConf)
}

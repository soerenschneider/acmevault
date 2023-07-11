package acme

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsRoute53 "github.com/aws/aws-sdk-go/service/route53"
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

func NewDynamicCredentialsProvider(vault *vault.VaultBackend) (*DynamicCredentialsProvider, error) {
	if nil == vault {
		return nil, errors.New("no vault backend provided")
	}

	return &DynamicCredentialsProvider{vault: vault}, nil
}

func (m *DynamicCredentialsProvider) Retrieve() (credentials.Value, error) {
	log.Info().Msg("Trying to read AWS credentials from Vault")
	dynamicCredentials, err := m.vault.ReadAwsCredentials()
	if err != nil {
		return credentials.Value{}, fmt.Errorf("could not login at vault: %v", err)
	}

	m.expiry = dynamicCredentials.Expiry
	cred := ConvertCredentials(*dynamicCredentials)

	log.Info().Msgf("Received AWS credentials with access id %s, waiting for %v for eventual consistency", cred.AccessKeyID, AwsIamPropagationImpediment)

	// The credentials we receive are usually not effective at AWS, yet, so we need to wait for a bit until
	// the changes on AWS are propagated
	time.Sleep(AwsIamPropagationImpediment)
	return cred, nil
}

func ConvertCredentials(dynamicCredentials vault.AwsDynamicCredentials) credentials.Value {
	return credentials.Value{
		AccessKeyID:     dynamicCredentials.AccessKeyId,
		SecretAccessKey: dynamicCredentials.SecretAccessKey,
		ProviderName:    "vault",
	}
}

func (m *DynamicCredentialsProvider) IsExpired() bool {
	return time.Now().After(m.expiry)
}

func BuildRoute53DnsProvider(credProvider ...DynamicCredentialsProvider) (challenge.Provider, error) {
	var awsSession *session.Session
	var err error
	if nil == credProvider || len(credProvider) == 0 {
		log.Info().Msg("Trying to use static credentials to build route53 session")
		awsSession, err = session.NewSession()
		if err != nil {
			return nil, err
		}
	} else {
		log.Info().Msg("Passing dynamic credentials provider to build route53 session")
		awsSession, err = session.NewSession(&aws.Config{
			Credentials: credentials.NewCredentials(&credProvider[0]),
		})
		if err != nil {
			return nil, err
		}
	}

	svc := awsRoute53.New(awsSession)
	conf := legoRoute53.NewDefaultConfig()
	conf.Client = svc
	return legoRoute53.NewDNSProviderConfig(conf)
}

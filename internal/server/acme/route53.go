package acme

import (
	vault2 "acmevault/pkg/certstorage/vault"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsRoute53 "github.com/aws/aws-sdk-go/service/route53"
	"github.com/go-acme/lego/v4/challenge"
	legoRoute53 "github.com/go-acme/lego/v4/providers/dns/route53"
	"github.com/rs/zerolog/log"
	"time"
)

const AwsIamPropagationImpediment = 20 * time.Second

type DynamicCredentialsProvider struct {
	vault  *vault2.VaultBackend
	expiry time.Time
}

func NewDynamicCredentialsProvider(vault *vault2.VaultBackend) (*DynamicCredentialsProvider, error) {
	if nil == vault {
		return nil, errors.New("no vault backend provided")
	}
	
	return &DynamicCredentialsProvider{vault: vault}, nil
}

func (m *DynamicCredentialsProvider) Retrieve() (credentials.Value, error) {
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

func ConvertCredentials(dynamicCredentials vault2.AwsDynamicCredentials) credentials.Value {
	return credentials.Value{
		AccessKeyID:     dynamicCredentials.AccessKeyId,
		SecretAccessKey: dynamicCredentials.SecretAccessKey,
		ProviderName:    "vault",
	}
}

func (m *DynamicCredentialsProvider) IsExpired() bool {
	return time.Now().Before(m.expiry)
}

func BuildRoute53DnsProvider(credProvider ...DynamicCredentialsProvider) (challenge.Provider, error) {
	var awsSession *session.Session
	if nil != credProvider && len(credProvider) > 0 {
		awsSession = session.Must(session.NewSession())
	} else {
		awsSession = session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewCredentials(&credProvider[0]),
		}))
	}

	svc := awsRoute53.New(awsSession)
	conf := legoRoute53.NewDefaultConfig()
	conf.Client = svc
	return legoRoute53.NewDNSProviderConfig(conf)
}
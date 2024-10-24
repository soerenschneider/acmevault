package vault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-acme/lego/v4/acme"
	"github.com/go-acme/lego/v4/registration"
	"github.com/hashicorp/vault/api"
	vault "github.com/hashicorp/vault/api"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/internal/metrics"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
)

const (
	timeout        = 10 * time.Second
	backoffRetries = 5
)

type VaultBackend struct {
	client   *api.Client
	conf     config.VaultConfig
	basePath string
}

func NewVaultBackend(vaultClient *api.Client, vaultConfig config.VaultConfig) (*VaultBackend, error) {
	vault := &VaultBackend{
		client:   vaultClient,
		conf:     vaultConfig,
		basePath: fmt.Sprintf("acmevault/%s", vaultConfig.PathPrefix),
		//basePath: fmt.Sprintf("%s/data/%s", vaultConfig.Kv2MountPath, vaultConfig.PathPrefix),
	}

	return vault, nil
}

func (vault *VaultBackend) WriteCertificate(resource *certstorage.AcmeCertificate) error {
	// save private key
	privateKey := resource.PrivateKey
	resource.PrivateKey = nil

	data := certstorage.CertToMap(resource)
	certPath := vault.getCertDataPath(resource.Domain)
	err := vault.writeKv2Secret(certPath, data)
	if err != nil {
		return fmt.Errorf("could not write certificate data for %s: %v", resource.Domain, err)
	}

	data = map[string]interface{}{
		"private_key": privateKey,
	}
	secretPath := vault.getSecretDataPath(resource.Domain)
	err = vault.writeKv2Secret(secretPath, data)
	if err != nil {
		return fmt.Errorf("could not write secrete data for domain %s: %v", resource.Domain, err)
	}

	return nil
}

func (vault *VaultBackend) ReadPublicCertificateData(domain string) (*certstorage.AcmeCertificate, error) {
	certPath := vault.getCertDataPath(domain)
	data, err := vault.readKv2Secret(certPath)
	if err != nil {
		return nil, fmt.Errorf("could not readKv2Secret public cert data from vault for domain %s: %v", domain, err)
	}
	return certstorage.MapToCert(data)
}

func (vault *VaultBackend) ReadFullCertificateData(domain string) (*certstorage.AcmeCertificate, error) {
	cert, err := vault.ReadPublicCertificateData(domain)
	if err != nil {
		return nil, err
	}

	privateKeyPath := vault.getSecretDataPath(domain)
	data, err := vault.readKv2Secret(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not readKv2Secret private data from vault for domain %s: %v", domain, err)
	}

	_, ok := data["private_key"]
	if !ok {
		return nil, fmt.Errorf("successfully readKv2Secret secret from vault but no private key data avaialble for domain: %s", domain)
	}

	privRaw := fmt.Sprintf("%s", data["private_key"])
	priv, err := base64.StdEncoding.DecodeString(privRaw)
	if err != nil {
		return nil, fmt.Errorf("can not decode private key: %v", err)
	}
	cert.PrivateKey = priv

	return cert, err
}

func (vault *VaultBackend) WriteAccount(acmeRegistration certstorage.AcmeAccount) error {
	jsonBytes, err := json.MarshalIndent(acmeRegistration.Registration.Body, "", "\t")
	if err != nil {
		return err
	}

	key, _ := certstorage.ConvertToPem(acmeRegistration.Key)

	data := map[string]interface{}{
		certstorage.VaultAccountKeyUri:     acmeRegistration.Registration.URI,
		certstorage.VaultAccountKeyAccount: jsonBytes,
		certstorage.VaultAccountKeyKey:     key,
		certstorage.VaultAccountKeyEmail:   acmeRegistration.Email,
	}

	accountPath := vault.getAccountPath(acmeRegistration.Email)

	err = vault.writeKv2Secret(accountPath, data)
	return err
}

func (vault *VaultBackend) ReadAccount(hash string) (*certstorage.AcmeAccount, error) {
	accountPath := vault.getAccountPath(hash)
	data, err := vault.readKv2Secret(accountPath)
	if err != nil {
		return nil, fmt.Errorf("could not readKv2Secret account from vault: %w", translateError(err))
	}

	var account acme.Account
	accountData := fmt.Sprintf("%v", data[certstorage.VaultAccountKeyAccount])
	accountJson, err := base64.StdEncoding.DecodeString(accountData)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode string: %v", err)
	}

	if err = json.Unmarshal(accountJson, &account); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal json: %v", err)
	}

	registration := &registration.Resource{
		URI:  fmt.Sprintf("%s", data[certstorage.VaultAccountKeyUri]),
		Body: account,
	}

	decodedKey, err := certstorage.FromPem([]byte(fmt.Sprintf("%v", data[certstorage.VaultAccountKeyKey])))
	if err != nil {
		return nil, fmt.Errorf("can not decode pem private key: %vault", err)
	}
	conf := certstorage.AcmeAccount{
		Email:        fmt.Sprintf("%v", data[certstorage.VaultAccountKeyEmail]),
		Key:          decodedKey,
		Registration: registration,
	}

	return &conf, nil
}

func (vault *VaultBackend) Authenticate() error {
	return nil
}

func (vault *VaultBackend) Logout() error {
	return vault.client.Auth().Token().RevokeSelf("xxx")
}

func (vault *VaultBackend) ReadAwsCredentials() (aws.Credentials, error) {
	metrics.AwsDynCredentialsRequested.Inc()
	path := vault.getAwsCredentialsPath()
	secret, err := vault.client.Logical().Read(path)
	if err != nil {
		metrics.AwsDynCredentialsRequestErrors.Inc()
		return aws.Credentials{}, fmt.Errorf("could not gather dynamic credentials: %v", err)
	}

	return mapVaultAwsCredentialResponse(secret)
}

func (vault *VaultBackend) writeKv2Secret(secretPath string, data map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := vault.client.KVv2(vault.conf.Kv2MountPath).Put(ctx, secretPath, data)
	return err
}

func (vault *VaultBackend) readKv2Secret(path string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	secret, err := vault.client.KVv2(vault.conf.Kv2MountPath).Get(ctx, path)
	if err != nil {
		return nil, translateError(err)
	}

	if secret == nil {
		return nil, certstorage.ErrNotFound
	}

	return secret.Data, nil
}

func translateError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, vault.ErrSecretNotFound) {
		return certstorage.ErrNotFound
	}

	vaultErr, ok := err.(*vault.ResponseError)
	if !ok {
		return err
	}

	if vaultErr.StatusCode == 404 {
		return certstorage.ErrNotFound
	}

	if vaultErr.StatusCode == 403 {
		return certstorage.ErrPermissionDenied
	}

	return err
}

func (vault *VaultBackend) formatDomain(domain string) string {
	if len(vault.conf.DomainPathFormat) == 0 {
		return domain
	}
	return fmt.Sprintf(vault.conf.DomainPathFormat, domain)
}

func (vault *VaultBackend) getAwsCredentialsPath() string {
	return fmt.Sprintf("%s/creds/%s", vault.conf.AwsMountPath, vault.conf.AwsRole)
}

func (vault *VaultBackend) getAccountPath(hash string) string {
	return fmt.Sprintf("%s/server/account/%s", vault.basePath, hash)
}

func (vault *VaultBackend) getCertDataPath(domain string) string {
	domainFormatted := vault.formatDomain(domain)
	return fmt.Sprintf("%s/client/%s/certificate", vault.basePath, domainFormatted)
}

func (vault *VaultBackend) getSecretDataPath(domain string) string {
	domainFormatted := vault.formatDomain(domain)
	return fmt.Sprintf("%s/client/%s/privatekey", vault.basePath, domainFormatted)
}

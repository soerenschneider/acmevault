package vault

import (
	"acmevault/internal"
	"acmevault/internal/config"
	"acmevault/pkg/certstorage"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/acme"
	"github.com/go-acme/lego/v4/registration"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	timeout        = 10 * time.Second
	backoffRetries = 5

	vaultAcmeApproleLoginPath = "/auth/approle/login"
	vaultAcmePathPrefix       = "/secret/acme"
)

type VaultBackend struct {
	client       *api.Client
	conf         config.VaultConfig
	tokenStorage TokenStorage
	revokeToken  bool
}

func NewVaultBackend(vaultConfig config.VaultConfig, storage TokenStorage) (*VaultBackend, error) {
	config := &api.Config{
		Timeout:    timeout,
		MaxRetries: backoffRetries,
		Backoff:    retryablehttp.DefaultBackoff,
		Address:    vaultConfig.VaultAddr,
	}

	var err error
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("couldn't build client: %v", err)
	}

	// set initial token, can be empty as well, ignore potential errors
	initialToken, _ := storage.ReadToken()
	client.SetToken(initialToken)

	vault := &VaultBackend{
		client:       client,
		conf:         vaultConfig,
		tokenStorage: storage,
	}

	_, err = vault.authenticate()
	if err != nil {
		log.Fatal().Msgf("All authentication options exhausted, giving up: %v", err)
	}

	if vault.conf.IsTokenIncreaseEnabled() {
		go func() {
			ticker := time.NewTicker(time.Duration(vault.conf.TokenIncreaseInterval) * time.Second)
			for {
				select {
				case <-ticker.C:
					vault.RenewToken(vault.conf.TokenIncreaseSeconds)
				}
			}
		}()
	}

	return vault, nil
}

func (vault *VaultBackend) WriteCertificate(resource *certstorage.AcmeCertificate) error {
	location := getCertPath(resource.Domain)
	data := certstorage.CertToMap(resource)
	_, err := vault.client.Logical().Write(location, data)
	return err
}

func (vault *VaultBackend) ReadCertificate(domain string) (*certstorage.AcmeCertificate, error) {
	location := getCertPath(domain)
	log.Info().Msgf("Trying to read secret from %s", location)
	secret, err := vault.client.Logical().Read(location)
	if err != nil {
		return nil, fmt.Errorf("could not read secret %s: %v", location, err)
	}

	if secret == nil {
		return nil, fmt.Errorf("no cert found for domain %s", domain)
	}

	return certstorage.MapToCert(secret.Data)
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

	path := getAccountPath(acmeRegistration.Email)
	log.Info().Msgf("Trying to write acme account %s", acmeRegistration.Email)
	_, err = vault.client.Logical().Write(path, data)
	return err
}

func (vault *VaultBackend) ReadAccount(hash string) (*certstorage.AcmeAccount, error) {
	secret, err := vault.client.Logical().Read(getAccountPath(hash))
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, errors.New("entry does not exist, yet")
	}

	var account acme.Account
	accountData := fmt.Sprintf("%v", secret.Data[certstorage.VaultAccountKeyAccount])
	accountJson, err := base64.StdEncoding.DecodeString(accountData)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode string: %v", err)
	}
	err = json.Unmarshal(accountJson, &account)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal json: %v", err)
	}

	registration := &registration.Resource{
		URI:  fmt.Sprintf("%s", secret.Data[certstorage.VaultAccountKeyUri]),
		Body: account,
	}

	decodedKey, err := certstorage.FromPem([]byte(fmt.Sprintf("%vault", secret.Data[certstorage.VaultAccountKeyKey])))
	if err != nil {
		return nil, fmt.Errorf("can not decode pem private key: %vault", err)
	}
	conf := certstorage.AcmeAccount{
		Email:        fmt.Sprintf("%vault", secret.Data[certstorage.VaultAccountKeyEmail]),
		Key:          decodedKey,
		Registration: registration,
	}

	return &conf, nil
}

// Renew lookups the currently used token and tries to renew it by a given TTL if it's renewable.
// Returns true if the token was successfully renewed, otherwise false.
func (vault *VaultBackend) RenewToken(tokenIncrement int) (bool, error) {
	log.Info().Msgf("Trying to renew token by %d seconds", tokenIncrement)
	currentToken, err := vault.authenticate()
	if err != nil {
		return false, err
	}

	if currentToken.Renewable {
		secret, err := vault.client.Auth().Token().RenewSelf(tokenIncrement)
		if err == nil && secret != nil {
			ttl, _ := secret.TokenTTL()
			log.Info().Msgf("Renewed token valid until %v", ttl)
			return true, nil
		}

		log.Info().Msgf("Could not renew token: %v", err)
		return false, err
	}

	log.Info().Msg("Token is not renewable")
	return false, nil
}

// authenticate tries to authenticate against vault using all possible configured options. Upon authentication
// a lookup on the token is performed to verify it. The resulting token is returned. If no authentication is possible,
// nil and an appropriate error is returned.
func (vault *VaultBackend) authenticate() (*TokenData, error) {
	log.Info().Msg("Trying authentication against vault")
	// try to lookup ourself, maybe we're already authenticated
	tokenData, err := vault.lookupToken()
	if err == nil && tokenData != nil {
		log.Info().Msgf("Already successfully authenticated against vault, token valid until %s", tokenData.PrettyExpiryDate())
		return tokenData, nil
	}

	// that didn't work, try token from storage
	log.Info().Msgf("Trying token from storage backend")
	token, err := vault.tokenStorage.ReadToken()
	if err == nil && len(token) > 0 {
		log.Info().Msg("ReadCertificate token from storage, testing token validity")
		tokenData, err := vault.confirmToken(token)
		if err == nil && tokenData != nil {
			log.Info().Msgf("Successfully authenticated against vault, token valid until %s", tokenData.PrettyExpiryDate())
			return tokenData, nil
		}
	}
	log.Info().Msgf("Trying to login via AppRole using roleId %s", vault.conf.RoleId)

	token, err = vault.loginAppRole(vault.conf.RoleId, vault.conf.SecretId)
	if err != nil {
		return nil, fmt.Errorf("could not login via AppRole %s: %v", vault.conf.RoleId, err)
	}

	tokenData, err = vault.confirmToken(token)
	if tokenData == nil || err != nil {
		return nil, errors.New("acquired token doesn't work, giving up")
	}

	log.Info().Msgf("Successfully authenticated via AppRole %s, token valid until %s", vault.conf.RoleId, tokenData.PrettyExpiryDate())
	log.Debug().Msg("Storing newly acquired token")
	err = vault.tokenStorage.StoreToken(token)
	if err != nil {
		log.Info().Msgf("could not store token: %v\n", err)
	}

	return tokenData, nil
}

// confirmToken sets the appropriate token for the vault client and performs a lookup on itself.
// Returns TokenData if the authentication was successful, otherwise nil and an error.
func (vault *VaultBackend) confirmToken(token string) (*TokenData, error) {
	vault.client.SetToken(token)
	tokenData, err := vault.lookupToken()

	// Update token lifetime metric
	if err == nil && tokenData != nil && !tokenData.ExpireTime.IsZero() {
		internal.VaultTokenExpiryTimestamp.Set(float64(tokenData.ExpireTime.Unix()))
	}
	return tokenData, err
}

// lookupToken looks up the currently set token and upon returns the TokenData that's associated to it. If the
// token is invalid, nil and an error is returned.
func (vault *VaultBackend) lookupToken() (*TokenData, error) {
	secret, err := vault.client.Auth().Token().LookupSelf()
	if err != nil || secret == nil {
		return nil, err
	}
	return FromSecret(secret), nil
}

// loginAppRole performs a login against vault using the "App Role" mechanism. Returns a vault token upon successful
// login, otherwise an empty string and an appropriate error.
func (vault *VaultBackend) loginAppRole(roleId, secretId string) (string, error) {
	data := map[string]interface{}{
		"role_id":   roleId,
		"secret_id": secretId,
	}

	resp, err := vault.client.Logical().Write(vaultAcmeApproleLoginPath, data)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("received empty reply")
	}

	return resp.Auth.ClientToken, nil
}

func (vault *VaultBackend) Cleanup() {
	if !vault.revokeToken {
		return
	}
	log.Info().Msg("Performing revokeToken, trying to revoke token...")
	err := vault.client.Auth().Token().RevokeSelf("")
	if err != nil {
		log.Info().Msgf("Error while revoking token: %v", err)
	}
}

func (vault *VaultBackend) ReadAwsCredentials() (*AwsDynamicCredentials, error) {
	internal.AwsDynCredentialsRequested.Inc()
	path := fmt.Sprintf("/aws/creds/%s", awsRole)
	secret, err := vault.client.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("could not gather dynamic credentials: %v", err)
	}
	return mapVaultAwsCredentialResponse(*secret)
}

func getAccountPath(hash string) string {
	return fmt.Sprintf("%s/server/account/%s", vaultAcmePathPrefix, hash)
}

func getCertPath(domain string) string {
	return fmt.Sprintf("%s/client/%s", vaultAcmePathPrefix, domain)
}

package vault

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/acme"
	"github.com/go-acme/lego/v4/registration"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
	"net/url"
	"path"
	"time"
)

const (
	timeout        = 10 * time.Second
	backoffRetries = 5

	vaultAcmeApproleLoginPath = "/auth/approle/login"
	vaultSecretPathMount = "/secret/data"
	maxVersions = 1
)

type VaultBackend struct {
	client           *api.Client
	conf             config.VaultConfig
	tokenStorage     TokenStorage
	revokeToken      bool
	namespacedPrefix string
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
		client:           client,
		conf:             vaultConfig,
		tokenStorage:     storage,
		namespacedPrefix: fmt.Sprintf("%s/%s", vaultSecretPathMount, vaultConfig.PathPrefix),
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
	certPath := vault.getCertDataPath(resource.Domain)
	data := certstorage.CertToMap(resource)
	err := vault.writeSecretV2(certPath, data)
	if err != nil {
		return fmt.Errorf("could not write certificate data for %s: %v", resource.Domain, err)
	}
	return err
}

func (vault *VaultBackend) writeSecretV1(secretPath string, data map[string]interface{}) error {
	secret, err := vault.client.Logical().Write(secretPath, data)
	printWarning("Received warnings while writing secretV1", secret)
	return err
}

func printWarning(msg string, secret *api.Secret) {
	if len(secret.Warnings) > 0 {
		var warningMsg string
		for _, warn := range secret.Warnings {
			warningMsg += warn
			warningMsg += " / "
		}
		log.Warn().Msgf("%s: %s", msg, warningMsg)
	}
}

func (vault *VaultBackend) writeSecretV2(secretPath string, data map[string]interface{}) error {
	vaultUrl, err := url.Parse(vault.conf.VaultAddr)
	if err != nil {
		return err
	}
	vaultUrl.Path = path.Join(vaultUrl.Path, "/v1" + secretPath)

	payload, err := wrapPayload(data)
	if err != nil {
		return err
	}
	req := &api.Request{
		Method:         "POST",
		URL:            vaultUrl,
		ClientToken:    vault.client.Token(),
		BodyBytes: payload,
	}

	_, err = vault.client.RawRequest(req)
	return err
}

// wrapPayload wraps a map containing the payload into a another map, all contained within the
// `data` field to use the KV2 API of vault. Returns the data as json-encoded []byte slice.
func wrapPayload(data map[string]interface{}) ([]byte, error) {
	if data == nil {
		data = map[string]interface{}{}
	}
	bla := struct {
		Data map[string]interface{} `json:"data"`
	}{
		Data: data,
	}

	return json.Marshal(bla)
}

func (vault *VaultBackend) ReadCertificate(domain string) (*certstorage.AcmeCertificate, error) {
	certPath := vault.getCertDataPath(domain)
	log.Info().Msgf("Trying to read secret from %s", certPath)
	return vault.read(certPath)
}

func (vault *VaultBackend) read(certPath string) (*certstorage.AcmeCertificate, error) {
	secret, err := vault.client.Logical().Read(certPath)
	if err != nil {
		return nil, fmt.Errorf("could not read secret %s: %v", certPath, err)
	}

	if secret == nil {
		return nil, fmt.Errorf("nothing found at %s", certPath)
	}

	var data map[string]interface{}
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("could not parse map")
	}
	return certstorage.MapToCert(data)
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
	log.Info().Msgf("Trying to write acme account %s", acmeRegistration.Email)

	err = vault.writeSecretV2(accountPath, data)
	return err
}

func (vault *VaultBackend) readPathV1(path string) (map[string]interface{}, error) {
	secret, err := vault.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, errors.New("entry does not exist, yet")
	}

	return secret.Data, nil
}

func (vault *VaultBackend) readPathV2(path string) (map[string]interface{}, error) {
	secret, err := vault.readPathV1(path)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	_, ok := secret["data"]
	if !ok {
		return nil, errors.New("no data portion available")
	}
	data, ok = secret["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("could not convert map")
	}

	return data, nil
}

func (vault *VaultBackend) ReadAccount(hash string) (*certstorage.AcmeAccount, error) {
	accountPath := vault.getAccountPath(hash)
	data, err := vault.readPathV2(accountPath)
	if err != nil {
		return nil, fmt.Errorf("could not read account from vault: %v", err)
	}

	var account acme.Account
	accountData := fmt.Sprintf("%v", data[certstorage.VaultAccountKeyAccount])
	accountJson, err := base64.StdEncoding.DecodeString(accountData)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode string: %v", err)
	}
	err = json.Unmarshal(accountJson, &account)
	if err != nil {
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

// RenewToken lookups the currently used token and tries to renew it by a given TTL if it's renewable.
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
	path := vault.getAwsCredentialsPath()
	secret, err := vault.client.Logical().Read(path)
	if err != nil {
		internal.AwsDynCredentialsRequestErrors.Inc()
		return nil, fmt.Errorf("could not gather dynamic credentials: %v", err)
	}
	return mapVaultAwsCredentialResponse(*secret)
}

func (vault *VaultBackend) getAwsCredentialsPath() string {
	return fmt.Sprintf("/aws/creds/%s", awsRole)
}

func (vault *VaultBackend) getAccountPath(hash string) string {
	return fmt.Sprintf("%s/server/account/%s", vault.namespacedPrefix, hash)
}

func (vault *VaultBackend) getCertDataPath(domain string) string {
	return fmt.Sprintf("%s/client/data/%s", vault.namespacedPrefix, domain)
}

func (vault *VaultBackend) getCertMetadataPath(domain string) string {
	return fmt.Sprintf("%s/client/expiry/%s", vault.namespacedPrefix, domain)
}

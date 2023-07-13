package config

import (
	"net/url"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type VaultConfig struct {
	Addr       string `json:"vaultAddr" validate:"required,http_url"`
	AuthMethod string `json:"vaultAuthMethod" validate:"required,oneof=token approle"`
	Token      string `json:"vaultToken" validate:"required_if=RoleId ''"`

	RoleId       string `json:"vaultRoleId" validate:"required_if=Token ''"`
	SecretId     string `json:"vaultSecretId" validate:"excluded_unless=SecretIdFile '',required_if=SecretIdFile '' Token ''"`
	SecretIdFile string `json:"vaultSecretIdFile" validate:"excluded_unless=SecretId '',required_if=SecretId '' Token ''"`

	PathPrefix       string `json:"vaultPathPrefix" validate:"required,startsnotwith=/,startsnotwith=/secret,endsnotwith=/,ne=acmevault"`
	DomainPathFormat string `json:"domainPathFormat" validate:"omitempty,containsrune=%"`

	Kv2MountPath string `json:"vaultKv2MountPath" validate:"required,endsnotwith=/,startsnotwith=/"`

	AwsMountPath string `json:"vaultAwsMountPath" validate:"required,endsnotwith=/,startsnotwith=/"`
	AwsRole      string `json:"vaultAwsRole" validate:"required"`
}

func (conf *VaultConfig) Print() {
	PrintFields(conf, SensitiveFields...)
}

func DefaultVaultConfig() VaultConfig {
	var pathPrefix string
	parsed, err := url.Parse(letsEncryptStagingUrl)
	if err == nil {
		pathPrefix = strings.ToLower(parsed.Host)
	}

	return VaultConfig{
		PathPrefix:   pathPrefix,
		Token:        os.Getenv("VAULT_TOKEN"),
		Addr:         os.Getenv("VAULT_ADDR"),
		AwsRole:      "acmevault",
		AwsMountPath: "aws",
		Kv2MountPath: "secret",
	}
}

func (conf *VaultConfig) Validate() error {
	return validate.Struct(conf)
}

func (conf *VaultConfig) LoadSecretIdFromFile() bool {
	return len(conf.SecretIdFile) > 0
}

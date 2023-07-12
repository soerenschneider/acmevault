package config

import (
	"net/url"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type VaultConfig struct {
	AuthMethod       string `json:"vaultAuthMethod" validate:"required,oneof=token approle"`
	VaultToken       string `json:"vaultToken" validate:"required_if=RoleId ''"`
	VaultAddr        string `json:"vaultAddr" validate:"required,http_url"`
	RoleId           string `json:"vaultRoleId" validate:"required_if=VaultToken ''"`
	SecretId         string `json:"vaultSecretId" validate:"excluded_unless=SecretIdFile '',required_if=SecretIdFile '' VaultToken ''"`
	SecretIdFile     string `json:"vaultSecretIdFile" validate:"excluded_unless=SecretId '',required_if=SecretId '' VaultToken ''"`
	PathPrefix       string `json:"vaultPathPrefix" validate:"required,startsnotwith=/,startsnotwith=/secret"`
	DomainPathFormat string `json:"domainPathFormat" validate:"omitempty,containsrune=%"`
	Kv2MountPath     string `json:"vaultKv2MountPath" validate:"required,endsnotwith=/,startsnotwith=/"`
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
		VaultToken:   os.Getenv("VAULT_TOKEN"),
		VaultAddr:    os.Getenv("VAULT_ADDR"),
		Kv2MountPath: "secret",
	}
}

func (conf *VaultConfig) Validate() error {
	return validate.Struct(conf)
}

func (conf *VaultConfig) LoadSecretIdFromFile() bool {
	return len(conf.SecretIdFile) > 0
}

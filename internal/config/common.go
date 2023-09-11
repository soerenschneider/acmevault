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
	AuthMethod string `json:"vaultAuthMethod" validate:"required,oneof=token approle k8s implicit"`
	Token      string `json:"vaultToken" validate:"required_if=AuthMethod 'token'"`

	RoleId       string `json:"vaultRoleId" validate:"required_if=AuthMethod 'approle'"`
	SecretId     string `json:"vaultSecretId" validate:"excluded_unless=SecretIdFile '',required_if=SecretIdFile '' AuthMethod 'approle'"`
	SecretIdFile string `json:"vaultSecretIdFile" validate:"excluded_unless=SecretId '',required_if=SecretId '' AuthMethod 'approle'"`

	K8sRoleId    string `json:"vaultK8sRoleId" validate:"required_if=AuthMethod 'k8s'"`
	K8sMountPath string `json:"vaultK8sMountPath"`

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

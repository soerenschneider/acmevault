package config

import (
	"net/url"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type VaultConfig struct {
	Addr       string `yaml:"addr" validate:"required,http_url"`
	AuthMethod string `yaml:"authMethod" validate:"required,oneof=token approle k8s implicit"`
	Token      string `yaml:"token" validate:"required_if=AuthMethod 'token'"`

	RoleId       string `yaml:"roleId" validate:"required_if=AuthMethod 'approle'"`
	SecretId     string `yaml:"secretId" validate:"excluded_unless=SecretIdFile '',required_if=SecretIdFile '' AuthMethod 'approle'"`
	SecretIdFile string `yaml:"secretIdFile" validate:"excluded_unless=SecretId '',required_if=SecretId '' AuthMethod 'approle'"`

	K8sRoleId    string `yaml:"k8sRoleId" validate:"required_if=AuthMethod 'k8s'"`
	K8sMountPath string `yaml:"k8sMountPath"`

	PathPrefix       string `yaml:"pathPrefix" validate:"required,startsnotwith=/,startsnotwith=/secret,endsnotwith=/,ne=acmevault"`
	DomainPathFormat string `yaml:"domainPathFormat" validate:"omitempty,containsrune=%"`

	Kv2MountPath string `yaml:"kv2MountPath" validate:"required,endsnotwith=/,startsnotwith=/"`

	AwsMountPath string `yaml:"awsMountPath" validate:"required,endsnotwith=/,startsnotwith=/"`
	AwsRole      string `yaml:"awsRole" validate:"required"`
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

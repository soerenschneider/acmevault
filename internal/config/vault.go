package config

import (
	"net/url"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type VaultConfig struct {
	Addr       string `yaml:"addr" env:"ADDR" validate:"required,http_url"`
	AuthMethod string `yaml:"authMethod" env:"AUTH_METHOD" validate:"required,oneof=token approle kubernetes implicit"`
	Token      string `yaml:"token" env:"TOKEN" validate:"required_if=AuthMethod 'token'"`

	RoleId       string `yaml:"roleId" env:"APPROLE_ROLE_ID" validate:"required_if=AuthMethod 'approle'"`
	SecretId     string `yaml:"secretId" env:"APPROLE_SECRET_ID" validate:"excluded_unless=SecretIdFile '',required_if=SecretIdFile '' AuthMethod 'approle'"`
	SecretIdFile string `yaml:"secretIdFile" env:"APPROLE_SECRET_ID_FILE" validate:"excluded_unless=SecretId '',required_if=SecretId '' AuthMethod 'approle'"`

	K8sRoleId    string `yaml:"k8sRoleId" env:"K8S_ROLE_ID" validate:"required_if=AuthMethod 'kubernetes'"`
	K8sMountPath string `yaml:"k8sMountPath" env:"K8S_MOUNT" `

	PathPrefix       string `yaml:"pathPrefix" env:"PATH_PREFIX" validate:"required,startsnotwith=/,startsnotwith=/secret,endsnotwith=/,ne=acmevault"`
	DomainPathFormat string `yaml:"domainPathFormat" env:"DOMAIN_PATH_FORMAT" validate:"omitempty,containsrune=%"`

	Kv2MountPath string `yaml:"kv2MountPath" env:"KV2_MOUNT" validate:"required,endsnotwith=/,startsnotwith=/"`

	AwsMountPath string `yaml:"awsMountPath" env:"AWS_MOUNT" validate:"required,endsnotwith=/,startsnotwith=/"`
	AwsRole      string `yaml:"awsRole" env:"AWS_ROLE" validate:"required"`
}

func defaultVaultConfig() VaultConfig {
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

func (conf *VaultConfig) UseAutoRenewAuth() bool {
	return conf.AuthMethod != "token" && conf.AuthMethod != "implicit"
}

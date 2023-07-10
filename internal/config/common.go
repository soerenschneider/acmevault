package config

import (
	"net/url"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type VaultConfig struct {
	VaultToken       string `json:"vaultToken" validate:"required_if=RoleId ''"`
	VaultAddr        string `json:"vaultAddr" validate:"required,http_url"`
	RoleId           string `json:"vaultRoleId" validate:"required_if=VaultToken ''"`
	SecretId         string `json:"vaultSecretId" validate:"excluded_unless=SecretIdFile '',required_if=SecretIdFile '' VaultToken ''"`
	SecretIdFile     string `json:"vaultSecretIdFile" validate:"excluded_unless=SecretId '',required_if=SecretId '' VaultToken ''"`
	PathPrefix       string `json:"vaultPathPrefix" validate:"required,startsnotwith=/,startsnotwith=/secret"`
	DomainPathFormat string `json:"domainPathFormat" validate:"omitempty,containsrune=%"`
}

func (conf *VaultConfig) Print() {
	log.Info().Msgf("VaultAddr=%s", conf.VaultAddr)
	log.Info().Msgf("VaultPathPrefix=%s", conf.PathPrefix)
	if len(conf.RoleId) > 0 {
		log.Info().Msgf("VaultRoleId=%s", conf.RoleId)
	}
	if len(conf.SecretId) > 0 {
		log.Info().Msg("VaultSecretId=*** (Redacted)")
	}
	if len(conf.SecretIdFile) > 0 {
		log.Info().Msgf("VaultSecretIdFile=%s", conf.SecretIdFile)
	}
	if len(conf.VaultWrappedToken) > 0 {
		log.Info().Msg("VaultWrappedToken=*** (Redacted)")
	}
	if len(conf.VaultWrappedTokenFile) > 0 {
		log.Info().Msgf("VaultWrappedFile=%s", conf.VaultWrappedTokenFile)
	}
	if len(conf.VaultToken) > 0 {
		log.Info().Msgf("VaultToken=%s", "*** (Redacted)")
	}
	if conf.TokenIncreaseSeconds > 0 {
		log.Info().Msgf("TokenIncreaseSeconds=%d", conf.TokenIncreaseSeconds)
	}
	if conf.TokenIncreaseInterval > 0 {
		log.Info().Msgf("TokenIncreaseInterval=%d", conf.TokenIncreaseInterval)
	}
	if len(conf.DomainPathFormat) > 0 {
		log.Info().Msgf("DomainPathFormat=%s", conf.DomainPathFormat)
	}
}

func DefaultVaultConfig() VaultConfig {
	var pathPrefix string
	parsed, err := url.Parse(letsEncryptStagingUrl)
	if err == nil {
		pathPrefix = strings.ToLower(parsed.Host)
	}

	return VaultConfig{
		PathPrefix: pathPrefix,
		VaultToken: os.Getenv("VAULT_TOKEN"),
		VaultAddr:  os.Getenv("VAULT_ADDR"),
	}
}

func (conf *VaultConfig) Validate() error {
	return validate.Struct(conf)
}

func isFileWritable(fileName string) bool {
	file, err := os.OpenFile(fileName, os.O_WRONLY, 0600)
	if err != nil {
		if os.IsPermission(err) {
			return false
		}
	} else {
		defer file.Close()
	}
	return true
}

func (conf *VaultConfig) HasWrappedToken() bool {
	return len(conf.VaultWrappedToken) > 0 || len(conf.VaultWrappedTokenFile) > 0
}

func (conf *VaultConfig) LoadWrappedTokenFromFile() bool {
	return len(conf.VaultWrappedTokenFile) > 0
}

func (conf *VaultConfig) LoadSecretIdFromFile() bool {
	return len(conf.SecretIdFile) > 0
}

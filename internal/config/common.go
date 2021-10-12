package config

import (
	"errors"
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
)

type VaultConfig struct {
	VaultToken            string `json:"vaultToken"`
	VaultAddr             string `json:"vaultAddr"`
	SecretId              string `json:"vaultSecretId"`
	RoleId                string `json:"vaultRoleId"`
	TokenIncreaseSeconds  int    `json:"tokenIncreaseSeconds"`
	TokenIncreaseInterval int    `json:"tokenIncreaseInterval"`
	PathPrefix            string `json:"vaultPathPrefix"`
}

func (conf *VaultConfig) IsTokenIncreaseEnabled() bool {
	return conf.TokenIncreaseInterval > 0 || conf.TokenIncreaseSeconds > 0
}

func (conf *VaultConfig) Print() {
	log.Info().Msgf("VaultAddr=%s", conf.VaultAddr)
	if len(conf.RoleId) > 0 {
		log.Info().Msgf("VaultRoleId=%s", conf.RoleId)
	}
	if len(conf.SecretId) > 0 {
		log.Info().Msgf("VaultSecretId=%s", "*** (Redacted)")
	}
	if len(conf.VaultToken) > 0 {
		log.Info().Msgf("VaultRoleId=%s", "*** (Redacted)")
	}
	if conf.TokenIncreaseSeconds > 0 {
		log.Info().Msgf("TokenIncreaseSeconds=%d", conf.TokenIncreaseSeconds)
	}
	if conf.TokenIncreaseInterval > 0 {
		log.Info().Msgf("TokenIncreaseInterval=%d", conf.TokenIncreaseInterval)
	}
	// TODO: Check pathPrefix
}

func DefaultVaultConfig() VaultConfig {
	var pathPrefix string
	parsed, err := url.Parse(letsEncryptStagingUrl)
	if err == nil {
		pathPrefix = strings.ToLower(parsed.Host)
	}

	return VaultConfig{
		PathPrefix: pathPrefix,
	}
}

func (conf *VaultConfig) Validate() error {
	if len(conf.VaultAddr) == 0 {
		return errors.New("no Vault address defined")
	}

	validRoleIdCredentials := len(conf.SecretId) > 0 && len(conf.RoleId) > 0
	if !validRoleIdCredentials && len(conf.VaultToken) == 0 {
		return errors.New("neither valid 'App Role' credentials nor plain Vault token provided")
	}

	return nil
}

func (conf *VaultConfig) HasLoginToken() bool {
	return false
}

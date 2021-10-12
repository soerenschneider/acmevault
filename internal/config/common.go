package config

import (
	"errors"
	"fmt"
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
	log.Info().Msgf("PathPrefix=%s", conf.PathPrefix)
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
	addr, err := url.ParseRequestURI(conf.VaultAddr)
	if err != nil || addr.Scheme == "" || addr.Host == "" || addr.Port() == "" {
		return errors.New("can not parse supplied vault addr as url")
	}

	for _, prefix := range []string{"/", "secret/"} {
		if strings.HasPrefix(conf.PathPrefix, prefix) {
			return fmt.Errorf("vault path prefix must not start with %s", prefix)
		}
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

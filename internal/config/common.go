package config

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"strings"
)

type VaultConfig struct {
	VaultToken            string `json:"vaultToken"`
	VaultAddr             string `json:"vaultAddr"`
	SecretId              string `json:"vaultSecretId"`
	SecretIdFile          string `json:"vaultSecretIdFile"`
	VaultWrappingToken    string `json:"vaultWrappingToken"`
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
	if len(conf.VaultWrappingToken) > 0 {
		log.Info().Msg("VaultWrappingToken=*** (Redacted)")
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

	if len(conf.PathPrefix) == 0 {
		return errors.New("empty path prefix provided")
	}
	for _, prefix := range []string{"/", "secret/"} {
		if strings.HasPrefix(conf.PathPrefix, prefix) {
			return fmt.Errorf("vault path prefix must not start with %s", prefix)
		}
	}

	validRoleIdCredentials := (len(conf.SecretId) > 0 || len(conf.SecretIdFile) > 0) && len(conf.RoleId) > 0
	if !validRoleIdCredentials && len(conf.VaultToken) == 0 {
		return errors.New("neither valid 'App Role' credentials nor plain Vault token provided")
	}

	if conf.HasWrappingToken() && !conf.LoadSecretIdFromFile() {
		return errors.New("vaultWrappingToken specified but no vaultSecretIdFile specified to write acquired secret to")
	}

	if len(conf.SecretId) > 0 && conf.LoadSecretIdFromFile() {
		return errors.New("both secretId and secretIdFile specified, unsure what to do")
	}

	if conf.LoadSecretIdFromFile() && !isFileWritable(conf.SecretIdFile) {
		return errors.New("specified secretIdFile is not writable, quitting")
	}

	return nil
}

func isFileWritable(fileName string) bool {
	file, err := os.OpenFile(fileName, os.O_WRONLY, 0600)
	defer file.Close()
	if err != nil {
		if os.IsPermission(err) {
			return false
		}
	}
	return true
}

func (conf *VaultConfig) HasWrappingToken() bool {
	return len(conf.VaultWrappingToken) > 0
}

func (conf *VaultConfig) LoadSecretIdFromFile() bool {
	return len(conf.SecretIdFile) > 0
}

func (conf *VaultConfig) HasLoginToken() bool {
	return false
}

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"regexp"
)

const letsEncryptUrl = "https://acme-v02.api.letsencrypt.org/directory"
const letsEncryptStagingUrl = "https://acme-staging-v02.api.letsencrypt.org/directory"

var (
	emailRegex             = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	defaultIntervalSeconds = 60 * 60 * 12
	defaultMetricsAddr     = "127.0.0.1:9112"
)

type AcmeVaultServerConfig struct {
	VaultConfig
	AcmeConfig
	IntervalSeconds int      `json:"intervalSeconds"`
	Domains         []string `json:"domains"`
	MetricsAddr     string   `json:"metricsAddr"`
}

type AcmeConfig struct {
	Email           string `json:"email"`
	AcmeUrl         string `json:"acmeUrl"`
	AcmeDnsProvider string `json:"acmeDnsProvider"`
}

func isValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

func (conf AcmeVaultServerConfig) Validate() error {
	if len(conf.AcmeDnsProvider) == 0 {
		return errors.New("field `acmeDnsProvider` not configured")
	}

	if len(conf.Domains) == 0 {
		return errors.New("field `domains` not configured")
	}

	if !isValidEmail(conf.Email) {
		return errors.New("field `email` not configured (correctly)")
	}

	if conf.IntervalSeconds < 0 {
		return fmt.Errorf("field `intervalSeconds` can not be a negative number: %d", conf.IntervalSeconds)
	}

	if conf.IntervalSeconds > 86400 {
		return fmt.Errorf("field `intervalSeconds` shouldn't be > 86400: %d", conf.IntervalSeconds)
	}

	return conf.VaultConfig.Validate()
}

func (conf AcmeVaultServerConfig) Print() {
	log.Info().Msg("--- Server Config Start ---")
	conf.VaultConfig.Print()
	log.Info().Msgf("AcmeDomains=%s", conf.Domains)
	log.Info().Msgf("AcmeEmail=%s", conf.Email)
	log.Info().Msgf("AcmeUrl=%s", conf.AcmeUrl)
	log.Info().Msgf("IntervalSeconds=%d", conf.IntervalSeconds)
	log.Info().Msgf("MetricsAddr=%s", conf.MetricsAddr)
	log.Info().Msgf("AcmeDnsProvider=%s", conf.AcmeDnsProvider)
	log.Info().Msg("--- Server Config End ---")
}

func getDefaultServerConfig() AcmeVaultServerConfig {
	return AcmeVaultServerConfig{
		AcmeConfig: AcmeConfig{
			AcmeUrl: letsEncryptUrl,
		},
		IntervalSeconds: defaultIntervalSeconds,
		MetricsAddr:     defaultMetricsAddr,
	}
}

func AcmeVaultServerConfigFromFile(path string) (AcmeVaultServerConfig, error) {
	conf := getDefaultServerConfig()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, fmt.Errorf("can not read config from file %s: %v", path, err)
	}

	err = json.Unmarshal(content, &conf)
	return conf, err
}

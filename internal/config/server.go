package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/asaskevich/govalidator"
	"github.com/rs/zerolog/log"
)

const letsEncryptUrl = "https://acme-v02.api.letsencrypt.org/directory"
const letsEncryptStagingUrl = "https://acme-staging-v02.api.letsencrypt.org/directory"

var (
	defaultIntervalSeconds = 60 * 60 * 12
	defaultMetricsAddr     = "127.0.0.1:9112"
)

type AcmeVaultServerConfig struct {
	VaultConfig
	AcmeConfig
	IntervalSeconds int                 `json:"intervalSeconds"`
	Domains         []AcmeServerDomains `json:"domains"`
	MetricsAddr     string              `json:"metricsAddr"`
}

type AcmeServerDomains struct {
	Domain string   `json:"domain"`
	Sans   []string `json:"sans,omitempty"`
}

func (a AcmeServerDomains) Verify() error {
	if ok := govalidator.IsDNSName(a.Domain); !ok {
		return fmt.Errorf("invalid domain name: '%s' is not a domain name", a.Domain)
	}
	for _, domain := range a.Sans {
		if ok := govalidator.IsDNSName(domain); !ok {
			return fmt.Errorf("invalid sans domain name: '%s' is not a domain name", domain)
		}
	}

	return nil
}

func (a AcmeServerDomains) String() string {
	if len(a.Sans) > 0 {
		return fmt.Sprintf("%s (%v)", a.Domain, a.Sans)
	}

	return a.Domain
}

type AcmeConfig struct {
	Email                string   `json:"email" validate:"email"`
	AcmeUrl              string   `json:"acmeUrl"`
	AcmeDnsProvider      string   `json:"acmeDnsProvider"`
	AcmeCustomDnsServers []string `json:"acmeCustomDnsServers,omitempty" validate:"dive,ip"`
}

func (conf AcmeConfig) Validate() error {
	if len(conf.AcmeDnsProvider) == 0 {
		return errors.New("field `acmeDnsProvider` not configured")
	}
	_, err := url.Parse(conf.AcmeDnsProvider)
	if err != nil {
		return fmt.Errorf("could not parse `acmeDnsProvider`: %v", err)
	}

	if !govalidator.IsEmail(conf.Email) {
		return fmt.Errorf("field `email` not configured (correctly): %s", conf.Email)
	}

	return nil
}

func (conf AcmeConfig) Print() {
	log.Info().Msgf("AcmeEmail=%s", conf.Email)
	log.Info().Msgf("AcmeUrl=%s", conf.AcmeUrl)
	log.Info().Msgf("AcmeDnsProvider=%s", conf.AcmeDnsProvider)
	if len(conf.AcmeCustomDnsServers) > 0 {
		log.Info().Msgf("AcmeCustomDnsServers=%v", conf.AcmeCustomDnsServers)
	}
}

func (conf AcmeVaultServerConfig) Validate() error {
	if conf.IntervalSeconds < 0 {
		return fmt.Errorf("field `intervalSeconds` can not be a negative number: %d", conf.IntervalSeconds)
	}

	if conf.IntervalSeconds > 86400 {
		return fmt.Errorf("field `intervalSeconds` shouldn't be > 86400: %d", conf.IntervalSeconds)
	}

	if len(conf.Domains) == 0 {
		return errors.New("no domains configured")
	}
	for _, domain := range conf.Domains {
		err := domain.Verify()
		if err != nil {
			return err
		}
	}
	err := conf.AcmeConfig.Validate()
	if err != nil {
		return err
	}

	return conf.VaultConfig.Validate()
}

func (conf AcmeVaultServerConfig) Print() {
	log.Info().Msg("--- Server Config Start ---")
	conf.VaultConfig.Print()
	conf.AcmeConfig.Print()
	for index, domain := range conf.Domains {
		log.Info().Msgf("AcmeDomains[%d]=%s", index, domain.String())
	}

	log.Info().Msgf("IntervalSeconds=%d", conf.IntervalSeconds)
	log.Info().Msgf("MetricsAddr=%s", conf.MetricsAddr)
	log.Info().Msg("--- Server Config End ---")
}

func getDefaultServerConfig() AcmeVaultServerConfig {
	return AcmeVaultServerConfig{
		AcmeConfig: AcmeConfig{
			AcmeUrl: letsEncryptUrl,
		},
		IntervalSeconds: defaultIntervalSeconds,
		MetricsAddr:     defaultMetricsAddr,
		VaultConfig:     DefaultVaultConfig(),
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

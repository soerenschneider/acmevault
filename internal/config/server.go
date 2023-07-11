package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const letsEncryptUrl = "https://acme-v02.api.letsencrypt.org/directory"
const letsEncryptStagingUrl = "https://acme-staging-v02.api.letsencrypt.org/directory"

var (
	defaultIntervalSeconds = 43200
	defaultMetricsAddr     = "127.0.0.1:9112"
)

type AcmeVaultServerConfig struct {
	VaultConfig
	AcmeEmail            string              `json:"email" validate:"required,email"`
	AcmeUrl              string              `json:"acmeUrl" validate:"required,oneof=https://acme-v02.api.letsencrypt.org/directory https://acme-staging-v02.api.letsencrypt.org/directory"`
	AcmeDnsProvider      string              `json:"acmeDnsProvider"`
	AcmeCustomDnsServers []string            `json:"acmeCustomDnsServers,omitempty" validate:"dive,ip"`
	IntervalSeconds      int                 `json:"intervalSeconds" validate:"min=3600,max=86400"`
	Domains              []AcmeServerDomains `json:"domains" validate:"required,dive"`
	MetricsAddr          string              `json:"metricsAddr" validate:"tcp_addr"`
}

type AcmeServerDomains struct {
	Domain string   `json:"domain" validate:"required,fqdn"`
	Sans   []string `json:"sans,omitempty" validate:"dive,fqdn"`
}

func (a AcmeServerDomains) String() string {
	if len(a.Sans) > 0 {
		return fmt.Sprintf("%s (%v)", a.Domain, a.Sans)
	}

	return a.Domain
}

func (conf AcmeVaultServerConfig) Validate() error {
	return validate.Struct(conf)
}

func getDefaultServerConfig() AcmeVaultServerConfig {
	return AcmeVaultServerConfig{
		AcmeUrl:         letsEncryptUrl,
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

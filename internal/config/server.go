package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const letsEncryptUrl = "https://acme-v02.api.letsencrypt.org/directory"
const letsEncryptStagingUrl = "https://acme-staging-v02.api.letsencrypt.org/directory"

var (
	defaultIntervalSeconds = 43200
	defaultMetricsAddr     = "127.0.0.1:9112"
)

type AcmeVaultServerConfig struct {
	Vault                VaultConfig         `yaml:"vault" validate:"required"`
	AcmeEmail            string              `yaml:"email" validate:"required,email"`
	AcmeUrl              string              `yaml:"acmeUrl" validate:"required,oneof=https://acme-v02.api.letsencrypt.org/directory https://acme-staging-v02.api.letsencrypt.org/directory"`
	AcmeDnsProvider      string              `yaml:"acmeDnsProvider"`
	AcmeCustomDnsServers []string            `yaml:"acmeCustomDnsServers,omitempty" validate:"dive,ip"`
	IntervalSeconds      int                 `yaml:"intervalSeconds" validate:"min=3600,max=86400"`
	Domains              []AcmeServerDomains `yaml:"domains" validate:"required,dive"`
	MetricsAddr          string              `yaml:"metricsAddr" validate:"tcp_addr"`
}

type AcmeServerDomains struct {
	Domain string   `yaml:"domain" validate:"required,fqdn"`
	Sans   []string `yaml:"sans,omitempty" validate:"dive,fqdn"`
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
		Vault:           DefaultVaultConfig(),
	}
}

func AcmeVaultServerConfigFromFile(path string) (AcmeVaultServerConfig, error) {
	conf := getDefaultServerConfig()
	content, err := os.ReadFile(path)
	if err != nil {
		return conf, fmt.Errorf("can not read config from file %s: %v", path, err)
	}

	err = yaml.Unmarshal(content, &conf)
	return conf, err
}

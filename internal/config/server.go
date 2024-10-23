package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
	"gopkg.in/yaml.v3"
)

const letsEncryptUrl = "https://acme-v02.api.letsencrypt.org/directory"
const letsEncryptStagingUrl = "https://acme-staging-v02.api.letsencrypt.org/directory"

var (
	defaultIntervalSeconds = 43200
	defaultMetricsAddr     = "127.0.0.1:9112"
)

type AcmeVaultConfig struct {
	Vault                VaultConfig     `yaml:"vault" envPrefix:"VAULT_" validate:"required"`
	AcmeEmail            string          `yaml:"email" env:"ACME_EMAIL" validate:"required,email"`
	AcmeUrl              string          `yaml:"acmeUrl" env:"ACME_URL" validate:"required,oneof=https://acme-v02.api.letsencrypt.org/directory https://acme-staging-v02.api.letsencrypt.org/directory"`
	AcmeDnsProvider      string          `yaml:"acmeDnsProvider" env:"ACME_DNS_PROVIDER"`
	AcmeCustomDnsServers []string        `yaml:"acmeCustomDnsServers,omitempty" env:"ACME_CUSTOM_DNS_SERVERS" validate:"dive,ip"`
	IntervalSeconds      int             `yaml:"intervalSeconds" env:"INTERVAL_SECONDS" validate:"min=3600,max=86400"`
	Domains              []DomainsConfig `yaml:"domains" validate:"required,dive"`
	MetricsAddr          string          `yaml:"metricsAddr" env:"METRICS_ADDR" validate:"omitempty,tcp_addr"`
	Verbose              bool            `yaml:"verbose" env:"VERBOSE"`
}

type DomainsConfig struct {
	Domain string   `yaml:"domain" validate:"required,fqdn"`
	Sans   []string `yaml:"sans,omitempty" validate:"dive,fqdn"`
}

func (a DomainsConfig) String() string {
	if len(a.Sans) > 0 {
		return fmt.Sprintf("%s (%v)", a.Domain, a.Sans)
	}

	return a.Domain
}

func (conf AcmeVaultConfig) Validate() error {
	return validate.Struct(conf)
}

func getDefaultConfig() AcmeVaultConfig {
	return AcmeVaultConfig{
		AcmeUrl:         letsEncryptUrl,
		IntervalSeconds: defaultIntervalSeconds,
		MetricsAddr:     defaultMetricsAddr,
		Vault:           defaultVaultConfig(),
	}
}

func read(path string) (AcmeVaultConfig, error) {
	conf := getDefaultConfig()
	content, err := os.ReadFile(path)
	if err != nil {
		return conf, fmt.Errorf("can not read config from file %s: %v", path, err)
	}

	err = yaml.Unmarshal(content, &conf)
	return conf, err
}

func GetConfig(path string) (AcmeVaultConfig, error) {
	conf, err := read(path)
	if err != nil {
		return AcmeVaultConfig{}, err
	}

	opts := env.Options{
		Prefix: "ACMEVAULT_",
	}

	if err := env.ParseWithOptions(&conf, opts); err != nil {
		return AcmeVaultConfig{}, err
	}

	return conf, nil
}

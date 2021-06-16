package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

type AcmeVaultClientConfig struct {
	VaultConfig
	PrivateKeyPath string   `json:"privateKeyFile"`
	CertPath       string   `json:"certFile"`
	User           string   `json:"user"`
	Group          string   `json:"group"`
	Hook           []string `json:"hooks"`
	MetricsPath    string   `json:"metricsPath"`
}

func AcmeVaultClientConfigFromFile(path string) (AcmeVaultClientConfig, error) {
	conf := AcmeVaultClientConfig{}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, fmt.Errorf("can not read config from file %s: %v", path, err)
	}

	err = json.Unmarshal(content, &conf)
	return conf, err
}

func (conf AcmeVaultClientConfig) Validate() error {
	if len(conf.User) == 0 {
		return errors.New("no user specified")
	}

	if len(conf.Group) == 0 {
		return errors.New("no group specified")
	}

	if len(conf.CertPath) == 0 {
		return errors.New("no certPath defined")
	}

	if len(conf.PrivateKeyPath) == 0 {
		return errors.New("no privateKeyPath defined")
	}

	return conf.VaultConfig.Validate()
}

func (conf AcmeVaultClientConfig) Print() {
	log.Info().Msg("--- Client Config Start ---")
	conf.VaultConfig.Print()
	log.Info().Msgf("PrivateKeyPath=%s", conf.PrivateKeyPath)
	log.Info().Msgf("CertificatePath=%s", conf.CertPath)
	log.Info().Msgf("User=%s", conf.User)
	log.Info().Msgf("Group=%s", conf.Group)
	if len(conf.Hook) > 0 {
		log.Info().Msgf("Hooks=%v", conf.Hook)
	}
	if len(conf.MetricsPath) > 0 {
		log.Info().Msgf("MetricsPath=%v", conf.MetricsPath)
	}
	log.Info().Msg("--- Client Config End ---")
}

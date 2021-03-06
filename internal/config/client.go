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
	FsWriterConfig
	Domain      string   `json:"domain"`
	Hook        []string `json:"hooks"`
	MetricsPath string   `json:"metricsPath"`
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
	if len(conf.Domain) == 0 {
		return errors.New("Missing field `domain`")
	}

	err := conf.FsWriterConfig.Validate()
	if err != nil {
		return err
	}

	return conf.VaultConfig.Validate()
}

func (conf AcmeVaultClientConfig) Print() {
	log.Info().Msg("--- Client Config Start ---")
	conf.VaultConfig.Print()
	conf.FsWriterConfig.Print()
	log.Info().Msgf("Domain=%s", conf.Domain)
	if len(conf.Hook) > 0 {
		log.Info().Msgf("Hooks=%v", conf.Hook)
	}
	if len(conf.MetricsPath) > 0 {
		log.Info().Msgf("MetricsPath=%v", conf.MetricsPath)
	}
	log.Info().Msg("--- Client Config End ---")
}

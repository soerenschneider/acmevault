package config

import (
	"errors"
	"github.com/rs/zerolog/log"
)

type FsWriterConfig struct {
	PrivateKeyPath string `json:"privateKeyFile"`
	CertPath       string `json:"certFile"`
	PemPath        string `json:"pemPath"`
	Username       string `json:"user"`
	Group          string `json:"group"`
}

func (conf FsWriterConfig) Validate() error {
	if len(conf.Username) == 0 {
		return errors.New("no user specified")
	}

	if len(conf.Group) == 0 {
		return errors.New("no group specified")
	}

	emptyCertPath := len(conf.CertPath) == 0
	emptyPrivateKeyPath := len(conf.PrivateKeyPath) == 0
	emptyPemPath := len(conf.PemPath) == 0
	if emptyCertPath && emptyPrivateKeyPath && emptyPemPath {
		return errors.New("missing either certPath and privateKeyPath or pemPath")
	}

	if emptyPemPath && (emptyCertPath || emptyPrivateKeyPath) {
		return errors.New("certPath and privateKeyPath must both be specified when no pemPath is defined")
	}

	return nil
}

func (conf FsWriterConfig) Print() {
	if len(conf.PemPath) > 0 {
		log.Info().Msgf("PemPath=%s", conf.PemPath)
	}
	if len(conf.PrivateKeyPath) > 0 {
		log.Info().Msgf("PrivateKeyPath=%s", conf.PrivateKeyPath)
	}
	if len(conf.CertPath) > 0 {
		log.Info().Msgf("CertificatePath=%s", conf.CertPath)
	}
	log.Info().Msgf("Username=%s", conf.Username)
	log.Info().Msgf("Group=%s", conf.Group)
}

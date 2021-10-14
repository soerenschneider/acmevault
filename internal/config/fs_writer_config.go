package config

import (
	"errors"
	"github.com/rs/zerolog/log"
)

type FsWriterConfig struct {
	PrivateKeyFile string `json:"privateKeyFile"`
	CertFile       string `json:"certFile"`
	PemFile        string `json:"pemFile"`
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

	emptyCertFile := len(conf.CertFile) == 0
	emptyPrivateKeyFile := len(conf.PrivateKeyFile) == 0
	emptyPemFile := len(conf.PemFile) == 0
	if emptyCertFile && emptyPrivateKeyFile && emptyPemFile {
		return errors.New("missing either certPath and privateKeyPath or pemPath")
	}

	if emptyPemFile && (emptyCertFile || emptyPrivateKeyFile) {
		return errors.New("certPath and privateKeyPath must both be specified when no pemPath is defined")
	}

	return nil
}

func (conf FsWriterConfig) Print() {
	if len(conf.PemFile) > 0 {
		log.Info().Msgf("PemFile=%s", conf.PemFile)
	}
	if len(conf.PrivateKeyFile) > 0 {
		log.Info().Msgf("PrivateKeyFile=%s", conf.PrivateKeyFile)
	}
	if len(conf.CertFile) > 0 {
		log.Info().Msgf("CertificateFile=%s", conf.CertFile)
	}
	log.Info().Msgf("Username=%s", conf.Username)
	log.Info().Msgf("Group=%s", conf.Group)
}

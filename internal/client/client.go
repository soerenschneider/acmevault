package client

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
	"os/exec"
	"time"
)

const metricsSubsystem = "client"

type CertificateWriter interface {
	WriteBundle(*certstorage.AcmeCertificate) (bool, error)
}

type VaultAcmeClient struct {
	conf    config.AcmeVaultClientConfig
	storage certstorage.CertStorage
	writer  CertificateWriter
}

func NewAcmeVaultClient(conf config.AcmeVaultClientConfig, storage certstorage.CertStorage, writer CertificateWriter) (*VaultAcmeClient, error) {
	err := conf.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}

	if storage == nil {
		return nil, errors.New("supplied storage is nil")
	}

	if writer == nil {
		return nil, errors.New("supplied writer is nil")
	}

	return &VaultAcmeClient{
		conf:    conf,
		storage: storage,
		writer:  writer,
	}, nil
}

func (client VaultAcmeClient) RetrieveAndSave(domain string) error {
	defer client.storage.Logout()

	log.Info().Msg("Logging in to storage...")
	err := client.storage.Authenticate()
	if err != nil {
		return fmt.Errorf("could not login to storage subsystem: %v", err)
	}

	log.Info().Msgf("Trying to read full cert data from storage for domain %s", domain)
	cert, err := client.storage.ReadFullCertificateData(domain)
	if err != nil {
		return fmt.Errorf("could not read secret bundle from vault: %v", err)
	}

	expiryTimestamp, err := cert.GetExpiryTimestamp()
	if err != nil {
		internal.CertErrors.WithLabelValues("unknown-expiry")
		log.Error().Msgf("Can not determine lifetime of certificate: %v", err)
	} else {
		daysLeft := int64(expiryTimestamp.Sub(time.Now().UTC()).Hours() / 24)
		log.Info().Msgf("Successfully read secret for domain %s from vault, valid for %d days", cert.Domain, daysLeft)
	}

	log.Info().Msg("Writing received data to configured backend...")
	runHook, err := client.writer.WriteBundle(cert)
	if err != nil {
		internal.CertWriteError.WithLabelValues(metricsSubsystem).Inc()
		return fmt.Errorf("writing the data was not successful: %v", err)
	}
	internal.CertWrites.WithLabelValues(metricsSubsystem).Inc()
	log.Info().Msg("Successfully written data")

	if !runHook {
		log.Info().Msg("No update detected, not running any hooks")
		return nil
	}

	log.Info().Msgf("Noticed update when writing secrets for domain %s", cert.Domain)
	return executeHook(client.conf.Hook)
}

func executeHook(hook []string) error {
	if len(hook) == 0 {
		return nil
	}

	log.Info().Msgf("Executing hook '%s'", hook)
	cmd := exec.Command(hook[0], hook[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Error().Msgf("Error running hook: %v", err)
		internal.HooksExecutionErrors.Inc()
		log.Error().Msgf("Error while running hook: %v", err)
		return err
	}

	log.Info().Msg("Hook successfully invoked")
	return nil
}

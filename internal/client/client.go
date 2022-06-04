package client

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal/client/hooks"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/internal/metrics"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
	"time"
)

const metricsSubsystem = "client"

// CertificateWriter defines a method to write a received certificate to a pluggable backend.
type CertificateWriter interface {
	WriteBundle(*certstorage.AcmeCertificate) (bool, error)
}

// PostHook defines a mechanism to run a hook after a certificate has been updated.
type PostHook interface {
	Invoke() error
}

type VaultAcmeClient struct {
	conf     config.AcmeVaultClientConfig
	storage  certstorage.CertStorage
	writer   CertificateWriter
	postHook PostHook
}

func NewAcmeVaultClient(conf config.AcmeVaultClientConfig, storage certstorage.CertStorage, writer CertificateWriter, hook PostHook) (*VaultAcmeClient, error) {
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

	if hook == nil {
		hook = &hooks.NopHook{}
	}

	return &VaultAcmeClient{
		conf:     conf,
		storage:  storage,
		writer:   writer,
		postHook: hook,
	}, nil
}

func (client VaultAcmeClient) RetrieveAndSave(domain string) error {
	defer client.storage.Logout()

	log.Info().Msg("Logging in to storage...")
	err := client.storage.Authenticate()
	if err != nil {
		metrics.AuthErrors.Inc()
		return fmt.Errorf("could not login to storage subsystem: %v", err)
	}

	log.Info().Msgf("Trying to read full cert data from storage for domain %s", domain)
	cert, err := client.storage.ReadFullCertificateData(domain)
	if err != nil {
		metrics.CertReadErrors.WithLabelValues(domain).Inc()
		return fmt.Errorf("could not read secret bundle from vault: %v", err)
	}

	if cert == nil {
		metrics.CertErrors.WithLabelValues(domain, "empty-cert").Inc()
		return fmt.Errorf("no cert returned")
	}

	expiryTimestamp, err := cert.GetExpiryTimestamp()
	if err != nil {
		metrics.CertErrors.WithLabelValues(domain, "unknown-expiry").Inc()
		log.Error().Msgf("Can not determine lifetime of certificate: %v", err)
	} else {
		daysLeft := int64(expiryTimestamp.Sub(time.Now().UTC()).Hours() / 24)
		log.Info().Msgf("Successfully read secret for domain %s from vault, valid for %d days", cert.Domain, daysLeft)
		metrics.CertClientExpiryTimestamp.WithLabelValues(domain).Set(float64(expiryTimestamp.Unix()))
	}

	log.Info().Msg("Writing received data to configured backend...")
	runHook, err := client.writer.WriteBundle(cert)
	if err != nil {
		metrics.PersistCertErrors.WithLabelValues(domain).Inc()
		return fmt.Errorf("writing the data was not successful: %v", err)
	}

	if !runHook {
		log.Info().Msg("No update detected, not running any hooks")
		return nil
	}

	log.Info().Msgf("Noticed update when writing secrets for domain %s", cert.Domain)
	err = client.postHook.Invoke()
	if err != nil {
		metrics.HooksExecutionErrors.Inc()
		return fmt.Errorf("error while running hook: %v", err)
	}

	return nil
}

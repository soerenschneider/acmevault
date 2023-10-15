package server

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/internal/metrics"
	"github.com/soerenschneider/acmevault/internal/server/acme"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
	"go.uber.org/multierr"
)

type AcmeVault struct {
	acmeClient  acme.AcmeDealer
	certStorage CertStorage
	domains     []config.DomainsConfig
}

type CertStorage interface {
	// Authenticate authenticates against the storage subsystem and returns an error about the success of the operation.
	Authenticate() error

	// WriteCertificate writes the full certificate to the underlying storage.
	WriteCertificate(resource *certstorage.AcmeCertificate) error

	// ReadPublicCertificateData reads the public portion of the certificate data (without the private key) from the
	// storage subsystem. This is intended to be used by the server component that does not need to have permission
	// to read the full certificate data.
	ReadPublicCertificateData(domain string) (*certstorage.AcmeCertificate, error)

	// ReadFullCertificateData reads all data for a given certificate and is intended to be used by the client component.
	ReadFullCertificateData(domain string) (*certstorage.AcmeCertificate, error)

	// Logout cleans up and logs out of the storage subsystem.
	Logout() error
}

func New(domains []config.DomainsConfig, acmeClient acme.AcmeDealer, storage CertStorage) (*AcmeVault, error) {
	if len(domains) == 0 {
		return nil, errors.New("no domains given")
	}

	if nil == acmeClient {
		return nil, errors.New("no acmeClient client provided")
	}

	if nil == storage {
		return nil, errors.New("no storage provider given")
	}

	return &AcmeVault{
		acmeClient:  acmeClient,
		certStorage: storage,
		domains:     domains,
	}, nil
}

func (c *AcmeVault) CheckCerts(ctx context.Context, wg *sync.WaitGroup) error {
	err := c.certStorage.Authenticate()
	if err != nil {
		return err
	}

	metrics.ServerLatestIterationTimestamp.SetToCurrentTime()
	for _, domain := range c.domains {
		err = multierr.Append(err, c.obtainAndHandleCert(domain))
	}

	if err := c.certStorage.Logout(); err != nil {
		log.Error().Err(err).Msg("logging out of storage failed")
	}

	return errs
}

func (c *AcmeVault) obtainAndHandleCert(domain config.DomainsConfig) error {
	read, err := c.certStorage.ReadPublicCertificateData(domain.Domain)
	if err != nil || read == nil {
		log.Error().Err(err).Msgf("Error reading cert data from storage for domain '%s'", domain.Domain)
		log.Info().Msgf("Trying to obtain cert from configured ACME provider for domain %s", domain.Domain)
		obtained, err := c.acmeClient.ObtainCert(domain)
		metrics.CertificatesRetrieved.Inc()
		if err != nil {
			metrics.CertificatesRetrievalErrors.Inc()
			return fmt.Errorf("obtaining cert for domain %s failed: %v", domain.Domain, err)
		}
		return handleReceivedCert(obtained, c.certStorage)
	}

	log.Info().Msgf("Read cert data for domain %q", domain)
	renewCert, err := read.NeedsRenewal()
	if err != nil {
		log.Info().Msgf("Could not determine cert lifetime for %s, probably the cert is broken", domain)
	}

	if renewCert {
		renewed, err := c.acmeClient.RenewCert(read)
		metrics.CertificatesRenewals.Inc()
		if err != nil {
			metrics.CertificatesRenewErrors.Inc()
			return fmt.Errorf("renewing cert failed for domain %s: %v", domain, err)
		}
		return handleReceivedCert(renewed, c.certStorage)
	}
	return nil
}

func handleReceivedCert(cert *certstorage.AcmeCertificate, storage CertStorage) error {
	if cert == nil {
		return fmt.Errorf("received empty cert for domain %s, this is weird and should not happen", cert.Domain)
	}

	expiry, err := cert.GetExpiryTimestamp()
	if err != nil {
		metrics.CertServerExpiryTimestamp.WithLabelValues(cert.Domain).Set(float64(expiry.Unix()))
	} else {
		metrics.CertErrors.WithLabelValues(cert.Domain, "unknown-expiry")
	}

	err = storage.WriteCertificate(cert)
	if err != nil {
		metrics.CertWriteError.WithLabelValues("server").Inc()
		return fmt.Errorf("received valid certificate for domain %s but storing it failed: %v", cert.Domain, err)
	}

	metrics.CertWrites.WithLabelValues("server").Inc()
	return nil
}

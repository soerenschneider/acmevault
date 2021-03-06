package server

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/internal/metrics"
	"github.com/soerenschneider/acmevault/internal/server/acme"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
	"time"
)

const (
	// MinCertLifetime defines a certs minimum validity. If a certificate's lifetime is less than this threshold, it's
	// being renewed.
	MinCertLifetime = time.Duration(30*24) * time.Hour
)

type AcmeVaultServer struct {
	acmeClient  acme.AcmeDealer
	certStorage certstorage.CertStorage
	domains     []config.AcmeServerDomains
}

func NewAcmeVaultServer(domains []config.AcmeServerDomains, acmeClient acme.AcmeDealer, storage certstorage.CertStorage) (*AcmeVaultServer, error) {
	if len(domains) == 0 {
		return nil, errors.New("no domains given")
	}

	if nil == acmeClient {
		return nil, errors.New("no acmeClient client provided")
	}

	if nil == storage {
		return nil, errors.New("no storage provider given")
	}

	return &AcmeVaultServer{
		acmeClient:  acmeClient,
		certStorage: storage,
		domains:     domains,
	}, nil
}

func (c *AcmeVaultServer) CheckCerts() {
	c.certStorage.Authenticate()
	metrics.ServerLatestIterationTimestamp.SetToCurrentTime()
	for _, domain := range c.domains {
		err := c.obtainAndHandleCert(domain)
		if err != nil {
			log.Error().Msgf("error while handling received certificate: %v", err)
		}
	}
	c.certStorage.Logout()
}

func (c *AcmeVaultServer) obtainAndHandleCert(domain config.AcmeServerDomains) error {
	log.Info().Msgf("Trying to read certificate data for domain %s from storage", domain.Domain)
	read, err := c.certStorage.ReadPublicCertificateData(domain.Domain)
	if err != nil || read == nil {
		log.Error().Msgf("Error reading cert data from storage for domain %s: %v", domain.Domain, err)
		log.Info().Msgf("Trying to obtain cert from configured ACME provider for domain %s", domain.Domain)
		obtained, err := c.acmeClient.ObtainCert(domain)
		metrics.CertificatesRetrieved.Inc()
		if err != nil {
			metrics.CertificatesRetrievalErrors.Inc()
			return fmt.Errorf("obtaining cert for domain %s failed: %v", domain.Domain, err)
		}
		return handleReceivedCert(obtained, c.certStorage)
	}

	log.Info().Msgf("Successfully read cert data from storage for domain %s", domain)
	expiry, err := read.GetExpiryTimestamp()
	if err != nil {
		log.Info().Msgf("Could not determine cert lifetime for %s, probably the cert is broken", domain)
	} else {
		timeLeft := expiry.Sub(time.Now().UTC())
		if timeLeft > MinCertLifetime {
			metrics.CertServerExpiryTimestamp.WithLabelValues(domain.Domain).Set(float64(expiry.Unix()))
			log.Info().Msgf("Not renewing cert for domain %s, still valid for %v", domain, timeLeft)
			return nil
		}
		log.Info().Msgf("Cert for domain %s is only valid for %v, renewing...", domain, timeLeft)
	}

	renewed, err := c.acmeClient.RenewCert(read)
	metrics.CertificatesRenewals.Inc()
	if err != nil {
		metrics.CertificatesRenewErrors.Inc()
		return fmt.Errorf("renewing cert failed for domain %s: %v", domain, err)
	}
	return handleReceivedCert(renewed, c.certStorage)
}

func handleReceivedCert(cert *certstorage.AcmeCertificate, storage certstorage.CertStorage) error {
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

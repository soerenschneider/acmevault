package server

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal"
	"github.com/soerenschneider/acmevault/internal/server/acme"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
	"time"
)

const (
	MinCertLifetime = time.Duration(30*24) * time.Hour
)

type AcmeVaultServer struct {
	acmeClient  acme.AcmeDealer
	certStorage certstorage.CertStorage
	domains     []string
}

func NewAcmeVaultServer(domains []string, acmeClient acme.AcmeDealer, storage certstorage.CertStorage) (*AcmeVaultServer, error) {
	if nil == domains || len(domains) == 0 {
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
	internal.ServerLatestIterationTimestamp.SetToCurrentTime()
	for _, domain := range c.domains {
		log.Info().Msgf("Acquiring certificate for domain %s", domain)
		err := c.obtainCertificate(domain)
		if err != nil {
			log.Error().Msgf("error while handling received certificate: %v", err)
		}
	}
}

func (c *AcmeVaultServer) obtainCertificate(domain string) error {
	read, err := c.certStorage.ReadCertificate(domain)
	if err != nil || read == nil {
		obtained, err := c.acmeClient.ObtainCert(domain)
		internal.CertificatesRetrieved.WithLabelValues(domain).Inc()
		if err != nil {
			internal.CertificatesRetrievalErrors.WithLabelValues(domain).Inc()
		}
		return handleReceivedCert(obtained, c.certStorage)
	}

	timeLeft, err := read.GetDurationUntilExpiry()
	if err != nil {
		log.Info().Msgf("Could not determine cert lifetime for %s, probably the cert is broken", domain)
	}

	if timeLeft > MinCertLifetime {
		log.Info().Msgf("Not renewing cert for domain %s, still valid for %v", domain, timeLeft)
		return nil
	}

	renewed, err := c.acmeClient.RenewCert(read)
	internal.CertificatesRenewals.WithLabelValues(domain).Inc()
	if err != nil {
		internal.CertificatesRenewErrors.WithLabelValues(domain).Inc()
	}
	return handleReceivedCert(renewed, c.certStorage)
}

func handleReceivedCert(cert *certstorage.AcmeCertificate, storage certstorage.CertStorage) error {
	if cert == nil {
		return fmt.Errorf("received empty cert for domain %s, this is weird and should not happen", cert.Domain)
	}

	expiry, err := cert.GetExpiryTimestamp()
	if err != nil {
		internal.CertExpiryTimestamp.WithLabelValues(cert.Domain).Set(float64(expiry.Unix()))
	} else {
		internal.CertErrors.WithLabelValues(cert.Domain, "unknown-expiry")
	}

	err = storage.WriteCertificate(cert)
	if err != nil {
		internal.CertWriteError.WithLabelValues(cert.Domain, "server").Inc()
		return fmt.Errorf("received valid certificate for domain %s but storing it failed: %v", cert.Domain, err)
	}

	internal.CertWrites.WithLabelValues(cert.Domain, "server").Inc()
	return nil
}

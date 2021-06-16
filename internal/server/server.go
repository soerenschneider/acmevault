package server

import (
	"acmevault/internal"
	"acmevault/internal/server/acme"
	"acmevault/pkg/certstorage"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
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
		c.obtainCertificate(domain)
	}
}

func (c *AcmeVaultServer) obtainCertificate(domain string) {
	res, err := c.certStorage.ReadCertificate(domain)
	if err != nil || res == nil {
		res, err = c.acmeClient.ObtainCert(domain)
		if err != nil {
			internal.CertificatesRetrievalErrors.WithLabelValues(domain).Inc()
		}
	} else {
		timeLeft, err := res.GetDurationUntilExpiry()
		if err != nil {
			log.Info().Msgf("Could not determine cert lifetime for %s, probably the cert is broken", domain)
		}

		if timeLeft > MinCertLifetime {
			log.Info().Msgf("Not renewing cert for domain %s, still valid for %v", domain, timeLeft)
			return
		}

		res, err = c.acmeClient.RenewCert(res)
		if err != nil {
			internal.CertificatesRenewErrors.WithLabelValues(domain).Inc()
		} else {
			internal.CertificatesRenewed.WithLabelValues(domain).Inc()
		}
	}

	err = handleReceivedCert(res, err, c.certStorage)
	if err != nil {
		log.Error().Msgf("error while handling received certificate: %v", err)
	}
}

func handleReceivedCert(cert *certstorage.AcmeCertificate, err error, storage certstorage.CertStorage) error {
	if err != nil {
		return fmt.Errorf("receiving certificate unsuccessful: %v", err)
	}

	if cert == nil {
		return fmt.Errorf("received empty cert for domain %s, this is weird and should not happen", cert.Domain)
	}

	internal.CertificatesRetrieved.WithLabelValues(cert.Domain).Inc()
	expiry, err := cert.GetExpiryTimestamp()
	if err != nil {
		internal.CertExpiryTimestamp.WithLabelValues(cert.Domain).Set(float64(expiry.Unix()))
	} else {
		internal.CertErrors.WithLabelValues(cert.Domain, "expiry")
	}

	err = storage.WriteCertificate(cert)
	if err != nil {
		internal.CertWriteError.WithLabelValues(cert.Domain, "server").Inc()
		return fmt.Errorf("received valid certificate for domain %s but storing it failed: %v", cert.Domain, err)
	}

	internal.CertWrites.WithLabelValues(cert.Domain, "server").Inc()
	return nil
}

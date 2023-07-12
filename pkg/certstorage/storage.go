package certstorage

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-acme/lego/v4/registration"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal/metrics"
)

const (
	// MinCertLifetime defines a certs minimum validity. If a certificate's lifetime is less than this threshold, it's
	// being renewed.
	MinCertLifetime = time.Duration(24*30) * time.Hour
	Skew            = time.Duration(24*60) * time.Hour
)

var ErrNotFound = errors.New("not found")
var ErrPermissionDenied = errors.New("permission denied")
var rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404

type CertStorage interface {
	// Authenticate authenticates against the storage subsystem and returns an error about the success of the operation.
	Authenticate() error

	// WriteCertificate writes the full certificate to the underlying storage.
	WriteCertificate(resource *AcmeCertificate) error

	// ReadPublicCertificateData reads the public portion of the certificate data (without the private key) from the
	// storage subsystem. This is intended to be used by the server component that does not need to have permission
	// to read the full certificate data.
	ReadPublicCertificateData(domain string) (*AcmeCertificate, error)

	// ReadFullCertificateData reads all data for a given certificate and is intended to be used by the client component.
	ReadFullCertificateData(domain string) (*AcmeCertificate, error)

	// Logout cleans up and logs out of the storage subsystem.
	Logout() error
}

var ErrAccountNotFound = errors.New("account not found")

type AccountStorage interface {
	// Authenticate authenticates against the storage subsystem and returns an error about the success of the operation.
	Authenticate() error

	// WriteAccount writes an ACME account to the storage.
	WriteAccount(AcmeAccount) error

	// ReadAccount reads the ACME account data for a given email address from the storage.
	ReadAccount(email string) (*AcmeAccount, error)

	// Logout cleans up and logs out of the storage subsystem.
	Logout() error
}

type AcmeCertificate struct {
	Domain            string `json:"domain"`
	CertURL           string `json:"certUrl"`
	CertStableURL     string `json:"certStableUrl"`
	PrivateKey        []byte `json:"-"`
	Certificate       []byte `json:"-"`
	IssuerCertificate []byte `json:"-"`
	CSR               []byte `json:"-"`
}

func (cert *AcmeCertificate) NeedsRenewal() (bool, error) {
	expiry, err := cert.GetExpiryTimestamp()
	if err != nil {
		return false, fmt.Errorf("could not determine cert expiry for domain '%s': %v", cert.Domain, err)
	}

	metrics.CertServerExpiryTimestamp.WithLabelValues(cert.Domain).Set(float64(expiry.Unix()))
	timeLeft := expiry.Sub(time.Now().UTC())
	log.Info().Msgf("Not renewing cert for domain %s, still valid for %v", cert.Domain, timeLeft)

	if timeLeft > MinCertLifetime && timeLeft <= Skew {
		if rnd.Intn(100) >= 97 {
			log.Info().Msgf("Earlier renewal of cert for domain %s to distribute cert expires (%v)", cert.Domain, timeLeft)
			return true, nil
		}
	}

	return timeLeft <= MinCertLifetime, nil
}

func (cert *AcmeCertificate) GetExpiryTimestamp() (time.Time, error) {
	block, _ := pem.Decode(cert.Certificate)
	if block == nil {
		return time.Time{}, errors.New("could not parse pem block from cert")
	}

	parsed, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse certificate: %v", err)
	}

	if len(parsed) == 0 {
		return time.Time{}, errors.New("no (valid) certificate data found")
	}
	return parsed[0].NotAfter, nil
}

func (cert *AcmeCertificate) GetDurationUntilExpiry() (time.Duration, error) {
	expiry, err := cert.GetExpiryTimestamp()
	if err != nil {
		metrics.CertErrors.WithLabelValues(cert.Domain, "unknown-expiry")
		return -1, err
	}

	return expiry.Sub(time.Now().UTC()), nil
}

func (cert *AcmeCertificate) AsPem() (pem string) {
	pem += fmt.Sprintln(string(cert.Certificate))
	pem += fmt.Sprint(string(cert.PrivateKey))
	return
}

type AcmeAccount struct {
	Email        string
	Key          crypto.PrivateKey
	Registration *registration.Resource
}

func (account AcmeAccount) IsInitialized() bool {
	if account.Key == nil {
		return false
	}

	if account.Registration == nil {
		return false
	}

	return true
}

func (account AcmeAccount) Validate() {
	if account.Key == nil {
		log.Fatal().Msg("No registration key provided")
	}

	if account.Registration == nil {
		log.Fatal().Msg("Registration empty")
	}
}

func (account AcmeAccount) GetEmail() string {
	return account.Email
}

func (account AcmeAccount) GetRegistration() *registration.Resource {
	return account.Registration
}

func (account AcmeAccount) GetPrivateKey() crypto.PrivateKey {
	return account.Key
}

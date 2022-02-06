package certstorage

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/registration"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal/server/metrics"
	"time"
)

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
	Logout()
}

type AccountStorage interface {
	// Authenticate authenticates against the storage subsystem and returns an error about the success of the operation.
	Authenticate() error

	// WriteAccount writes an ACME account to the storage.
	WriteAccount(AcmeAccount) error

	// ReadAccount reads the ACME account data for a given email address from the storage.
	ReadAccount(email string) (*AcmeAccount, error)

	// Logout cleans up and logs out of the storage subsystem.
	Logout()
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
		metrics.CertErrors.WithLabelValues("unknown-expiry")
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

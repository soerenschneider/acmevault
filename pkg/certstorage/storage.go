package certstorage

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/registration"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal"
	"time"
)

type CertStorage interface {
	WriteCertificate(resource *AcmeCertificate) error
	ReadCertificate(domain string) (*AcmeCertificate, error)
	Cleanup()
}

type AccountStorage interface {
	WriteAccount(AcmeAccount) error
	ReadAccount(hash string) (*AcmeAccount, error)
	Cleanup()
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
		internal.CertErrors.WithLabelValues(cert.Domain, "expiry")
		return -1, err
	}

	return expiry.Sub(time.Now().UTC()), nil
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

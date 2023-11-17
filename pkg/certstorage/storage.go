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
	MinCertLifetime = time.Duration(24) * time.Hour * 30
	Skew            = time.Duration(24) * time.Hour * 45
)

var ErrNotFound = errors.New("not found")
var ErrPermissionDenied = errors.New("permission denied")
var rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404

var ErrAccountNotFound = errors.New("account not found")

type AcmeCertificate struct {
	Domain            string `json:"domain"`
	CertURL           string `json:"certUrl"`
	CertStableURL     string `json:"certStableUrl"`
	PrivateKey        []byte `json:"-"`
	Certificate       []byte `json:"-"`
	IssuerCertificate []byte `json:"-"`
	CSR               []byte `json:"-"`
}

func niceTimeLeft(duration time.Duration) string {
	if duration.Hours() < 24 {
		return duration.String()
	}
	days := duration / (24 * time.Hour)
	duration = duration % (24 * time.Hour)

	hours := duration / time.Hour

	return fmt.Sprintf("%d days %d hours", days, hours)
}

func (cert *AcmeCertificate) NeedsRenewal() (bool, error) {
	expiry, err := cert.GetExpiryTimestamp()
	if err != nil {
		return false, fmt.Errorf("could not determine cert expiry for domain '%s': %v", cert.Domain, err)
	}

	metrics.CertServerExpiryTimestamp.WithLabelValues(cert.Domain).Set(float64(expiry.Unix()))
	timeLeft := expiry.Sub(time.Now().UTC())
	log.Debug().Msgf("Not renewing cert for domain %s, still valid for %v", cert.Domain, niceTimeLeft(timeLeft))

	if timeLeft > MinCertLifetime && timeLeft <= Skew {
		if rnd.Intn(100) >= 97 {
			log.Info().Msgf("Earlier renewal of cert for domain %s to distribute cert expires (%v)", cert.Domain, niceTimeLeft(timeLeft))
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

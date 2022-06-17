package acme

import (
	"crypto"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/registration"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
)

const (
	accountPrivateKeyType = certcrypto.RSA4096
	certPrivateKeyType    = certcrypto.RSA4096
)

type AcmeDealer interface {
	RegisterAccount() (*registration.Resource, error)
	ObtainCert(domain config.AcmeServerDomains) (*certstorage.AcmeCertificate, error)
	RenewCert(cert *certstorage.AcmeCertificate) (*certstorage.AcmeCertificate, error)
}

func GeneratePrivateKey() (crypto.PrivateKey, error) {
	return certcrypto.GeneratePrivateKey(accountPrivateKeyType)
}

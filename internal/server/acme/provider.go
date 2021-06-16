package acme

import (
	"acmevault/pkg/certstorage"
	"crypto"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/registration"
)

const (
	accountPrivateKeyType = certcrypto.RSA4096
	certPrivateKeyType    = certcrypto.RSA4096
)

type AcmeDealer interface {
	RegisterAccount() (*registration.Resource, error)
	ObtainCert(domain string) (*certstorage.AcmeCertificate, error)
	RenewCert(cert *certstorage.AcmeCertificate) (*certstorage.AcmeCertificate, error)
}

func GeneratePrivateKey() (crypto.PrivateKey, error) {
	return certcrypto.GeneratePrivateKey(accountPrivateKeyType)
}

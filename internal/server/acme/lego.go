package acme

import (
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
)

const DnsProviderRoute53 = "route53"

type GoLego struct {
	client *lego.Client
}

func NewGoLegoDealer(accountStorage certstorage.AccountStorage, acmeConfig config.AcmeConfig, dnsProvider challenge.Provider) (*GoLego, error) {
	log.Info().Msgf("Trying to read account details for %s from vault...", acmeConfig.Email)
	account, err := accountStorage.ReadAccount(acmeConfig.Email)
	// TODO: Introduce customize error types that signal whether to continue or not
	if err != nil {
		log.Info().Msgf("Received error when reading account %s: %v", acmeConfig.Email, err)
	}
	registerNewAccount := account == nil || err != nil
	if registerNewAccount {
		log.Info().Msg("No (valid) account data found in vault, attempting to register a new account")
		key, _ := GeneratePrivateKey()
		account = &certstorage.AcmeAccount{
			Email: acmeConfig.Email,
			Key:   key,
		}
	} else {
		log.Info().Msgf("Successfully read account from storage for %s", account.GetEmail())
	}

	legoConfig := lego.NewConfig(account)
	legoConfig.Certificate.KeyType = certPrivateKeyType
	if len(acmeConfig.AcmeUrl) > 0 {
		legoConfig.CADirURL = acmeConfig.AcmeUrl
	}

	l := &GoLego{}
	l.client, err = lego.NewClient(legoConfig)
	if err != nil {
		return nil, err
	}

	if registerNewAccount {
		registration, err := l.RegisterAccount()
		if err != nil {
			return nil, fmt.Errorf("can not register new account: %v", err)
		}
		account.Registration = registration
		err = accountStorage.WriteAccount(*account)
		if err != nil {
			log.Info().Msgf("Error while writing account information: %v", err)
		}
	}

	err = l.client.Challenge.SetDNS01Provider(dnsProvider)
	if err != nil {
		return nil, fmt.Errorf("could not set dns challenge: %v", err)
	}

	return l, nil
}

func (l *GoLego) RegisterAccount() (*registration.Resource, error) {
	return l.client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
}

func (l *GoLego) ObtainCert(domain string) (*certstorage.AcmeCertificate, error) {
	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}

	legoCert, err := l.client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}
	acmeCert := fromLego(legoCert)
	return &acmeCert, nil
}

func (l *GoLego) RenewCert(cert *certstorage.AcmeCertificate) (*certstorage.AcmeCertificate, error) {
	if cert == nil {
		return nil, errors.New("empty certificate provided")
	}

	oldLego := toLego(cert)
	oldLego.PrivateKey = nil
	newlegoCert, err := l.client.Certificate.Renew(oldLego, false, false, "")
	if err != nil {
		return nil, err
	}
	acmeCert := fromLego(newlegoCert)
	return &acmeCert, nil
}

func toLego(other *certstorage.AcmeCertificate) certificate.Resource {
	if nil == other {
		return certificate.Resource{}
	}

	return certificate.Resource{
		Domain:            other.Domain,
		CertURL:           other.CertURL,
		CertStableURL:     other.CertStableURL,
		PrivateKey:        other.PrivateKey,
		Certificate:       other.Certificate,
		IssuerCertificate: other.IssuerCertificate,
		CSR:               other.CSR,
	}
}

func fromLego(other *certificate.Resource) certstorage.AcmeCertificate {
	if nil == other {
		return certstorage.AcmeCertificate{}
	}

	return certstorage.AcmeCertificate{
		Domain:            other.Domain,
		CertURL:           other.CertURL,
		CertStableURL:     other.CertStableURL,
		PrivateKey:        other.PrivateKey,
		Certificate:       other.Certificate,
		IssuerCertificate: other.IssuerCertificate,
		CSR:               other.CSR,
	}
}

package acme

import (
	"errors"
	"fmt"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
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

func buildLegoClient(account *certstorage.AcmeAccount, acmeUrl string) (*GoLego, error) {
	legoConfig := lego.NewConfig(account)
	legoConfig.Certificate.KeyType = certPrivateKeyType
	if len(acmeUrl) > 0 {
		legoConfig.CADirURL = acmeUrl
	}

	var err error
	l := &GoLego{}
	l.client, err = lego.NewClient(legoConfig)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func getAccount(accountStorage AccountStorage, email string) (*certstorage.AcmeAccount, bool, error) {
	account, err := accountStorage.ReadAccount(email)
	if err == nil {
		log.Info().Msgf("retrieved account data for '%s'", email)
		return account, false, nil
	}

	if !errors.Is(err, certstorage.ErrNotFound) {
		return nil, false, err
	}

	log.Warn().Msg("No (valid) account data found in vault, attempting to register a new account")
	key, err := GeneratePrivateKey()
	if err != nil {
		return nil, false, fmt.Errorf("could not generate private key for ACME account: %w", err)
	}

	return &certstorage.AcmeAccount{
		Email: email,
		Key:   key,
	}, true, nil
}

func NewGoLegoDealer(accountStorage AccountStorage, conf config.AcmeVaultConfig, dnsProvider challenge.Provider) (*GoLego, error) {
	log.Info().Msgf("Trying to read account details for %s from vault...", conf.AcmeEmail)
	account, registerNewAccount, err := getAccount(accountStorage, conf.AcmeEmail)
	if err != nil {
		return nil, err
	}

	l, err := buildLegoClient(account, conf.AcmeUrl)
	if err != nil {
		return nil, err
	}

	if registerNewAccount {
		registration, err := l.RegisterAccount()
		if err != nil {
			return nil, fmt.Errorf("can not register new account: %v", err)
		}
		account.Registration = registration
		if err = accountStorage.WriteAccount(*account); err != nil {
			log.Warn().Err(err).Msg("Error while writing account information")
		}
	}

	var opts []dns01.ChallengeOption
	if len(conf.AcmeCustomDnsServers) > 0 {
		opts = append(opts, dns01.AddRecursiveNameservers(conf.AcmeCustomDnsServers))
	}

	if err := l.client.Challenge.SetDNS01Provider(dnsProvider, opts...); err != nil {
		return nil, fmt.Errorf("could not set dns challenge: %v", err)
	}

	return l, nil
}

func (l *GoLego) RegisterAccount() (*registration.Resource, error) {
	return l.client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
}

func (l *GoLego) ObtainCert(domain config.DomainsConfig) (*certstorage.AcmeCertificate, error) {
	domains := []string{domain.Domain}
	domains = append(domains, domain.Sans...)
	request := certificate.ObtainRequest{
		Domains: domains,
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
	opts := &certificate.RenewOptions{
		Bundle:         false,
		PreferredChain: "",
		MustStaple:     false,
	}
	newlegoCert, err := l.client.Certificate.RenewWithOptions(oldLego, opts)
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

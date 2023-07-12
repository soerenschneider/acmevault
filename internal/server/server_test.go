package server

import (
	"testing"

	"github.com/go-acme/lego/v4/registration"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
	"github.com/stretchr/testify/mock"
)

func TestServerHappyPathRenewal(t *testing.T) {
	dealer := &MockAcmeDealer{}
	certStorage := &MockStorage{}
	server := AcmeVaultServer{
		acmeClient:  dealer,
		certStorage: certStorage,
		domains:     []config.AcmeServerDomains{{Domain: "example.com"}},
	}

	old := &certstorage.AcmeCertificate{}
	new := &certstorage.AcmeCertificate{}
	certStorage.On("ReadPublicCertificateData", mock.Anything).Return(old, nil)
	dealer.On("RenewCert").Return(new, nil)
	certStorage.On("WriteCertificate", new).Return(nil)
	err := server.obtainAndHandleCert(server.domains[0])
	if err != nil {
		t.Fail()
	}
}

type MockAcmeDealer struct {
	mock.Mock
}

func (m *MockAcmeDealer) RegisterAccount() (*registration.Resource, error) {
	args := m.Called()
	if nil == args.Get(0) {
		return nil, args.Error(1)
	}
	return args.Get(0).(*registration.Resource), args.Error(1)
}

func (m *MockAcmeDealer) ObtainCert(domains config.AcmeServerDomains) (*certstorage.AcmeCertificate, error) {
	args := m.Called()
	if nil == args.Get(0) {
		return nil, args.Error(1)
	}
	return args.Get(0).(*certstorage.AcmeCertificate), args.Error(1)
}

func (m *MockAcmeDealer) RenewCert(cert *certstorage.AcmeCertificate) (*certstorage.AcmeCertificate, error) {
	args := m.Called()
	if nil == args.Get(0) {
		return nil, args.Error(1)
	}

	return args.Get(0).(*certstorage.AcmeCertificate), args.Error(1)
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Authenticate() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) WriteCertificate(resource *certstorage.AcmeCertificate) error {
	args := m.Called(resource)
	return args.Error(0)
}

func (m *MockStorage) ReadPublicCertificateData(domain string) (*certstorage.AcmeCertificate, error) {
	args := m.Called(domain)

	if args.Get(0) == nil {
		return nil, nil
	}
	return args.Get(0).(*certstorage.AcmeCertificate), args.Error(1)
}

func (m *MockStorage) ReadFullCertificateData(domain string) (*certstorage.AcmeCertificate, error) {
	args := m.Called(domain)

	if args.Get(0) == nil {
		return nil, nil
	}
	return args.Get(0).(*certstorage.AcmeCertificate), args.Error(1)
}

func (m *MockStorage) Logout() error {
	args := m.Called()
	return args.Error(0)
}

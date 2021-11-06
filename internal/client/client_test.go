package client

import (
	"errors"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestClientHappyPath(t *testing.T) {
	mockWriter := &MockCertWriter{}
	mockStorage := &MockStorage{}
	mockHook := &MockHook{}

	mockStorage.On("Authenticate").Return(nil)
	mockStorage.On("Logout")

	cert := &certstorage.AcmeCertificate{}
	mockStorage.On("ReadFullCertificateData", mock.Anything).Return(cert, nil)

	mockWriter.On("WriteBundle", cert).Return(true, nil)

	mockHook.On("Invoke").Return(nil)

	client := VaultAcmeClient{
		conf:     config.AcmeVaultClientConfig{},
		storage:  mockStorage,
		writer:   mockWriter,
		postHook: mockHook,
	}
	err := client.RetrieveAndSave("domain")
	mockWriter.AssertExpectations(t)
	mockStorage.AssertExpectations(t)

	if err != nil {
		t.Fail()
	}
}

func TestClientAuthFailure(t *testing.T) {
	mockWriter := &MockCertWriter{}
	mockStorage := &MockStorage{}
	mockHook := &MockHook{}

	mockStorage.On("Authenticate").Return(errors.New("auth error"))
	mockStorage.On("Logout")

	client := VaultAcmeClient{
		conf:     config.AcmeVaultClientConfig{},
		storage:  mockStorage,
		writer:   mockWriter,
		postHook: mockHook,
	}
	err := client.RetrieveAndSave("domain")
	mockWriter.AssertExpectations(t)
	mockStorage.AssertExpectations(t)

	if err == nil {
		t.Fail()
	}
}

func TestClientCertReadFailure(t *testing.T) {
	mockWriter := &MockCertWriter{}
	mockStorage := &MockStorage{}
	mockHook := &MockHook{}

	mockStorage.On("Authenticate").Return(nil)
	mockStorage.On("Logout")
	mockStorage.On("ReadFullCertificateData", mock.Anything).Return(nil, nil)

	client := VaultAcmeClient{
		conf:     config.AcmeVaultClientConfig{},
		storage:  mockStorage,
		writer:   mockWriter,
		postHook: mockHook,
	}
	err := client.RetrieveAndSave("domain")
	mockWriter.AssertExpectations(t)
	mockStorage.AssertExpectations(t)

	if err == nil {
		t.Fail()
	}
}

type MockCertWriter struct {
	mock.Mock
}

func (m *MockCertWriter) WriteBundle(bundle *certstorage.AcmeCertificate) (bool, error) {
	args := m.Called(bundle)
	return args.Bool(0), args.Error(1)
}

type MockHook struct {
	mock.Mock
}

func (m *MockHook) Invoke() error {
	args := m.Called()
	return args.Error(0)
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

func (m *MockStorage) Logout() {
	m.Called()
}

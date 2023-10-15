package acme

import "github.com/soerenschneider/acmevault/pkg/certstorage"

type AccountStorage interface {
	// Authenticate authenticates against the storage subsystem and returns an error about the success of the operation.
	Authenticate() error

	// WriteAccount writes an ACME account to the storage.
	WriteAccount(account certstorage.AcmeAccount) error

	// ReadAccount reads the ACME account data for a given email address from the storage.
	ReadAccount(email string) (*certstorage.AcmeAccount, error)

	// Logout cleans up and logs out of the storage subsystem.
	Logout() error
}

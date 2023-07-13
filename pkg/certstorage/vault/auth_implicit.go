package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
)

const defaultTokenFile = "~/.vault-token" // #nosec G101

type ImplicitAuth struct {
}

func NewImplicitAuth(token string) (*TokenAuth, error) {
	if len(token) == 0 {
		return nil, errors.New("empty token provided")
	}

	return &TokenAuth{token: token}, nil
}

func (t *ImplicitAuth) Logout(_ context.Context, _ *api.Client) error {
	return nil
}

func (t *ImplicitAuth) Login(_ context.Context, _ *api.Client) (*api.Secret, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	return &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken: token,
		},
	}, nil
}

func getToken() (string, error) {
	token := os.Getenv("VAULT_TOKEN")
	if len(token) > 0 {
		return token, nil
	}

	tokenFile := expandPath(defaultTokenFile)
	read, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("error reading file '%s': %v", tokenFile, err)
	}

	log.Info().Msgf("Using vault token from file '%s'", tokenFile)
	return string(read), nil
}

func expandPath(file string) string {
	if len(file) > 0 && file[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return file
		}

		if len(file) > 1 {
			return filepath.Join(homeDir, file[1:])
		}
		return homeDir
	}

	return file
}

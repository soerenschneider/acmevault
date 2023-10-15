package vault

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"
)

const (
	ImplicitAuthName = "implicit"
	defaultTokenFile = "~/.vault-token" // #nosec G101
)

type ImplicitAuth struct {
	tokenLocations []string
}

func NewImplicitAuth(tokenLocations ...string) (*ImplicitAuth, error) {
	if nil == tokenLocations {
		tokenLocations = []string{}
	}
	tokenLocations = append(tokenLocations, defaultTokenFile)

	return &ImplicitAuth{
		tokenLocations: tokenLocations,
	}, nil
}

func (t *ImplicitAuth) Logout(_ context.Context, _ *api.Client) error {
	return nil
}

func (t *ImplicitAuth) Login(_ context.Context, _ *api.Client) (*api.Secret, error) {
	token, err := t.getToken()
	if err != nil {
		return nil, err
	}

	return &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken: token,
		},
	}, nil
}

func (t *ImplicitAuth) getToken() (string, error) {
	token := os.Getenv("VAULT_TOKEN")
	if len(token) > 0 {
		return token, nil
	}

	var errs error
	for _, file := range t.tokenLocations {
		tokenFile := expandPath(file)
		log.Info().Msgf("Trying vault token from file '%s'", tokenFile)
		read, err := os.ReadFile(tokenFile)
		if err == nil {
			return string(read), nil
		}
		log.Warn().Msgf("Specified token file %s could not be read: %v", tokenFile, err)
		errs = multierr.Append(errs, err)
	}

	return "", fmt.Errorf("could not read token file(s): %w", errs)
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

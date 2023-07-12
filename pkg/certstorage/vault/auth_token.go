package vault

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/api"
)

type TokenAuth struct {
	token string
}

func NewTokenAuth(token string) (*TokenAuth, error) {
	if len(token) == 0 {
		return nil, errors.New("empty token provided")
	}

	return &TokenAuth{token: token}, nil
}

func (t *TokenAuth) Logout(_ context.Context, _ *api.Client) error {
	return nil
}

func (t *TokenAuth) Login(_ context.Context, _ *api.Client) (*api.Secret, error) {
	return &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken: t.token,
		},
	}, nil
}

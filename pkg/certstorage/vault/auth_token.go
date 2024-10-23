package vault

import (
	"context"

	"github.com/hashicorp/vault/api"
)

type TokenAuth struct {
	token string
}

func NewTokenAuth(token string) (*TokenAuth, error) {
	return &TokenAuth{token}, nil
}

func (t *TokenAuth) Login(_ context.Context, _ *api.Client) (*api.Secret, error) {
	ret := &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken: t.token,
		},
	}

	return ret, nil
}

package vault

import (
	"context"

	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

const ApproleAuthName = "approle"

func NewApproleAuth(roleId string, secretId *approle.SecretID) (*ApproleAuth, error) {
	va, err := approle.NewAppRoleAuth(roleId, secretId)
	if err != nil {
		return nil, err
	}

	return &ApproleAuth{
		auth: va,
	}, nil
}

type ApproleAuth struct {
	auth *approle.AppRoleAuth
}

func (a *ApproleAuth) Login(ctx context.Context, client *vault.Client) (*vault.Secret, error) {
	return a.auth.Login(ctx, client)
}

func (a *ApproleAuth) Logout(ctx context.Context, client *vault.Client) error {
	return client.Auth().Token().RevokeSelf("")
}

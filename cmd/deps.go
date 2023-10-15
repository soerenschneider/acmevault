package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/internal/server"
	"github.com/soerenschneider/acmevault/internal/server/acme"
	"github.com/soerenschneider/acmevault/pkg/certstorage/vault"
)

type deps struct {
	vaultAuth           vault.Auth
	storage             Storage
	credentialsProvider aws.CredentialsProvider
	dnsProvider         challenge.Provider
	acmeClient          acme.AcmeDealer

	done chan bool

	acmeVault *server.AcmeVault
}

type Storage interface {
	server.CertStorage
	acme.AccountStorage
	acme.AwsDynamicCredentialsBackend
}

func buildDeps(conf config.AcmeVaultConfig) *deps {
	deps := &deps{}
	var err error

	deps.vaultAuth, err = buildVaultAuth(conf.Vault)
	dieOnError(err, "could not build token auth")

	if conf.Vault.UseAutoRenewAuth() {
		log.Info().Msg("Building Vault auth auto renew wrapper...")
		deps.done = make(chan bool)
		deps.vaultAuth, err = vault.NewAutoRenew(deps.vaultAuth, deps.done)
		dieOnError(err, "could not build token auth")
	}

	deps.storage, err = vault.NewVaultBackend(conf.Vault, deps.vaultAuth)
	dieOnError(err, "could not generate desired backend")

	err = deps.storage.Authenticate()
	dieOnError(err, "Could not authenticate against storage")

	deps.credentialsProvider, err = acme.NewAwsDynamicCredentialsProvider(deps.storage)
	dieOnError(err, "could not build dynamic credentials provider")

	deps.dnsProvider, err = acme.BuildRoute53DnsProvider(deps.credentialsProvider)
	dieOnError(err, "could not build dns provider")

	deps.acmeClient, err = acme.NewGoLegoDealer(deps.storage, conf, deps.dnsProvider)
	dieOnError(err, "Could not initialize acme client")

	deps.acmeVault, err = server.New(conf.Domains, deps.acmeClient, deps.storage)
	dieOnError(err, "Couldn't build server")

	return deps
}

func dieOnError(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}

func buildVaultAuth(conf config.VaultConfig) (vault.Auth, error) {
	switch conf.AuthMethod {
	case vault.TokenAuthName:
		return vault.NewTokenAuth(conf.Token)
	case vault.ApproleAuthName:
		secretId := &approle.SecretID{
			FromFile:   conf.SecretIdFile,
			FromString: conf.SecretId,
		}
		return vault.NewApproleAuth(conf.RoleId, secretId)
	case vault.KubernetesAuthName:
		return vault.NewVaultKubernetesAuth(conf.K8sRoleId, conf.K8sMountPath)
	case vault.ImplicitAuthName:
		return vault.NewImplicitAuth()
	default:
		return nil, fmt.Errorf("no valid auth method: %s", conf.AuthMethod)
	}
}

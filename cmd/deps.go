package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/auth/kubernetes"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/internal/server"
	"github.com/soerenschneider/acmevault/internal/server/acme"
	"github.com/soerenschneider/acmevault/pkg/certstorage/vault"
)

const (
	vaultAuthToken      = "token"
	vaultAuthApprole    = "approle"
	vaultAuthKubernetes = "kubernetes"
	vaultAuthImplicit   = "implicit"
)

type deps struct {
	vaultAuth         api.AuthMethod
	vaultTokenRenewer *vault.TokenRenewer

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

	vaultClient, err := buildVaultClient(conf.Vault)
	dieOnError(err, "could not build vault client")

	deps.vaultAuth, err = buildVaultAuth(conf.Vault)
	dieOnError(err, "could not build token auth")

	if conf.Vault.UseAutoRenewAuth() {
		log.Info().Msg("Building Vault auth auto renew wrapper...")
		deps.done = make(chan bool)
		deps.vaultTokenRenewer, err = vault.NewTokenRenewer(vaultClient, deps.vaultAuth)
		dieOnError(err, "could not build token auth")
	}

	deps.storage, err = vault.NewVaultBackend(vaultClient, conf.Vault)
	dieOnError(err, "could not generate desired backend")

	// TODO
	//err = deps.storage.Authenticate()
	//dieOnError(err, "Could not authenticate against storage")

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

func buildVaultClient(conf config.VaultConfig) (*api.Client, error) {
	vaultConf := api.DefaultConfig()
	vaultConf.Address = conf.Addr
	vaultConf.MaxRetries = 3
	return api.NewClient(vaultConf)
}

func buildVaultAuth(conf config.VaultConfig) (api.AuthMethod, error) {
	switch conf.AuthMethod {
	case vaultAuthToken:
		return vault.NewTokenAuth(conf.Token)
	case vaultAuthApprole:
		secretId := &approle.SecretID{
			FromFile:   conf.SecretIdFile,
			FromString: conf.SecretId,
		}
		return approle.NewAppRoleAuth(conf.RoleId, secretId)
	case vaultAuthKubernetes:
		opts := []kubernetes.LoginOption{
			kubernetes.WithMountPath(conf.K8sMountPath),
		}
		return kubernetes.NewKubernetesAuth(conf.K8sRoleId, opts...)
	case vaultAuthImplicit:
		return vault.NewImplicitAuth()
	default:
		return nil, fmt.Errorf("no valid auth method: %s", conf.AuthMethod)
	}
}

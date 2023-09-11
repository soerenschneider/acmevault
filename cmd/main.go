package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/internal/metrics"
	"github.com/soerenschneider/acmevault/internal/server"
	"github.com/soerenschneider/acmevault/internal/server/acme"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
	"github.com/soerenschneider/acmevault/pkg/certstorage/vault"
)

func main() {
	configPath := parseCli()
	log.Info().Msgf("acmevault-server version %s, commit %s", internal.BuildVersion, internal.CommitHash)
	conf, err := config.AcmeVaultServerConfigFromFile(configPath)
	if err != nil {
		log.Fatal().Msgf("Could not load config: %v", err)
	}
	conf.Print()
	err = conf.Validate()
	if err != nil {
		log.Fatal().Msgf("Invalid configuration provided: %v", err)
	}
	NewAcmeVaultServer(conf)
}

const (
	envConfFile = "ACME_VAULT_CONFIG_FILE"
	cliConfFile = "conf"
	cliVersion  = "version"
)

func parseCli() (configFile string) {
	flag.StringVar(&configFile, cliConfFile, os.Getenv(envConfFile), "path to the config file")
	version := flag.Bool(cliVersion, false, "Print version and exit")
	flag.Parse()

	if *version {
		fmt.Printf("%s (revision %s)", internal.BuildVersion, internal.CommitHash)
		os.Exit(0)
	}

	if len(configFile) == 0 {
		log.Fatal().Msgf("No config file specified, use flag '-%s' or env var '%s'", cliConfFile, envConfFile)
	}

	if strings.HasPrefix(configFile, "~/") {
		configFile = path.Join(getUserHomeDirectory(), configFile[2:])
	}

	return
}

func getUserHomeDirectory() string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	return dir
}

func buildVaultAuth(conf config.AcmeVaultServerConfig) (vault.Auth, error) {
	switch conf.AuthMethod {
	case "token":
		return vault.NewTokenAuth(conf.Token)
	case "approle":
		secretId := &approle.SecretID{
			FromFile:   conf.SecretIdFile,
			FromString: conf.SecretId,
		}
		return vault.NewApproleAuth(conf.RoleId, secretId)
	case "k8s":
		return vault.NewVaultKubernetesAuth(conf.K8sRoleId, conf.K8sMountPath)
	case "implicit":
		return vault.NewImplicitAuth()
	default:
		return nil, fmt.Errorf("no valid auth method: %s", conf.AuthMethod)
	}
}

func dieOnError(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}

func NewAcmeVaultServer(conf config.AcmeVaultServerConfig) {
	vaultAuth, err := buildVaultAuth(conf)
	dieOnError(err, "could not build token auth")

	storage, err := vault.NewVaultBackend(conf.VaultConfig, vaultAuth)
	dieOnError(err, "could not generate desired backend")

	err = storage.Authenticate()
	dieOnError(err, "Could not authenticate against storage")

	dynamicCredentialsProvider, err := acme.NewDynamicCredentialsProvider(storage)
	dieOnError(err, "could not build dynamic credentials provider")

	dnsProvider, err := acme.BuildRoute53DnsProvider(dynamicCredentialsProvider)
	dieOnError(err, "could not build dns provider")

	acmeClient, err := acme.NewGoLegoDealer(storage, conf, dnsProvider)
	dieOnError(err, "Could not initialize acme client")

	acmeVaultServer, err := server.NewAcmeVaultServer(conf.Domains, acmeClient, storage)
	dieOnError(err, "Couldn't build server")

	err = Run(acmeVaultServer, storage, conf)
	dieOnError(err, "Couldn't start server")
}

func Run(acmeVault *server.AcmeVaultServer, storage certstorage.CertStorage, conf config.AcmeVaultServerConfig) error {
	if acmeVault == nil {
		return errors.New("empty acmevault provided")
	}
	if storage == nil {
		return errors.New("storage provider not provided")
	}
	err := conf.Validate()
	if err != nil {
		return fmt.Errorf("config invalid: %v", err)
	}

	go metrics.StartMetricsServer(conf.MetricsAddr)

	ticker := time.NewTicker(time.Duration(conf.IntervalSeconds) * time.Second)
	done := make(chan os.Signal, 1)
	signal.Notify(done,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	err = acmeVault.CheckCerts()
	if err != nil {
		log.Error().Err(err).Msg("error checking certs")
	}
	for {
		select {
		case <-ticker.C:
			err = acmeVault.CheckCerts()
			if err != nil {
				log.Error().Err(err).Msg("error checking certs")
			}
		case <-done:
			log.Info().Msg("Received signal, quitting")
			if err := storage.Logout(); err != nil {
				log.Warn().Err(err).Msg("Logging out failed")
			}
			ticker.Stop()
			os.Exit(0)
		}
	}
}

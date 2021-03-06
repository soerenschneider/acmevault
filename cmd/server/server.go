package main

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/cmd"
	"github.com/soerenschneider/acmevault/internal"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/internal/metrics"
	"github.com/soerenschneider/acmevault/internal/server"
	"github.com/soerenschneider/acmevault/internal/server/acme"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
	"github.com/soerenschneider/acmevault/pkg/certstorage/vault"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configPath := cmd.ParseCliFlags()
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

func NewAcmeVaultServer(conf config.AcmeVaultServerConfig) {
	storage, err := vault.NewVaultBackend(conf.VaultConfig)
	if err != nil {
		log.Fatal().Msgf("Could not generate desired backend: %v", err)
	}

	err = storage.Authenticate()
	if err != nil {
		log.Fatal().Msgf("Could not authenticate against storage: %v", err)
	}

	dynamicCredentialsProvider, _ := acme.NewDynamicCredentialsProvider(storage)
	dnsProvider, _ := acme.BuildRoute53DnsProvider(*dynamicCredentialsProvider)
	acmeClient, err := acme.NewGoLegoDealer(storage, conf.AcmeConfig, dnsProvider)
	if err != nil {
		log.Fatal().Msgf("Could not initialize acme client: %v", err)
	}

	acmeVaultServer, err := server.NewAcmeVaultServer(conf.Domains, acmeClient, storage)
	if err != nil {
		log.Fatal().Msgf("Couldn't build server: %v", err)
	}

	err = Run(acmeVaultServer, storage, conf)
	if err != nil {
		log.Fatal().Msgf("Couldn't start server: %v", err)
	}
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

	acmeVault.CheckCerts()
	for {
		select {
		case <-ticker.C:
			acmeVault.CheckCerts()
		case <-done:
			log.Info().Msg("Received signal, quitting")
			storage.Logout()
			ticker.Stop()
			os.Exit(0)
		}
	}
}

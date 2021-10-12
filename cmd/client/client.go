package main

import (
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/cmd"
	"github.com/soerenschneider/acmevault/internal"
	"github.com/soerenschneider/acmevault/internal/client"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/pkg/certstorage/vault"
	"os"
	"strings"
)

// Prefix of the configured AppRole role_ids for this tool
const roleIdPrefix = "acme-client-"

var conf config.AcmeVaultClientConfig

func main() {
	configPath := cmd.ParseCliFlags()
	log.Info().Msgf("acmevault-client version %s, commit %s", internal.BuildVersion, internal.CommitHash)
	conf, err := config.AcmeVaultClientConfigFromFile(configPath)
	if err != nil {
		log.Fatal().Msgf("Can't parse config: %v", err)
	}

	err = conf.Validate()
	if err != nil {
		log.Fatal().Msgf("Invalid config: %v", err)
	}
	conf.Print()

	storage, err := vault.NewVaultBackend(conf.VaultConfig, vault.NewPopulatedInMemoryTokenStorage(conf.VaultToken))
	if err != nil {
		log.Fatal().Msgf("Could not generate desired backend: %v", err)
	}

	writer, err := client.NewFsWriter(conf.CertPath, conf.PrivateKeyPath, conf.User, conf.Group)
	if err != nil {
		log.Fatal().Msgf("Could not create writer: %v", err)
	}

	client, err := client.NewAcmeVaultClient(conf, storage, writer)
	if err != nil {
		log.Fatal().Msgf("Could not build client: %v", err)
	}

	err = pickUpCerts(client, conf)
	exitCode := 0
	if err != nil {
		log.Error().Msgf("error while picking up and storing certificates: %v", err)
		exitCode = 1
	}
	os.Exit(exitCode)
}

func pickUpCerts(client *client.VaultAcmeClient, conf config.AcmeVaultClientConfig) error {
	if client == nil {
		return errors.New("empty client passed")
	}

	domain := strings.ReplaceAll(conf.RoleId, roleIdPrefix, "")
	err := client.RetrieveAndSave(domain)
	if len(conf.MetricsPath) > 0 {
		// shadow outer error
		err := internal.WriteMetrics(conf.MetricsPath)
		if err != nil {
			log.Error().Msgf("couldn't write metrics to file: %v", err)
		}
	}

	return err
}

package main

import (
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/cmd"
	"github.com/soerenschneider/acmevault/internal"
	"github.com/soerenschneider/acmevault/internal/client"
	"github.com/soerenschneider/acmevault/internal/client/hooks"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/pkg/certstorage/vault"
	"strings"
)

// Prefix of the configured AppRole role_ids for this tool
const roleIdPrefix = "acme-client-"

func main() {
	configPath := cmd.ParseCliFlags()
	log.Info().Msgf("acmevault-client version %s, commit %s", internal.BuildVersion, internal.CommitHash)
	conf, err := config.AcmeVaultClientConfigFromFile(configPath)
	if err != nil {
		die("Can't parse config: %v", err, conf)
	}

	conf.Print()
	err = conf.Validate()
	if err != nil {
		die("Invalid config: %v", err, conf)
	}

	storage, err := vault.NewVaultBackend(conf.VaultConfig)
	if err != nil {
		die("Could not generate desired backend: %v", err, conf)
	}

	writer, err := client.NewFsWriter(conf.FsWriterConfig)
	if err != nil {
		die("Could not create writer: %v", err, conf)
	}

	hook, err := hooks.NewCommandPostHook(conf.Hook)
	if err != nil {
		die("Could not create post hook: %v", err, conf)
	}

	client, err := client.NewAcmeVaultClient(conf, storage, writer, hook)
	if err != nil {
		die("Could not build client: %v", err, conf)
	}

	err = pickUpCerts(client, conf)
	die("error while picking up and storing certificates: %v", err, conf)
}

func die(msg string, err error, conf config.AcmeVaultClientConfig) {
	writeMetrics(conf)
	if err != nil {
		log.Fatal().Msgf(msg, err)
	}
}

func pickUpCerts(client *client.VaultAcmeClient, conf config.AcmeVaultClientConfig) error {
	if client == nil {
		return errors.New("empty client passed")
	}

	domain := strings.ReplaceAll(conf.RoleId, roleIdPrefix, "")
	return client.RetrieveAndSave(domain)
}

func writeMetrics(conf config.AcmeVaultClientConfig) {
	if len(conf.MetricsPath) > 0 {
		// shadow outer error
		err := internal.WriteMetrics(conf.MetricsPath)
		if err != nil {
			log.Error().Msgf("couldn't write metrics to file: %v", err)
		}
	}
}

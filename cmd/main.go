package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/internal/metrics"
)

func main() {
	configPath := parseCli()
	log.Info().Msgf("acmevault-server version %s, commit %s", internal.BuildVersion, internal.CommitHash)
	conf, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not load config")
	}

	if err := conf.Validate(); err != nil {
		log.Fatal().Err(err).Msgf("Invalid configuration provided")
	}

	deps := buildDeps(conf)
	run(conf, deps)
}

const (
	envConfFile = "ACMEVAULT_CONFIG_FILE"
	cliConfFile = "config"
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

func run(conf config.AcmeVaultConfig, deps *deps) {
	if len(conf.MetricsAddr) > 0 {
		go metrics.StartMetricsServer(conf.MetricsAddr)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Duration(conf.IntervalSeconds) * time.Second)
	defer func() {
		ticker.Stop()
		cancel()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	wg := &sync.WaitGroup{}

	if err := deps.acmeVault.CheckCerts(ctx, wg); err != nil {
		log.Error().Err(err).Msg("error checking certs")
	}

	stop := false
	for !stop {
		select {
		case <-ticker.C:
			err := deps.acmeVault.CheckCerts(ctx, wg)
			if err != nil {
				log.Error().Err(err).Msg("error checking certs")
			}
		case <-done:
			log.Info().Msg("Received signal, quitting")
			if err := deps.storage.Logout(); err != nil {
				log.Warn().Err(err).Msg("Logging out failed")
			}
			deps.done <- true
			cancel()
			ticker.Stop()
			stop = true
		}
	}

	log.Info().Msg("Waiting on other components")
	wg.Wait()
	log.Info().Msg("Done, bye!")
}

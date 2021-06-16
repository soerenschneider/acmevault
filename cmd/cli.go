package cmd

import (
	"flag"
	"github.com/rs/zerolog/log"
	"os"
	"os/user"
	"path"
	"strings"
)

const (
	envConfFile = "ACME_VAULT_CONFIG_FILE"
	cliConfFile = "conf"
)

func ParseCliFlags() (configFile string) {
	flag.StringVar(&configFile, cliConfFile, os.Getenv(envConfFile), "path to the config file")
	flag.Parse()

	if len(configFile) == 0 {
		log.Fatal().Msgf("No configfile specified, use flag '-%s' or env var '%s'", cliConfFile, envConfFile)
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

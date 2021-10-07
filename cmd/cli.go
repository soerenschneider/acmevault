package cmd

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal"
	"os"
	"os/user"
	"path"
	"strings"
)

const (
	envConfFile = "ACME_VAULT_CONFIG_FILE"
	cliConfFile = "conf"
	cliVersion  = "version"
)

func ParseCliFlags() (configFile string) {
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

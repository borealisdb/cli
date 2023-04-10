package config

import (
	"fmt"
	"log"
	"os"
)

const (
	TemplatesPath                  = "pkg/templates"
	CliTemplateFilename            = "config-template.ini"
	CliCredentialsTemplateFilename = "credentials-template.ini"

	CliConfigFilename      = "config.ini"
	CliCredentialsFilename = "credentials.ini"
	CliConfigType          = "ini"

	DefaultProjectPath                = "./borealis"
	DefaultProductionEnvironmentPath  = "/production"
	DefaultDevelopmentEnvironmentPath = "/development"
	DefaultNamespace                  = "default"

	CredentialsAuthSection = "auth"
	CredentialsTokenKey    = "token"
	SsoStartUrlKey         = "sso_start_url"
	HostKey                = "host"
	RootUrlPathKey         = "root_url_path"
	HelmChartUrl           = "../charts/borealis" // TODO change this once chart is hosted
	HelmReleaseName        = "borealis"
)

var CliConfigDefaultPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	CliConfigDefaultPath = fmt.Sprintf("%v/%v", home, ".borealis")
}

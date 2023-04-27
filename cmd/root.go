package cmd

import (
	"fmt"
	"github.com/borealisdb/cli/pkg/config"
	"github.com/borealisdb/cli/pkg/logger"
	"github.com/borealisdb/cli/pkg/utils"
	"github.com/borealisdb/go-sdk/api"
	"github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagEnvironment        = "environment" // For now only CliConfig is supported
	FlagLogLevel           = "log"
	FlagProjectFolder      = "project-folder"
	FlagBorealisConfigName = "borealis-config-name"
	FlagNamespace          = "namespace"
	FlagDryRun             = "dry-run"
	FlagHost               = "host"
	FlagRootPathUrl        = "root-path-url"
	FlagClusterName        = "cluster-name"
	FlagChartUrl           = "chart"
)

var log *logrus.Entry
var logLevel string
var projectFolder string
var borealisConfigName string
var dryRun bool
var namespace string
var credentialsConfigLocation string
var cliConfig config.Wrapper
var host string
var rootPathUrl string
var clusterName string
var environment string
var chartUrl string
var sdk api.API
var Version = "0.0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "borealisdb",
	Short:   "Borealis CLI",
	Version: Version,
	Long:    `Official CLI for Borealisdb https://github.com/borealisdb`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("Borealis CLI: %s", Version)
	},
	// Validation and flags cleanup
	// TODO add flags validation here
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		projectFolder = strings.TrimSuffix(projectFolder, "/")
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&environment, FlagEnvironment, "development", "")
	rootCmd.PersistentFlags().StringVar(&logLevel, FlagLogLevel, "info", "")
	rootCmd.PersistentFlags().StringVar(&projectFolder, FlagProjectFolder, config.DefaultProjectPath, "")
	rootCmd.PersistentFlags().StringVar(&namespace, FlagNamespace, config.DefaultNamespace, "")
	rootCmd.PersistentFlags().BoolVar(&dryRun, FlagDryRun, false, "")
	rootCmd.PersistentFlags().StringVar(&host, FlagHost, "", "")
	rootCmd.PersistentFlags().StringVar(&rootPathUrl, FlagRootPathUrl, "/borealis", "")
	rootCmd.PersistentFlags().StringVar(&clusterName, FlagClusterName, "my-cluster", "")
}

func initConfig() {
	log = logger.New(logLevel)
	log.Debugf("reading config...")
	viper.SetConfigName(config.CliConfigFilename)
	viper.AddConfigPath(config.CliConfigDefaultPath)
	viper.SetConfigType(config.CliConfigType)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		credentialsConfigLocation = fmt.Sprintf("%v/%v", config.CliConfigDefaultPath, config.CliCredentialsFilename)
		cliConfig = config.Wrapper{
			CredentialsLocation: credentialsConfigLocation,
		}
		if err := cliConfig.LoadCredentials(); err != nil {
			log.WithError(err).Fatal()
		}
		log.Debugf("using config file: %v and credentials file", viper.ConfigFileUsed())

		baseUrl := cliConfig.GetBaseUrl()
		sdk, err = api.New(baseUrl, defaultTokenFunc, api.Config{
			TlsCaLocation: cliConfig.GetCALocation(),
		})
		if err != nil {
			cobra.CheckErr(fmt.Sprintf("error initializing sdk: %v", err))
		}
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if initCmd.CalledAs() == "" {
				cobra.CheckErr(fmt.Errorf("%v. Please run: %v init", err, utils.CliAppName))
			}
		} else {
			// Config file was found but another error was produced
		}
	}
}

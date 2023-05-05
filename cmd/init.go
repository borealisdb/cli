package cmd

import (
	"fmt"
	"github.com/borealisdb/cli/pkg/config"
	"github.com/spf13/cobra"
	"net/url"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "OnInit the cli",
	Run: func(cmd *cobra.Command, args []string) {
		base := config.New(config.Params{
			Host:          host,
			RootPathUrl:   rootPathUrl,
			ProjectFolder: projectFolder,
		}, log)
		if err := base.GenerateDefaultConfigFiles(); err != nil {
			cobra.CheckErr(err)
		}

		// helm repo add borealisdb https://borealisdb.github.io/charts
		// chart url could be local, thus adding remote repo won't work
		if _, err := url.ParseRequestURI(chartUrl); err == nil {
			RunCommandOrDie(
				"helm",
				"repo",
				"add",
				config.HelmChartName,
				config.HelmChartUrl,
			)
		}

		RunCommandOrDie(
			"helm",
			"repo",
			"update",
		)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().StringVar(&chartUrl, FlagChartUrl, fmt.Sprintf("%v/borealis", config.HelmChartName), "")
	initCmd.PersistentFlags().StringVar(&host, FlagHost, "", "")
	initCmd.MarkPersistentFlagRequired(FlagHost)
}

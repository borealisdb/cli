package cmd

import (
	"fmt"
	"github.com/borealisdb/cli/pkg/config"
	"github.com/spf13/cobra"
	"net/url"
	"os/exec"
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
		err := base.GenerateDefaultConfigFiles()

		// helm repo add borealisdb https://borealisdb.github.io/charts
		// chart url could be local, thus adding remote repo won't work
		if _, err := url.ParseRequestURI(chartUrl); err == nil {
			helmRepoAddOutput, err := exec.Command(
				"helm",
				"repo",
				"add",
				config.HelmChartName,
				config.HelmChartUrl,
			).CombinedOutput()
			if err != nil {
				cobra.CheckErr(fmt.Sprint(err) + ": " + string(helmRepoAddOutput))
			}
			log.Infof(string(helmRepoAddOutput))
		}

		helmRepoUpdateOutput, err := exec.Command(
			"helm",
			"repo",
			"update",
		).CombinedOutput()
		if err != nil {
			cobra.CheckErr(fmt.Sprint(err) + ": " + string(helmRepoUpdateOutput))
		}
		log.Infof(string(helmRepoUpdateOutput))

		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().StringVar(&chartUrl, FlagChartUrl, config.HelmChartUrl, "")
	initCmd.PersistentFlags().StringVar(&host, FlagHost, "", "")
	initCmd.MarkPersistentFlagRequired(FlagHost)
}

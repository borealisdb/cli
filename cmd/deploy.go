package cmd

import (
	"fmt"
	"github.com/borealisdb/cli/pkg/config"
	"github.com/spf13/cobra"
	"os/exec"
	"path/filepath"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy templates and helm chart",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		secretsFilePath := filepath.Join(projectFolder, environment, "infrastructures/secrets.yaml")
		secretOutput, err := exec.Command("kubectl", "apply", "-f", secretsFilePath).CombinedOutput()
		if err != nil {
			cobra.CheckErr(fmt.Sprint(err) + ": " + string(secretOutput))
		}
		log.Infof(string(secretOutput))

		helmRepoUpdateOutput, err := exec.Command(
			"helm",
			"repo",
			"update",
		).CombinedOutput()
		if err != nil {
			cobra.CheckErr(fmt.Sprint(err) + ": " + string(helmRepoUpdateOutput))
		}
		log.Infof(string(helmRepoUpdateOutput))

		valuesFilePath := filepath.Join(projectFolder, environment, "infrastructures/values.yaml")
		helmOutput, err := exec.Command(
			"helm",
			"upgrade",
			"--install",
			"-f",
			valuesFilePath,
			"--create-namespace",
			"--namespace",
			namespace,
			config.HelmReleaseName,
			config.HelmChartUrl,
		).CombinedOutput()
		if err != nil {
			cobra.CheckErr(fmt.Sprint(err) + ": " + string(helmOutput))
		}
		log.Infof(string(helmOutput))
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.PersistentFlags().StringVar(&chartUrl, FlagChartUrl, config.HelmChartUrl, "")
}

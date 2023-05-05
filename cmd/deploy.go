package cmd

import (
	"github.com/borealisdb/cli/pkg/config"
	"github.com/spf13/cobra"
	"path/filepath"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy templates and helm chart",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		RunCommandOrDie(
			"kubectl",
			"apply",
			"-f",
			filepath.Join(projectFolder, environment, "infrastructures/secrets.yaml"),
		)

		RunCommandOrDie(
			"helm",
			"repo",
			"update",
		)

		RunCommandOrDie(
			"helm",
			"upgrade",
			"--install",
			"-f",
			filepath.Join(projectFolder, environment, "infrastructures/values.yaml"),
			"--create-namespace",
			"--namespace",
			namespace,
			config.HelmReleaseName,
			chartUrl,
		)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

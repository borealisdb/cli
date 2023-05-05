package cmd

import (
	"github.com/spf13/cobra"
	"path/filepath"
)

// clusterCreateCmd represents the clusterCreate command
var clusterDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		RunCommandOrDie(
			"kubectl",
			"apply",
			"-f",
			filepath.Join(projectFolder, environment, clusterName),
		)
	},
}

func init() {
	clusterCmd.AddCommand(clusterDeployCmd)
	clusterDeployCmd.PersistentFlags().StringVar(&clusterName, FlagClusterName, "", "")
	clusterDeployCmd.PersistentFlags().StringVar(&environment, FlagEnvironment, "", "")
	if err := clusterDeployCmd.MarkPersistentFlagRequired(FlagEnvironment); err != nil {
		cobra.CheckErr(err)
	}
	if err := clusterDeployCmd.MarkPersistentFlagRequired(FlagClusterName); err != nil {
		cobra.CheckErr(err)
	}
}

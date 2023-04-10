package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"path/filepath"
)

// clusterCreateCmd represents the clusterCreate command
var clusterDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		output, err := exec.Command("kubectl", "apply", "-f", filepath.Join(projectFolder, environment, clusterName)).CombinedOutput()
		if err != nil {
			cobra.CheckErr(fmt.Sprint(err) + ": " + string(output))
		}
		log.Infof(string(output))
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

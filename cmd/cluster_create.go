package cmd

import (
	"github.com/spf13/cobra"
)

// clusterCreateCmd represents the clusterCreate command
var clusterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Initialize a new cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	clusterCmd.AddCommand(clusterCreateCmd)
}

package cmd

import (
	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Group command for cluster functionality",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(clusterCmd)
}

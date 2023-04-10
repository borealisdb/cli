package cmd

import (
	"github.com/borealisdb/cli/pkg/config"
	"github.com/spf13/cobra"
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
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

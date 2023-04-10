package cmd

import (
	"github.com/borealisdb/cli/pkg/auth"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login and get a token back",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ssoStartUrl := cliConfig.GetSSOLoginURL()
		a := auth.Auth{
			Log:       log,
			CliConfig: cliConfig,
		}
		if err := a.Openbrowser(ssoStartUrl); err != nil {
			return
		}
		a.TokenListener()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

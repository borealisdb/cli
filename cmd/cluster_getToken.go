/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"github.com/borealisdb/go-sdk/api"

	"github.com/spf13/cobra"
)

// clusterGettokenCmd represents the clusterGettoken command
var clusterGettokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Get token to authenticate with postgres",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := sdk.GenerateClusterToken(api.GenerateClusterTokenRequest{
			ClusterName: clusterName,
		})
		if err != nil {
			log.Errorf(err.Error())
			cobra.CheckErr(err)
		}
		jsonBytes, err := json.MarshalIndent(token, "", "    ")
		if err != nil {
			log.Errorf(string(jsonBytes))
			cobra.CheckErr(err)
		}
		log.Infof(string(jsonBytes))
	},
}

func init() {
	clusterCmd.AddCommand(clusterGettokenCmd)
}

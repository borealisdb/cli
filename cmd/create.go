package cmd

import (
	"fmt"
	"github.com/borealisdb/cli/pkg/config"
	"github.com/borealisdb/cli/pkg/types"
	"github.com/borealisdb/cli/pkg/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create the project, it outputs an helm chart with default values and create the folder structures",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(projectFolder); err == nil || os.IsExist(err) {
			cobra.CheckErr(fmt.Sprintf("%v folder alreasy exist, skipping creation", projectFolder))
		}

		// Create folder structures
		for _, env := range []string{config.DefaultProductionEnvironmentPath, config.DefaultDevelopmentEnvironmentPath} {
			infraDest := filepath.Join(projectFolder, env, "infrastructures")
			if err := os.MkdirAll(infraDest, 0700); err != nil {
				cobra.CheckErr(fmt.Sprintf("could not create infrastructures destination folder: %v", err))
			}

			if err := utils.WriteAssetToFile(filepath.Join(config.TemplatesPath, "secrets.yaml"), filepath.Join(infraDest, "secrets.yaml")); err != nil {
				cobra.CheckErr(fmt.Sprintf("could not write secrets.yaml file to destination folder: %v", err))
			}

			if err := utils.RenderTemplate(
				types.Cluster{
					ClusterName: clusterName,
					Host:        host,
					RootUrlPath: rootPathUrl,
				},
				filepath.Join(config.TemplatesPath, "values.yaml"),
				filepath.Join(infraDest, "values.yaml"),
			); err != nil {
				cobra.CheckErr(fmt.Sprintf("could not render values.yaml: %v", err))
			}

			clusterDest := filepath.Join(projectFolder, env, clusterName)
			if err := os.MkdirAll(clusterDest, 0700); err != nil {
				cobra.CheckErr(fmt.Sprintf("could not create clusters template destination folder: %v", err))
			}

			clusterAssetPath := filepath.Join(config.TemplatesPath, "cluster.yaml")
			if err := utils.RenderTemplate(
				types.Cluster{
					ClusterName: clusterName,
					Host:        host,
					RootUrlPath: rootPathUrl,
				},
				clusterAssetPath,
				filepath.Join(clusterDest, "cluster.yaml"),
			); err != nil {
				cobra.CheckErr(fmt.Sprintf("could not render cluster: %v", err))
			}

			accountsAssetPath := filepath.Join(config.TemplatesPath, "accounts.yaml")
			if err := utils.RenderTemplate(
				types.Cluster{ClusterName: clusterName},
				accountsAssetPath,
				filepath.Join(clusterDest, "accounts.yaml"),
			); err != nil {
				cobra.CheckErr(fmt.Sprintf("could not render cluster: %v", err))
			}

			secretsAssetPath := filepath.Join(config.TemplatesPath, "cluster-secrets.yaml")
			if err := utils.RenderTemplate(
				types.Cluster{ClusterName: clusterName},
				secretsAssetPath,
				filepath.Join(clusterDest, "secrets.yaml"),
			); err != nil {
				cobra.CheckErr(fmt.Sprintf("could not render cluster secrets: %v", err))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVar(&clusterName, FlagClusterName, "", "")
	createCmd.PersistentFlags().StringVar(&host, FlagHost, "", "")
	if err := createCmd.MarkPersistentFlagRequired(FlagClusterName); err != nil {
		cobra.CheckErr(err)
	}
	if err := createCmd.MarkPersistentFlagRequired(FlagHost); err != nil {
		cobra.CheckErr(err)
	}
}

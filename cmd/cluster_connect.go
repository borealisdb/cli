package cmd

import (
	"context"
	"fmt"
	"github.com/anthhub/forwarder"
	"github.com/borealisdb/cli/pkg/config"
	"github.com/borealisdb/commons/constants"
	"github.com/borealisdb/commons/k8sutil"
	"github.com/borealisdb/go-sdk/api"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	DefaultNamespace = "default"
)

var (
	datname    string
	certPath   string
	username   string
	kubeconfig string
	role       string
)

// clusterConnectCmd represents the clusterConnect command
var clusterConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect to a cluster's database",
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfig()
		client, err := k8sutil.NewFromConfig(config)
		if err != nil {
			log.Fatal(err)
		}

		// Download certificates
		if err := downloadCerts(client); err != nil {
			cobra.CheckErr(err)
		}

		// Get token
		token, err := sdk.GenerateClusterToken(api.GenerateClusterTokenRequest{
			ClusterName: clusterName,
		})
		if err != nil {
			cobra.CheckErr(err)
		}

		endpoint, err := getEndpoint(client)
		if err != nil {
			cobra.CheckErr(err)
		}

		options := []*forwarder.Option{
			{
				LocalPort:  5432,
				RemotePort: 5432,
				Namespace:  namespace,
				Source:     endpoint,
			},
		}
		forwarders, err := forwarder.WithForwarders(context.Background(), options, kubeconfig)
		if err != nil {
			cobra.CheckErr(err)
		}

		defer forwarders.Close()

		ready, err := forwarders.Ready()
		if err != nil {
			cobra.CheckErr(err)
		}
		log.Infof("port forward is ready: %+v\n", ready)

		command := exec.Command(
			"psql",
			fmt.Sprintf("sslmode=verify-ca sslrootcert=%v host=localhost user=%v port=5432 dbname=%v", certPath, username, datname),
		)
		command.Env = os.Environ()
		command.Env = append(command.Env, fmt.Sprintf("PGPASSWORD=%v", token.AccessToken))
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		if err := command.Run(); err != nil {
			cobra.CheckErr(err)
		}
	},
}

func getConfig() *restclient.Config {
	envKube := os.Getenv("KUBECONFIG")
	if envKube != "" {
		kubeconfig = envKube
	}

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	return kubeConfig
}

func downloadCerts(client k8sutil.KubernetesClient) error {
	tlsSecret, err := client.Secrets(namespace).Get(context.Background(), constants.GetTLSSecretName(clusterName), metav1.GetOptions{})
	if err != nil {
		return err
	}

	if err := os.WriteFile(
		certPath,
		tlsSecret.Data[constants.RootCaCertName],
		0777,
	); err != nil {
		return err
	}

	return nil
}

func getEndpoint(client k8sutil.KubernetesClient) (string, error) {
	if role == "replica" {
		return fmt.Sprintf("svc/%v", clusterName), nil
	}

	pods, err := client.
		Pods(namespace).
		List(context.Background(), metav1.ListOptions{LabelSelector: fmt.Sprintf("spilo-role=%v,cluster-name=%v", role, clusterName)})
	if err != nil {
		return "", err
	}

	if len(pods.Items) == 0 {
		return "", fmt.Errorf("could not find any pods for %v for role %v", clusterName, role)
	}

	return fmt.Sprintf("po/%v", pods.Items[0].Name), nil
}

func init() {
	clusterCmd.AddCommand(clusterConnectCmd)
	clusterConnectCmd.PersistentFlags().StringVar(&clusterName, "cluster-name", "", "")
	clusterConnectCmd.PersistentFlags().StringVar(&namespace, "namespace", DefaultNamespace, "")
	clusterConnectCmd.PersistentFlags().StringVar(&datname, "db", "postgres", "")
	clusterConnectCmd.PersistentFlags().StringVar(&username, "username", "", "")
	clusterConnectCmd.PersistentFlags().StringVar(&certPath, "cert-path", filepath.Join(config.CliConfigDefaultPath, constants.RootCaCertName), "")
	clusterConnectCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "")
	clusterConnectCmd.PersistentFlags().StringVar(&role, "role", "master", "")

	clusterConnectCmd.MarkPersistentFlagRequired("cluster-name")
	clusterConnectCmd.MarkPersistentFlagRequired("username")
}

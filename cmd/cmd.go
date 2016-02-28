package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	rootLongDescription = "azkube is a kubernetes deployment helper for Azure"
)

var rootArgs util.RootArguments

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "azkube",
		Short: rootLongDescription,
		Long:  rootLongDescription,
	}

	rootCmd.PersistentFlags().StringVar(&rootArgs.TenantID, "tenant-id", "", "azure subscription id")
	rootCmd.PersistentFlags().StringVar(&rootArgs.SubscriptionID, "subscription-id", "", "azure tenant id")
	rootCmd.PersistentFlags().StringVar(&rootArgs.AuthMethod, "auth-method", "device", "auth method (default:`device`, `clientsecret`, `clientcertificate`)")
	rootCmd.PersistentFlags().StringVar(&rootArgs.ClientID, "client-id", "", "client id (used with --auth-method=[clientsecret|clientcertificate])")
	rootCmd.PersistentFlags().StringVar(&rootArgs.ClientSecret, "client-secret", "", "client secret (used with --auth-mode=clientsecret)")
	rootCmd.PersistentFlags().StringVar(&rootArgs.ClientCertificatePath, "certificate-path", "", "path to client certificate (used with --auth-method=clientcertificate)")
	rootCmd.PersistentFlags().StringVar(&rootArgs.PrivateKeyPath, "private-key-path", "", "path to private key (used with --auth-method=clientcertificate)")

	viper.SetEnvPrefix("azkube_")
	viper.BindEnv("tenant_id")
	viper.BindEnv("subscription_id")
	viper.BindEnv("client_secret")
	viper.BindEnv("client_certificate_path")
	viper.BindPFlag("tenant_id", rootCmd.PersistentFlags().Lookup("tenant-id"))
	viper.BindPFlag("subscription_id", rootCmd.Flags().Lookup("subscription-id"))
	viper.BindPFlag("client_id", rootCmd.PersistentFlags().Lookup("client-id"))
	viper.BindPFlag("client_secret", rootCmd.PersistentFlags().Lookup("client-secret"))
	viper.BindPFlag("client_certificate_path", rootCmd.PersistentFlags().Lookup("client-certificate-path"))

	rootCmd.AddCommand(NewDeployCmd())
	//rootCmd.AddCommand(NewScaleDeploymentCmd())
	//rootCmd.AddCommand(NewDestroyDeploymentCmd())

	return rootCmd
}

func validateRootArgs(rootArgs util.RootArguments) {
	if rootArgs.SubscriptionID == "" {
		log.Fatal("--subscription-id must be specified")
	}

	if rootArgs.TenantID == "" {
		log.Fatal("--tenant-id must be specified")
	}

	if rootArgs.AuthMethod == "client-secret" {
		if rootArgs.ClientID == "" || rootArgs.ClientSecret == "" {
			log.Fatal("--client-id and --client-secret must be specified when --auth-method=\"client_secret\"")
		}
	} else if rootArgs.AuthMethod == "client-certificate" {
		if rootArgs.ClientID == "" || rootArgs.ClientCertificatePath == "" {
			log.Fatal("--client-id and --client-certificate must be specified when --auth-method=\"client_certificate\"")
		}

	}
}

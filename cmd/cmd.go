package cmd

import (
	"flag"

	"github.com/colemickens/azkube/util"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
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

	rootCmd.AddCommand(NewDeployCmd())
	//rootCmd.AddCommand(NewScaleDeploymentCmd())
	//rootCmd.AddCommand(NewDestroyDeploymentCmd())

	return rootCmd
}

func validateRootArgs(rootArgs util.RootArguments) {
	flag.Parse() // make glog happy (why couldn't they just have an init?)
	_ = glog.Info
	// noop for now
	// validate auth_type, etc
}

package cmd

import (
	"io/ioutil"
	"log"
	"path/filepath"

	// "github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	// "github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/azure"
	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	certInstallLongDescription = "long desc"
)

func NewCertInstallCmd() *cobra.Command {
	var configPath string
	var machineType string
	var destination string

	var certInstallCmd = &cobra.Command{
		Use:   "certinstall",
		Short: "install certificates on the machine",
		Long:  certInstallLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting certinstall command")

			var deployProperties util.DeploymentProperties
			RunCertInstallCmd(deployProperties, machineType, destination)

			log.Println("finished certinstall command")
		},
	}

	certInstallCmd.Flags().StringVarP(&configPath, "config", "c", "/etc/kubernetes/azure.json", "path to config")
	certInstallCmd.Flags().StringVarP(&machineType, "machineType", "m", "", "machine type: 'master' or 'node'")
	certInstallCmd.Flags().StringVarP(&machineType, "destination", "d", "/etc/kubernetes/", "machine type: 'master' or 'node'")

	return certInstallCmd
}

func RunCertInstallCmd(deployProperties util.DeploymentProperties, machineType, destination string) {
	var err error

	d, err := util.NewDeployerWithCertificate("a", "b", "c", "d", "e")
	if err != nil {
		panic(err)
	}

	secretMap := map[string]map[string]string{
		"master": {
			"ca-crt":                               "ca.crt",
			"apiserver-crt":                        "apiserver.crt",
			"apiserver-key":                        "apiserver.key",
			"master-proxy-kubeconfig":              "master-proxy-kubeconfig",
			"master-kubelet-kubeconfig":            "master-kubelet-kubeconfig",
			"master-scheduler-kubeconfig":          "master-scheduler-kubeconfig",
			"master-controller-manager-kubeconfig": "master-controller-manager-kubeconfig",
		},
		"node": {
			"node-proxy-kubeconfig":   "node-proxy-kubeconfig",
			"node-kubelet-kubeconfig": "node-kubelet-kubeconfig",
		},
		"etcd": {},
	}

	log.Println("bootstrapping secrets for", machineType)
	secrets, ok := secretMap[machineType]
	if !ok {
		log.Fatalln("don't have a secret list for", machineType)
	}

	for secretName, secretPath := range secrets {
		log.Println("retrieving secret:", secretName)

		// TODO (colemickens): fix this!shit
		secretValue, err := d.GetSecret(d.Config.VaultName, secretName)
		if err != nil {
			// TODO(colemickens): retry?
			panic(err)
		}

		secretDestinationPath := filepath.Join(destination, secretPath)
		log.Println("writing secret:", secretDestinationPath)
		err = ioutil.WriteFile(secretDestinationPath, []byte(*secretValue), 0644)
		if err != nil {
			// TODO(colemickens): retry?
			panic(err)
		}
	}

	log.Println("done")
}

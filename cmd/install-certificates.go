package cmd

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"

	// "github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	// "github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/azure"
	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	installCertificatesLongDescription = "long desc"
)

func NewInstallCertificatesCmd() *cobra.Command {
	var statePath string
	var machineType string
	var destination string

	var installCertificatesCmd = &cobra.Command{
		Use:   "install-certificates",
		Short: "install certificates on the machine",
		Long:  installCertificatesLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting install-certificates command")

			var state *util.State
			var err error
			state, err = ReadAndValidateState(
				statePath,
				[]reflect.Type{
					reflect.TypeOf(state.Common),
					reflect.TypeOf(state.App),
					reflect.TypeOf(state.Ssh),
					reflect.TypeOf(state.Pki),
					reflect.TypeOf(state.Vault),
					reflect.TypeOf(state.Secrets),
					reflect.TypeOf(state.Myriad),
				},
				[]reflect.Type{},
			)
			if err != nil {
				panic(err)
			}

			state = RunInstallCertificatesCmd(state, machineType, destination)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished install-certificates command")
		},
	}

	installCertificatesCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")
	installCertificatesCmd.Flags().StringVarP(&machineType, "machineType", "m", "", "machine type: 'master' or 'node'")
	installCertificatesCmd.Flags().StringVarP(&machineType, "destination", "d", "/etc/kubernetes/", "machine type: 'master' or 'node'")

	return installCertificatesCmd
}

func RunInstallCertificatesCmd(stateIn *util.State, machineType, destination string) (stateOut *util.State) {
	*stateOut = *stateIn

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
		secretValue, err := d.VaultClient.GetSecret(stateIn.Vault.Name, secretName)
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

	return stateOut
}

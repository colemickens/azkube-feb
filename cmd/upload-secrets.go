package cmd

import (
	"io/ioutil"
	"log"
	"reflect"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	uploadSecretsLongDescription = "long desc"
)

func NewUploadSecretsCmd() *cobra.Command {
	var statePath string

	var uploadSecretsCmd = &cobra.Command{
		Use:   "upload-secrets",
		Short: "upload secrets for kubernetes deployment",
		Long:  uploadSecretsLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting upload-secrets command")

			state := &util.State{}
			var err error
			state, err = ReadAndValidateState(
				statePath,
				[]reflect.Type{
					reflect.TypeOf(state.Common),
					reflect.TypeOf(state.App),
					reflect.TypeOf(state.Pki),
					reflect.TypeOf(state.Vault),
				},
				[]reflect.Type{
					reflect.TypeOf(state.Secrets),
					reflect.TypeOf(state.Myriad),
				},
			)
			if err != nil {
				panic(err)
			}

			RunUploadSecretsCmd(state)
			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished upload-secrets command")
		},
	}

	uploadSecretsCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return uploadSecretsCmd
}

func RunUploadSecretsCmd(state *util.State) {
	d, err := util.NewDeployerFromState(*state)
	if err != nil {
		panic(err)
	}

	// TODO: factor this out so it can be shared with install-certificates
	secrets := map[string]string{
		"pki/ca.crt":                               "ca-crt",
		"pki/apiserver.crt":                        "apiserver-crt",
		"pki/apiserver.key":                        "apiserver-key",
		"pki/node-proxy-kubeconfig":                "node-proxy-kubeconfig",
		"pki/node-kubelet-kubeconfig":              "node-kubelet-kubeconfig",
		"pki/master-proxy-kubeconfig":              "master-proxy-kubeconfig",
		"pki/master-kubelet-kubeconfig":            "master-kubelet-kubeconfig",
		"pki/master-scheduler-kubeconfig":          "master-scheduler-kubeconfig",
		"pki/master-controller-manager-kubeconfig": "master-controller-manager-kubeconfig",
	}

	// upload the PFX in the special format that KV requires for this secret
	pfxBytes, err := state.App.ServicePrincipalPkcs12()
	if err != nil {
		panic(err)
	}
	servicePrincipalSecretURL, err := d.VaultClient.PutKeyVaultCertificate(
		state.Vault.Name,
		"servicePrincipal-pfx",
		pfxBytes,
	)
	if err != nil {
		panic(err)
	}

	for secretPath, secretName := range secrets {
		secretValue, err := ioutil.ReadFile(secretPath)
		if err != nil {
			panic(err)
		}
		_, err = d.VaultClient.PutSecret(state.Vault.Name, secretName, secretValue)
		if err != nil {
			panic(err)
		}
	}

	state.Secrets = &util.SecretsProperties{
		ServicePrincipalSecretURL: servicePrincipalSecretURL,
	}
}

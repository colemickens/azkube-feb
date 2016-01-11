package cmd

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	deployVaultLongDescription = "long desc"
)

func NewDeployVaultCmd() *cobra.Command {
	var statePath string
	var vaultName string
	var waitDns bool

	var deployVaultCmd = &cobra.Command{
		Use:   "deploy-vault",
		Short: "deploy the Azure KeyVault",
		Long:  deployVaultLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting deploy-vault command")

			state := &util.State{}
			var err error
			state, err = ReadAndValidateState(statePath,
				[]reflect.Type{
					reflect.TypeOf(state.Common),
					reflect.TypeOf(state.App),
				},
				[]reflect.Type{
					reflect.TypeOf(state.Vault),
					reflect.TypeOf(state.Secrets),
					reflect.TypeOf(state.Myriad),
				},
			)
			if err != nil {
				panic(err)
			}

			if vaultName == "" {
				vaultName = state.Common.DeploymentName + "-vault"
				vaultName = strings.Replace(vaultName, "-", "", -1)
			}

			RunDeployVaultCmd(state, vaultName, waitDns)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished deploy-vault command")
		},
	}

	deployVaultCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")
	deployVaultCmd.Flags().StringVarP(&vaultName, "vault-name", "v", "", "vault name (will be derived from deployment name if empty")
	deployVaultCmd.Flags().BoolVarP(&waitDns, "wait-dns", "w", true, "wait for vault dns to resolve before considering step complete")

	return deployVaultCmd
}

// TODO: should these get a copy of state and return just their subelement?
func RunDeployVaultCmd(state *util.State, vaultName string, waitDns bool) {
	d, err := util.NewDeployerFromState(*state)
	if err != nil {
		panic(err)
	}

	vaultTemplateInput := util.VaultTemplateInput{
		VaultName:                vaultName,
		TenantID:                 state.Common.TenantID,
		ServicePrincipalObjectID: state.App.ServicePrincipalObjectID,
	}

	vaultTemplate, err := util.PopulateTemplate(util.VaultTemplate, vaultTemplateInput)
	if err != nil {
		panic(err)
	}

	_, err = d.DoDeployment(*state.Common, "vault", vaultTemplate, true)

	state.Vault = &util.VaultProperties{
		Name: vaultName,
	}

	// looop until dns resolves
	if waitDns {
		for {
			vaultHostname := fmt.Sprintf("%s.vault.azure.net", vaultName)
			log.Println("waiting for vault dns to resolve: ", vaultHostname)

			ip, err := net.ResolveIPAddr("ip4", vaultHostname)
			if err != nil {
				log.Println("resolve err:", err)
			} else if ip != nil {
				break
			}

			time.Sleep(time.Second * 30)
		}
	}
}

package cmd

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"time"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	deployVaultLongDescription = "long desc"
)

func NewDeployVaultCmd() *cobra.Command {
	var statePath string

	var deployVaultCmd = &cobra.Command{
		Use:   "deploy-vault",
		Short: "deploy the Azure KeyVault",
		Long:  deployVaultLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting deploy-vault command")

			state = ReadAndValidate(statePath,
				[]reflect.Type{
					reflect.TypeOf(state.CommonProperties),
					reflect.TypeOf(state.AppProperties),
					reflect.TypeOf(state.SshProperties),
					reflect.TypeOf(state.PkiProperties),
					reflect.TypeOf(state.VaultProperites),
					reflect.TypeOf(state.SecretsProperties),
					reflect.TypeOf(state.MyriadProperties),
				},
				[]reflect.Type{},
			)

			state, err = RunDeployVaultCmd(state)
			if err != nil {
				panic(err)
			}

			state = RunValidateDeploymentCmd(stateIn)
			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished deploy-vault command")
		},
	}

	destroyVaultCmd.Flags().StringVar(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return deployVaultCmd
}

func RunDeployVaultCmd(stateIn util.State) (stateOut util.State) {
	stateOut = stateIn
}

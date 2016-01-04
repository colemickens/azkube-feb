package cmd

import (
	"log"
	"reflect"

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

			var state *util.State
			var err error
			state, err = ReadAndValidateState(statePath,
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

			state = RunDeployVaultCmd(state)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished deploy-vault command")
		},
	}

	deployVaultCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return deployVaultCmd
}

func RunDeployVaultCmd(stateIn *util.State) (stateOut *util.State) {
	*stateOut = *stateIn

	return stateOut
}

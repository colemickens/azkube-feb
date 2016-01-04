package cmd

import (
	"log"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	validateDeploymentLongDescription = "long desc"
)

func NewValidateDeploymentCmd() *cobra.Command {
	var statePath string

	var validateDeploymentCmd = &cobra.Command{
		Use:   "validate-deployment",
		Short: "validate the completed deployment",
		Long:  validateDeploymentLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting validate-deployment command")

			state, err = ReadAndValidateState(
				statePath,
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
			if err != nil {
				panic(err)
			}

			state = RunValidateDeploymentCmd(stateIn)
			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished validate-deployment command")
		},
	}

	validateDeploymentCmd.Flags().StringVar(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return validateDeploymentCmd
}

func RunValidateDeploymentCmd(stateIn util.State) (stateOut util.State) {
	stateOut = stateIn
}

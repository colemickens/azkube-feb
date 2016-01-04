package cmd

import (
	"log"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	destroyDeploymentLongDescription = "long desc"
)

func NewDestroyDeploymentCmd() *cobra.Command {
	var statePath string

	var destroyDeploymentCmd = &cobra.Command{
		Use:   "destroy-deployment",
		Short: "destroy a kubernetes deployment",
		Long:  destroyDeploymentLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting destroy-deployment command")

			state, err = ReadAndValidateState(
				statePath,
				[]reflect.Type{
					reflect.TypeOf(state.CommonProperties),
					reflect.TypeOf(state.AppProperties),
					reflect.TypeOf(state.SshProperties),
					reflect.TypeOf(state.PkiProperties),
					reflect.TypeOf(state.VaultProperites),
					reflect.TypeOf(state.SecretsProperties),
				},
				[]reflect.Type{
					reflect.TypeOf(state.MyriadProperties),
				},
			)
			if err != nil {
				panic(err)
			}

			state = RunScaleDeploymentCmd(stateIn)
			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished destroy-deployment command")
		},
	}

	uploadSecretsCmd.Flags().StringVar(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return destroyDeploymentCmd
}

func RunScaleDeploymentCmd(stateIn util.State) (stateOut util.State) {

}

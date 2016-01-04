package cmd

import (
	"log"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	deployMyriadLongDescription = "long desc"
)

func NewDeployMyriadCmd() *cobra.Command {
	var statePath string
	var deploymentName string

	var deployMyriadCmd = &cobra.Command{
		Use:   "deploy-myriad",
		Short: "deploy the coreos kubernetes machines",
		Long:  deployMyriadLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting deploy-myriad command")

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

			state, err = RunDeployMyriadCmd(state)
			if err != nil {
				panic(err)
			}

			state = RunValidateDeploymentCmd(stateIn)
			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished deploy-myriad command")
		},
	}

	deployMyriadCmd.Flags().StringVar(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return deployMyriadCmd
}

func RunDeployMyriadCmd(stateIn util.State) (stateOut util.State) {
	stateOut = stateIn
}

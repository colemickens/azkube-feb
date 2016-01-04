package cmd

import (
	"log"
	"reflect"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	deployMyriadLongDescription = "long desc"
)

func NewDeployMyriadCmd() *cobra.Command {
	var statePath string

	var deployMyriadCmd = &cobra.Command{
		Use:   "deploy-myriad",
		Short: "deploy the coreos kubernetes machines",
		Long:  deployMyriadLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting deploy-myriad command")

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
				},
				[]reflect.Type{
					reflect.TypeOf(state.Myriad),
				},
			)

			state = RunDeployMyriadCmd(state)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished deploy-myriad command")
		},
	}

	deployMyriadCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return deployMyriadCmd
}

func RunDeployMyriadCmd(stateIn *util.State) (stateOut *util.State) {
	*stateOut = *stateIn

	return stateOut
}

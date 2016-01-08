package cmd

import (
	"log"
	"reflect"

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

			state := &util.State{}
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

			state = RunValidateDeploymentCmd(state)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished validate-deployment command")
		},
	}

	validateDeploymentCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return validateDeploymentCmd
}

func RunValidateDeploymentCmd(stateIn *util.State) (stateOut *util.State) {
	*stateOut = *stateIn

	return stateOut
}

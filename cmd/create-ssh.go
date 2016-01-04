package cmd

import (
	"log"
	"reflect"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	createSshLongDescription = "long desc"
)

func NewCreateSshCmd() *cobra.Command {
	var statePath string

	var createSshCmd = &cobra.Command{
		Use:   "create-ssh",
		Short: "creates ssh configuration needed by subsequent stages",
		Long:  createSshLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting create-ssh command")

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
			if err != nil {
				panic(err)
			}

			state = RunCreateSshCmd(state)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished create-ssh command")
		},
	}

	createSshCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return createSshCmd
}

func RunCreateSshCmd(stateIn *util.State) (stateOut *util.State) {
	*stateOut = *stateIn

	return stateOut
}

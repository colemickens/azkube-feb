package cmd

import (
	"log"
	"reflect"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	createPkiLongDescription = "long desc"
)

func NewCreatePkiCmd() *cobra.Command {
	var statePath string

	var createPkiCmd = &cobra.Command{
		Use:   "create-pki",
		Short: "creates the public key infrastructure to be used by the Kubernetes cluster",
		Long:  createPkiLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting create-pki command")

			var state *util.State
			var err error
			state, err = ReadAndValidateState(statePath,
				[]reflect.Type{},
				[]reflect.Type{},
			)
			if err != nil {
				panic(err)
			}

			RunCreatePkiCmd(state)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished create-pki command")
		},
	}

	createPkiCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return createPkiCmd
}

func RunCreatePkiCmd(state *util.State) {
	var err error
	state.Pki, err = util.GeneratePki(*state.Common)
	if err != nil {
		panic(err)
	}
}

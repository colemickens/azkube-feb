package cmd

import (
	"log"

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

			state, err = ReadState(statePath)
			if err != nil {
				panic(err)
			}

			state = RunCreatePkiCmd(deployProperties, state)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished create-pki command")
		},
	}

	createPkiCmd.Flags().StringVar(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return createPkiCmd
}

func RunCreatePkiCmd(stateIn util.State, deploymentName string) (stateOut util.State) {
	stateOut = stateIn
}

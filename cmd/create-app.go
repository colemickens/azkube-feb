package cmd

import (
	"log"
	"reflect"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	createAppLongDescription = "long desc"
)

func NewCreateAppCmd() *cobra.Command {
	var statePath string
	var appName string
	var appIdentifierURL string

	var createAppCmd = &cobra.Command{
		Use:   "create-app",
		Short: "creates active directory application to be used by Kubernetes itself",
		Long:  createAppLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting create-config command")

			var state *util.State
			var err error
			state, err = ReadAndValidateState(statePath,
				[]reflect.Type{
					reflect.TypeOf(state.Common),
				},
				[]reflect.Type{
					reflect.TypeOf(state.App),
					reflect.TypeOf(state.Ssh),
					reflect.TypeOf(state.Pki),
					reflect.TypeOf(state.Vault),
					reflect.TypeOf(state.Secrets),
					reflect.TypeOf(state.Myriad),
				},
			)
			if err != nil {
				panic(err)
			}

			if appName == "" {
				appName = state.Common.DeploymentName
			}

			if appIdentifierURL == "" {
				appIdentifierURL = "http://" + state.Common.DeploymentName + "/"
			}

			state, err = RunCreateAppCmd(state, appName, appIdentifierURL)
			if err != nil {
				panic(err)
			}

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished create-app command")
		},
	}

	createAppCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")
	createAppCmd.Flags().StringVarP(&appName, "name", "n", "", "name of the app")
	createAppCmd.Flags().StringVarP(&appIdentifierURL, "identifier-url", "i", "", "identifier-url for the app")

	return createAppCmd
}

func RunCreateAppCmd(stateIn *util.State, appName, appURL string) (stateOut *util.State, err error) {
	*stateOut = *stateIn // copy inputs

	stateOut.App = &util.AppProperties{
	// make copy of inputs or something for now?
	// hard code these?
	// make it load from a file and pretend it actually went out and made it ?
	}

	panic("you must do this yourself for now")
}

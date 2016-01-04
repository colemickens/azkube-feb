package cmd

import (
	//"encoding/hex"
	//"io/ioutil"
	//"log"
	"time"

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

			state = ReadAndValidate(statePath,
				[]reflect.Type{
					reflect.TypeOf(state.CommonProperties),
				},
				[]reflect.Type{
					reflect.TypeOf(state.AppProperties),
					reflect.TypeOf(state.SshProperties),
					reflect.TypeOf(state.PkiProperties),
					reflect.TypeOf(state.VaultProperites),
					reflect.TypeOf(state.SecretsProperties),
					reflect.TypeOf(state.MyriadProperties),
				},
			)

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

	createAppCmd.Flags().StringVar(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")
	createAppCmd.Flags().StringVar(&appName, "name", "n", "", "name of the app")
	createAppCmd.Flags().StringVar(&appIdentifierURL, "identifier-url", "i", "", "identifier-url for the app")

	return createAppCmd
}

func RunCreateAppCmd(stateIn util.State, appName, appURL string) (stateOut util.State, err error) {
	stateOut.App = util.AppProperties{
	// make copy of inputs or something for now?
	// hard code these?
	// make it load from a file and pretend it actually went out and made it ?
	}
	panic("you must do this yourself for now")
}
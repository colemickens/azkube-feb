package cmd

import (
	"log"
	"reflect"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	uploadSecretsLongDescription = "long desc"
)

func NewUploadSecretsCmd() *cobra.Command {
	var statePath string

	var uploadSecretsCmd = &cobra.Command{
		Use:   "upload-secrets",
		Short: "upload secrets for kubernetes deployment",
		Long:  uploadSecretsLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting upload-secrets command")

			var state *util.State
			var err error
			state, err = ReadAndValidateState(
				statePath,
				[]reflect.Type{
					reflect.TypeOf(state.Common),
					reflect.TypeOf(state.App),
					reflect.TypeOf(state.Ssh),
					reflect.TypeOf(state.Pki),
					reflect.TypeOf(state.Vault),
				},
				[]reflect.Type{
					reflect.TypeOf(state.Secrets),
					reflect.TypeOf(state.Myriad),
				},
			)
			if err != nil {
				panic(err)
			}

			state = RunUploadSecretsCmd(state)
			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished upload-secrets command")
		},
	}

	uploadSecretsCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return uploadSecretsCmd
}

func RunUploadSecretsCmd(stateIn *util.State) (stateOut *util.State) {
	*stateOut = *stateIn

	return stateOut
}

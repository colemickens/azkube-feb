package cmd

import (
	"log"

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

			state, err = ReadAndValidateState(
				statePath,
				[]reflect.Type{
					reflect.TypeOf(state.CommonProperties),
					reflect.TypeOf(state.AppProperties),
					reflect.TypeOf(state.SshProperties),
					reflect.TypeOf(state.PkiProperties),
					reflect.TypeOf(state.VaultProperites),
				},
				[]reflect.Type{
					reflect.TypeOf(state.SecretsProperties),
					reflect.TypeOf(state.MyriadProperties),
				},
			)
			if err != nil {
				panic(err)
			}

			state = RunUploadSecretsCmd(stateIn)
			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished upload-secrets command")
		},
	}

	uploadSecretsCmd.Flags().StringVar(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return uploadSecretsCmd
}

func RunUploadSecretsCmd(stateIn util.State) (stateOut util.State) {
	stateOut = stateIn
}

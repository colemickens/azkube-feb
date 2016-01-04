package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	rootLongDescription = "longer description"
)

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "azkube",
		Short: "azure <-> kubernetes tool",
		Long:  rootLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// on deployer's box
	rootCmd.AddCommand(NewCreateCommonCmd())
	rootCmd.AddCommand(NewCreateAppCmd())
	rootCmd.AddCommand(NewCreateSshCmd())
	rootCmd.AddCommand(NewCreatePkiCmd())
	rootCmd.AddCommand(NewDeployVaultCmd())
	rootCmd.AddCommand(NewUploadSecretsCmd())
	rootCmd.AddCommand(NewDeployMyriadCmd())

	// on deployer's box (after cluster is up)
	rootCmd.AddCommand(NewScaleDeploymentCmd())
	rootCmd.AddCommand(NewDestroyDeploymentCmd())

	// on cluster's box
	rootCmd.AddCommand(NewInstallCertificatesCmd())

	return rootCmd
}

func ReadAndValidateState(path string, expected, forbidden []reflect.Type) (stateOut *util.State, err error) {
	var state util.State
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(contents, &state)
	if err != nil {
		return nil, err
	}

	stateValue := reflect.ValueOf(state)
	for i := 0; i < stateValue.NumField(); i++ {
		field := stateValue.FieldByIndex([]int{i})

		if field.Kind() != reflect.Ptr {
			continue
		}

		if !field.Elem().IsValid() {
			for _, expect := range expected {
				if field.Type() == expect {
					return nil, fmt.Errorf("field %s was nil and should not be", field.Type().Name())
				}
			}
		} else {
			for _, forbid := range forbidden {
				if field.Type() == forbid {
					return nil, fmt.Errorf("field %s was filled but should not be", field.Type().Name())
				}
			}
		}
	}
	return &state, nil
}

func WriteState(path string, state *util.State) (err error) {
	stateJson, err := json.Marshal(state)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, stateJson, 0777) // TODO: review file perms
	if err != nil {
		return err
	}

	return nil
}

package cmd

import (
	"reflect"

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
	rootCmd.AddCommand(NewCreatePkiCmd())
	rootCmd.AddCommand(NewCreateSshCmd())
	rootCmd.AddCommand(NewCreateAppCmd())
	rootCmd.AddCommand(NewDeployVaultCmd())
	rootCmd.AddCommand(NewUploadSecretsCmd())
	rootCmd.AddCommand(NewDeployMyriadCmd())

	// on deployer's box (after cluster is up)
	rootCmd.AddCommand(NewScaleCmd())
	rootCmd.AddCommand(NewDestroyCmd())

	// on cluster's box
	rootCmd.AddCommand(NewInstallCertificatesCmd())

	return rootCmd
}

func ReadAndValidateState(path string, expected, forbidden []reflect.Type) (state util.State, err error) {
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
	for i = 0; i < stateValue.NumField(); i++ {
		field := stateValue.FieldByIndex(i)

		if field.Kind() != reflect.Ptr {
			continue
		}

		if !field.Elem().IsValid() {
			for _, expect := range expected {
				if field.Type() == expect {
					return fmt.Errorf("field %s was nil and should not be", field.Name)
				}
			}
		} else {
			for _, forbid := range forbidden {
				if field.Type() == forbid {
					return fmt.Errorf("field %s was filled but should not be", field.Name)
				}
			}
		}
	}
	return state, nil
}

func WriteState(path string, state util.State) (err error) {
	stateJson, err = json.Marshal(state)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, stateJson, 0777) // TODO: review file perms
	if err != nil {
		return err
	}

	return nil
}

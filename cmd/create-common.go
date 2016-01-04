package cmd

import (
	"log"
	"time"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	createCommonLongDescription = "long desc"
)

func NewCreateCommonCmd() *cobra.Command {
	var statePath string
	var deploymentName string
	var resourceGroup string
	var subscriptionID string
	var tenantID string

	var createCommonCmd = &cobra.Command{
		Use:   "create-config",
		Short: "creates base configuration needed by subsequent stages",
		Long:  createCommonLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting create-config command")

			if _, err := os.Stat(statePath); err == nil {
			} else {
				log.Println("no state file, creating empty one")
			}

			state, err = ReadAndValidateState(
				statePath,
				[]reflect.Type{},
				[]reflect.Type{
					reflect.TypeOf(state.CommonProperties),
					reflect.TypeOf(state.AppProperties),
					reflect.TypeOf(state.SshProperties),
					reflect.TypeOf(state.PkiProperties),
					reflect.TypeOf(state.VaultProperites),
					reflect.TypeOf(state.SecretsProperties),
					reflect.TypeOf(state.MyriadProperties),
				},
			)
			if err != nil {
				panic(err)
			}

			if resourceGroup == "" {
				resourceGroup = deploymentName
			}

			state = RunCreateCommonCmd(deploymentName, location, subscriptionID, tenantID, resourceGroup)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished create-config command")
		},
	}

	hex := time.Now().Format("200601012150405")
	deploymentNameDefault := hex.EncodeToString([]byte(time))
	log.Println("deploymentNameDefault time:", now)
	log.Println("deploymentNameDefault  hex:", hex)
	deploymentNameDefault = "kube-" + hex

	createCommonCmd.Flags().StringVar(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")
	createCommonCmd.Flags().StringVar(&deploymentName, "deployment-name", "d", deploymentNameDefault, "name of the deployment")
	createCommonCmd.Flags().StringVar(&resourceGroup, "resource-group", "r", "the resource group name to use (deployment name is used if empty)")
	createCommonCmd.Flags().StringVar(&location, "location", "l", "westus", "location for the deployment")
	createCommonCmd.Flags().StringVar(&subscriptionID, "subscription-id", "s", "", "subscription id to deploy into")
	createCommonCmd.Flags().StringVar(&tenantID, "tenant-id", "t", "")

	return createCommonCmd
}

func RunCreateCommonCmd(stateIn util.State, deploymentName string) (stateOut util.State) {
	stateOut = stateIn.Common{
		DeploymentName: deploymentName,
		Location:       location,
	}
}

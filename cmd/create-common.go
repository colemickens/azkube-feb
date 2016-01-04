package cmd

import (
	"encoding/hex"
	"log"
	"os"
	"reflect"
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
	var location string
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

			var state *util.State
			var err error
			state, err = ReadAndValidateState(
				statePath,
				[]reflect.Type{},
				[]reflect.Type{
					reflect.TypeOf(state.Common),
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

			if deploymentName == "" {
				now := time.Now().Format("200601012150405")
				nowHex := hex.EncodeToString([]byte(now))
				log.Println("deploymentNameDefault time:", now)
				log.Println("deploymentNameDefault  hex:", nowHex)
				deploymentName = "kube-" + nowHex
			}

			if resourceGroup == "" {
				resourceGroup = deploymentName
			}

			state = RunCreateCommonCmd(state, deploymentName, location, subscriptionID, tenantID, resourceGroup)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished create-config command")
		},
	}

	// TODO(colemickens) stringvarp -> stringvar

	createCommonCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")
	createCommonCmd.Flags().StringVarP(&deploymentName, "deployment-name", "d", "", "name of the deployment (one will be auto-generated if empty)")
	createCommonCmd.Flags().StringVarP(&resourceGroup, "resource-group", "r", "", "the resource group name to use (deployment name is used if empty)")
	createCommonCmd.Flags().StringVarP(&location, "location", "l", "westus", "location for the deployment")
	createCommonCmd.Flags().StringVarP(&subscriptionID, "subscription-id", "s", "", "subscription id to deploy into")
	createCommonCmd.Flags().StringVarP(&tenantID, "tenant-id", "t", "", "tenant id of account")

	return createCommonCmd
}

func RunCreateCommonCmd(stateIn *util.State, deploymentName, location, subscriptionID, tenantID, resourceGroup string) (stateOut *util.State) {
	*stateOut = *stateIn

	stateOut.Common = &util.CommonProperties{
		DeploymentName: deploymentName,
		Location:       location,
		SubscriptionID: subscriptionID,
		TenantID:       tenantID,
		ResourceGroup:  resourceGroup,
	}

	return stateOut
}

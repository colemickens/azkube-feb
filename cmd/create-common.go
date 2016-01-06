package cmd

import (
	"fmt"
	"io/ioutil"
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
	var masterFQDN string

	var createCommonCmd = &cobra.Command{
		Use:   "create-common",
		Short: "creates common configuration needed by subsequent stages",
		Long:  createCommonLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting create-config command")

			if _, err := os.Stat(statePath); err == nil {
			} else {
				ioutil.WriteFile(statePath, []byte("{}"), os.ModePerm)
				log.Println("no state file, creating empty one")
			}

			state := &util.State{}
			state, err := ReadAndValidateState(
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
				deploymentName = fmt.Sprintf("kube-%s", time.Now().Format("20060102-150405"))
			}

			if resourceGroup == "" {
				resourceGroup = deploymentName
			}

			// TODO(colemickens): validate location

			if masterFQDN == "" {
				masterHostname := fmt.Sprintf("%s-master", deploymentName)
				masterFQDN = fmt.Sprintf("https://%s.%s.cloudapp.azure.net", masterHostname, location)
			}

			if subscriptionID == "" || tenantID == "" {
				panic("subscriptionId and tenantId must be specified!")
			}

			state = RunCreateCommonCmd(state, deploymentName, resourceGroup, location, subscriptionID, tenantID, masterFQDN)

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
	createCommonCmd.Flags().StringVarP(&subscriptionID, "subscription-id", "i", "", "subscription id to deploy into")
	createCommonCmd.Flags().StringVarP(&tenantID, "tenant-id", "t", "", "tenant id of account")
	createCommonCmd.Flags().StringVarP(&masterFQDN, "master-fqdn", "m", "", "master fqdn (will use automatic cloudapp hostname otherwise) (NOTE: THIS AFFECTS CERT SUBJECT NAME)")

	return createCommonCmd
}

func RunCreateCommonCmd(stateIn *util.State, deploymentName, resourceGroup, location, subscriptionID, tenantID, masterFQDN string) (stateOut *util.State) {
	log.Println("made it here")

	stateOut = &*stateIn

	stateOut.Common = &util.CommonProperties{
		DeploymentName: deploymentName,
		ResourceGroup:  resourceGroup,
		Location:       location,
		SubscriptionID: subscriptionID,
		TenantID:       tenantID,
		MasterFQDN:     masterFQDN,
	}

	return stateOut
}

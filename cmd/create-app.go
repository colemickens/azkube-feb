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
			log.Println("starting create-app command")

			state := &util.State{}
			var err error
			state, err = ReadAndValidateState(statePath,
				[]reflect.Type{
					reflect.TypeOf(state.Common),
				},
				[]reflect.Type{
					reflect.TypeOf(state.App),
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

			RunCreateAppCmd(state, appName, appIdentifierURL)

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

	// need lots of ways of getting creds down into here
	// - interactive
	// - cert
	// - client secret

	return createAppCmd
}

func RunCreateAppCmd(state *util.State, appName, appURL string) {
	// we need to do two things
	// One: start oauth device flow
	// Two: create a ManualSecret that takes the access/refresh tokens directly.
	//deviceCode := azure.InitiateDeviceAuth(client, azure.PublicCloud.OAuthConfigForTenant(tenantID))
	//token := azure.CheckForUserCompletion(client, deviceCode)

	// do refresh dance to get valid tokens for everything
	d, err := util.NewDeployerWithToken(
		state.Common.SubscriptionID,
		state.Common.TenantID,
		"http://azkube",                                // client id
		"fk17GvW4GYj7Ju1g/sUGB4Jr39HQ+hiBW3VXTHRvnRE=") // client secret
	if err != nil {
		panic(err)
	}

	_, err = d.EnsureResourceGroup(
		state.Common.ResourceGroup,
		state.Common.Location,
		true)
	if err != nil {
		panic(err)
	}

	state.App, err = d.AdClient.CreateApp(*state.Common, appName, appURL)
	if err != nil {
		panic(err)
	}

	err = d.CreateRoleAssignment(*state.Common, state.App.ServicePrincipalObjectID)
	if err != nil {
		panic(err)
	}
}

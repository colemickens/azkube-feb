package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"

	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/azure-sdk-for-go/arm/resources"
	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	deployLongDescription = "long desc"
)

func NewDeployCmd() *cobra.Command {
	var configPath string

	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploy kubernetes on azure",
		Long:  deployLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting deploy command")

			var configIn util.DeployConfigIn
			configContents, err := ioutil.ReadFile(configPath)
			if err != nil {
				log.Fatalln(err)
			}
			err = json.Unmarshal(configContents, &configIn)
			if err != nil {
				log.Fatalln(err)
			}

			var config util.DeployConfigOut
			config.DeployConfigIn = configIn
			config.AppName = config.ResourceGroup + "-app"
			config.AppURL = "http://" + config.AppName + "/"
			config.VaultName = config.ResourceGroup + "-vault"
			config.ClientNames = []string{"default"}
			config.MasterFqdn = getMasterFQDN(config)

			RunDeployCmd(config)
			log.Println("finished deploy command")
		},
	}

	deployCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config")

	return deployCmd
}

func RunDeployCmd(config util.DeployConfigOut) {
	// TODO: load DeployerObjectId somehow??

	util.EnsureResourceGroup(config, true)

	err := util.GeneratePki(path.Join(config.OutputDirectory, "pki"))
	if err != nil {
		panic(err)
	}

	// TODO: create active directory app
	config.ServicePrincipalObjectID,
		err = util.CreateApp(config)
	if err != nil {
		panic(err)
	}

	err = util.GenerateSsh(path.Join(config.OutputDirectory, "ssh"))
	if err != nil {
		panic(err)
	}

	vaultTemplate, vaultParams, err := util.LoadAndFormat("vault", config, nil)
	if err != nil {
		panic(err)
	}

	deployClient := resources.NewDeploymentsClient(config.SubscriptionID)
	deployClient.Authorizer, err = util.GetAuthorizer(config, azure.AzureResourceManagerScope)
	if err != nil {
		panic(err)
	}

	_, err = util.DoDeployment(config, "vault", vaultTemplate, vaultParams, true)
	if err != nil {
		panic(err)
	}

	vaultClient := &autorest.Client{}
	vaultClient.Authorizer, err = util.GetAuthorizer(config, util.AzureVaultScope)
	if err != nil {
		panic(err)
	}

	config.ServicePrincipalSecretURL, err = util.UploadSecrets(config, vaultClient)
	if err != nil {
		panic(err)
	}

	myriadTemplate, myriadParams, err := util.LoadAndFormat("myriad", config, util.InsertCloudConfig)
	if err != nil {
		panic(err)
	}

	_, err = util.DoDeployment(config, "myriad", myriadTemplate, myriadParams, true)
	if err != nil {
		panic(err)
	}

	log.Println("done")
}

func getMasterFQDN(config util.DeployConfigOut) string {
	// TODO: this should be overrideable
	return config.ResourceGroup + "-master." + config.Location + ".cloudapp.azure.com"
}

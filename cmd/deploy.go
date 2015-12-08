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
	var outputPath string

	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploy kubernetes on azure",
		Long:  deployLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting deploy command")

			var configIn util.DeploymentConfig
			configContents, err := ioutil.ReadFile(configPath)
			if err != nil {
				log.Fatalln(err)
			}
			err = json.Unmarshal(configContents, &configIn)
			if err != nil {
				log.Fatalln(err)
			}

			var config util.DeploymentProperties
			config.DeploymentConfig = configIn
			config.App.AppName = config.ResourceGroup + "-app"
			config.App.AppURL = "http://" + config.App.AppName + "/"
			config.Vault.Name = config.ResourceGroup + "-vault"
			config.ClientNames = []string{"default"}
			config.MasterFqdn = getMasterFQDN(config)

			RunDeployCmd(config, outputPath)
			log.Println("finished deploy command")
		},
	}

	deployCmd.Flags().StringVarP(&configPath, "config", "c", "/etc/kubernetes/azure.json", "path to config")
	deployCmd.Flags().StringVarP(&outputPath, "output", "o", "./", "where to place output")

	return deployCmd
}

func RunDeployCmd(config util.DeploymentProperties, outputPath string) {
	util.EnsureResourceGroup(config, true)

	config.Pki, err := util.GeneratePki(path.Join(outputPath, "pki"))
	if err != nil {
		panic(err)
	}

	// TODO: create active directory app
	config.App, err = util.CreateApp(config, configIn.AppName, configIn.AppURL)
	if err != nil {
		panic(err)
	}

	err = util.GenerateSsh(path.Join(outputPath, "ssh"))
	if err != nil {
		panic(err)
	}

	vaultTemplate := util.CreateVaultTemplate(config)
	vaultParams := make(map[string]interface{})

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

	config.Secrets.ServicePrincipalSecretURL, err = util.UploadSecrets(config, vaultClient)
	if err != nil {
		panic(err)
	}

	myriadTemplate := util.CreateMyriadTemplate(config)
	myriadParams := make(map[string]interface{})

	_, err = util.DoDeployment(config, "myriad", myriadTemplate, myriadParams, true)
	if err != nil {
		panic(err)
	}

	log.Println("done")
}

func getMasterFQDN(config util.DeploymentProperties) string {
	// TODO: this should be overrideable
	// TODO: or add SAN support
	return config.ResourceGroup + "-master." + config.Location + ".cloudapp.azure.com"
}

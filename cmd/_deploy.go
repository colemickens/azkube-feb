// +build ignore
package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	// "github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	// "github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/azure"
	// "github.com/Azure/azure-sdk-for-go/arm/resources"
	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	deployLongDescription = "long desc"
)

func NewDeployCmd() *cobra.Command {
	var deploymentName string
	var location string
	var tenantID string
	var subscriptionID string

	var appClientId string
	var appClientCertificatePath string
	var appClientPrivateKeyPath string

	var outputPath

	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploy kubernetes on azure",
		Long:  deployLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting deploy command")

			if useExistingApp {
				if clientID == "" {
					panic("must specify clientID when using useExistingApp")
				}

				if clientCertificatePath = "" {
					panic("must specify clientCertificatePath when using useExistingApp")
				}
			}

			if deploymentName == "" {
				panic("must specify deploymentName")
			}
			if tenantID == "" {
				panic("must specify tenantId")
			}
			if subscriptionID == "" {
				panic("must specify subscriptionId")
			}
			if location == "" {
				panic("must specify location")
			}
			if appClientId == "" {

			}
			if appClientCertificatePath == "" {

			}
			if appClientPrivateKeyPath == "" {

			}



			RunDeployCmd(
				deploymentName,
				location,
				tenantID,
				subscriptionID,
				,
				clientID,
				clientCertificatePath)

			log.Println("finished deploy command")
		},
	}

	deployCmd.Flags().StringVarP(&deploymentName, "deploymentName", "n", "", "name for the deployment")
	deployCmd.Flags().StringVarP(&location, "location", "l", "", "location for deployed resources")
	deployCmd.Flags().StringVarP(&tenantID, "tenantId", "t", "", "tenant id for deployment")
	deployCmd.Flags().StringVarP(&subscriptionID, "subscriptionId", "s", "", "subscription id for deployment")

	deployCmd.Flags().BoolVarP(&appConfig, "appConfig", "u", false, "file containing app config information include the name/id-url of the app, along with the private key and certificate registered with the application")

	deployCmd.Flags().StringVarP(&outputPath, "outputPath", "o", "", "where to store outputs")

	return deployCmd
}

func RunDeployCmd(deploymentName, location, tenantID, subscriptionID string, existingAppConfig *AppProperties) {
	var d *util.Deployer

	// Determine if we need to create an application, or if we are using an existing one
	// Prefer to create an app per-deployment, but we need this for CI scenarios due to AD security

	var appProperties AppProperties

	if existingAppConfig != nil {
		var err error
		d, err = util.NewDeployerWithCertificate(
			subscriptionID,
			tenantID,
			clientID,
			clientCertificatePath)
		if err != nil {
			panic(err)
		}

		d.State.App = *existingAppConfig
	} else {
		// we need to do the actual oauth dance to get a token
		// this work is on pause because I'm not sure what all it will require
		// I think it requires pre-deploying the app ahead of time, which might be more work than its worth

		// create the app and use it's credentials to initialize a Deployer
	}

	d.State.Ssh, err = d.GenerateSsh(filepath.Join(outputPath, "ssh"))
	if err != nil {
		panic(err)
	}

	/*
	sshDirectory := filepath.Join(outputPath, "ssh")
	sshPrivateKeyPath := filepath.Join(sshDirectory, "id_rsa")
	err = os.MkdirAll(sshDirectory, 0777)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(sshPrivateKeyPath, d.State.Ssh.PrivateKeyPem, 0644)
	if err != nil {
		panic(err)
	}
	*/

	_, err = d.EnsureResourceGroup(deployConfig.ResourceGroup, deployConfig.Location, true)
	if err != nil {
		panic(err)
	}

	d.State.Pki, err = d.GeneratePki()
	if err != nil {
		panic(err)
	}

	vaultTemplateInput := util.VaultTemplateInput{
		VaultName: vaultName,
		TenantID: tenantID,
		ServicePrincipalObjectID: d.State.App.ServicePrincipalObjectID,
	}
	vaultTemplate, err := util.PopulateTemplate(util.VaultTemplate, vaultTemplateInput)
	if err != nil {
		panic(err)
	}

	vts, _ := json.Marshal(vaultTemplate)
	ioutil.WriteFile("/home/cole/test.txt.txt", vts, 0777)

	_, err = d.DoDeployment("vault", vaultTemplate, true)
	if err != nil {
		panic(err)
	}

	d.State.Secrets, err = d.UploadSecrets(vaultName)
	if err != nil {
		panic(err)
	}

	masterCloudConfig := util.PopulateTemplate(util.MasterCloudConfigTemplate, cloudConfigTemplateInput)
	if err != nil {
		panic(err)
	}

	minionCloudConfig := util.PopulateTemplate(util.MinionCloudConfigTemplate, cloudConfigTemplateInput)
	if err != nil {
		panic(err)
	}

	masterCloudConfig = util.Flatten(masterCloudConfig)
	minionCloudConfig = util.Flatten(minionCloudConfig)

	myriadTemplateInput := util.MyriadTemplateInput{
		DeploymentName: deploymentName,

		MasterVmSize: "A1_Dynamic",
		NodeVmSize: "",
		NodeVmssInitialCount: 3,
		Username: "azkube",

		VaultName: vaultName,
		ServicePrincipalSecretURL: d.State.Secrets.ServicePrincipalSecretURL,

		PodCidr: podCidr,
		ServiceCidr: serviceCidr,

		SshPublicKeyData: d.State.Ssh.OpenSshPublicKey,

		MasterCloudConfig: masterCloudConfig,
		MinionCloudConfig: minionCloudConfig,
	}

	myriadTemplate, err := util.PopulateTemplate(util.MyriadTemplate, myriadTemplateInput)
	if err != nil {
		panic(err)
	}

	_, err = d.DoDeployment("myriad", myriadTemplate, true)
	if err != nil {
		panic(err)
	}

	log.Println("done")
}

func getMasterFQDN(d *util.DeploymentConfig) string {
	// TODO: this should be overrideable
	// TODO: or add SAN support
	return d.ResourceGroup + "-master." + d.Location + ".cloudapp.azure.com"
}

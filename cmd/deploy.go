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
	var configPath string
	var outputPath string

	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploy kubernetes on azure",
		Long:  deployLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting deploy command")

			var deployConfig util.DeploymentConfig
			configContents, err := ioutil.ReadFile(configPath)
			if err != nil {
				log.Fatalln(err)
			}
			err = json.Unmarshal(configContents, &deployConfig)
			if err != nil {
				log.Fatalln(err)
			}

			RunDeployCmd(deployConfig, outputPath)
			log.Println("finished deploy command")
		},
	}

	deployCmd.Flags().StringVarP(&configPath, "config", "c", "/etc/kubernetes/azure.json", "path to config")
	deployCmd.Flags().StringVarP(&outputPath, "output", "o", "./deploy-output", "where to place output")

	return deployCmd
}

func RunDeployCmd(deployConfig util.DeploymentConfig, outputPath string) {
	// TODO(colemickens): make Deployer created universally
	d, err := util.NewDeployerWithToken(
		deployConfig.SubscriptionID,
		deployConfig.TenantID,
		deployConfig.DeployerClientID,
		deployConfig.DeployerClientSecret,
	)
	if err != nil {
		panic(err)
	}

	d.Config = deployConfig

	// TODO: fix config
	d.State.VaultName = deployConfig.VaultName

	d.State.Ssh, err = d.GenerateSsh()
	if err != nil {
		panic(err)
	}

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

	_, err = d.EnsureResourceGroup(deployConfig.ResourceGroup, deployConfig.Location, true)
	if err != nil {
		panic(err)
	}

	d.State.Pki, err = d.GeneratePki()
	if err != nil {
		panic(err)
	}

	// TODO: create active directory app
	d.State.App, err = d.CreateApp(deployConfig.AppName, deployConfig.AppURL)
	if err != nil {
		panic(err)
	}

	vaultTemplate, err := d.PopulateTemplate(util.VaultTemplate)
	if err != nil {
		panic(err)
	}

	vts, _ := json.Marshal(vaultTemplate)
	ioutil.WriteFile("/home/cole/test.txt.txt", vts, 0777)

	_, err = d.DoDeployment("vault", vaultTemplate, true)
	if err != nil {
		panic(err)
	}

	d.State.Secrets, err = d.UploadSecrets(d.Config.VaultName)
	if err != nil {
		panic(err)
	}

	d.State.MyriadConfig, err = d.LoadMyriadCloudConfigs()
	if err != nil {
		panic(err)
	}

	myriadTemplate, err := d.PopulateTemplate(util.MyriadTemplate)
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

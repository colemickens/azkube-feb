package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"

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
	deployCmd.Flags().StringVarP(&outputPath, "output", "o", "./", "where to place output")

	return deployCmd
}

func RunDeployCmd(deployConfig util.DeploymentConfig, outputPath string) {
	// TODO(colemickens): make Deployer created universally
	// with flexing auth for inside-of- and outside-of- Azure
	d, err := util.NewDeployerWithSecret("a", "b", "c", "d")

	_, err = d.EnsureResourceGroup(deployConfig.ResourceGroup, deployConfig.Location, true)
	if err != nil {
		panic(err)
	}

	d.State.Pki, err = d.GeneratePki(path.Join(outputPath, "pki"))
	if err != nil {
		panic(err)
	}

	// TODO: create active directory app
	d.State.App, err = d.CreateApp(deployConfig.AppName, deployConfig.AppURL)
	if err != nil {
		panic(err)
	}

	d.State.Ssh, err = d.GenerateSsh(path.Join(outputPath, "ssh"))
	if err != nil {
		panic(err)
	}

	vaultTemplate, err := d.PopulateTemplate(util.VaultTemplate)

	_, err = d.DoDeployment("vault", vaultTemplate, true)
	if err != nil {
		panic(err)
	}

	d.State.Secrets, err = d.UploadSecrets(d.State.VaultConfig.Name)
	if err != nil {
		panic(err)
	}

	d.State.MyriadConfig, err = d.LoadMyriadCloudConfigs()
	if err != nil {
		panic(err)
	}

	myriadTemplate := d.PopulateTemplate(util.MyriadTemplate)

	_, _, err = d.DoDeployment("myriad", myriadTemplate, true)
	if err != nil {
		panic(err)
	}

	log.Println("done")
}

func getMasterFQDN(deployProperties util.DeploymentProperties) string {
	// TODO: this should be overrideable
	// TODO: or add SAN support
	return d.State.ResourceGroup + "-master." + d.State.Location + ".cloudapp.azure.com"
}

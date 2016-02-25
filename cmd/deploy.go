package cmd

import (
	"fmt"
	"net"
	"path"
	"time"

	"github.com/colemickens/azkube/util"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

const (
	deployLongDescription = "creates a new kubernetes cluster in Azure"
)

func NewDeployCmd() *cobra.Command {
	var deployArgs util.DeployArguments

	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "creates a new kubernetes cluster in Azure",
		Long:  deployLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			err := validateDeployArgs(&deployArgs)
			if err != nil {
				glog.Exit("Failed to validate arguments for `deploy` command: %q", err)
			}
			err = deployRun(cmd, args, deployArgs)
			if err != nil {
				glog.Exit("Error occurred during deployment: %q", err)
			}
		},
	}

	flags := deployCmd.Flags()

	flags.StringVar(&deployArgs.OutputDirectory, "output-directory", "", "output directory (this is derived from --deployment-name if omitted)")
	flags.StringVar(&deployArgs.DeploymentName, "deployment-id", "", "deployment identifier (used to name output, resource group, and other resources)")
	flags.StringVar(&deployArgs.ResourceGroup, "resource-group", "", "resource group to deploy to (this is derived from --deployment-name if omitted)")
	flags.StringVar(&deployArgs.Location, "location", "brazilsouth", "location to deploy Azure resource (these can be found by running `azure location list` with azure-xplat-cli)")
	flags.StringVar(&deployArgs.MasterSize, "master-size", "Standard_A1", "size of the master virtual machine")
	flags.StringVar(&deployArgs.NodeSize, "node-size", "Standard_A1", "size of the node virtual machines")
	flags.IntVar(&deployArgs.NodeCount, "node-count", 3, "initial number of node virtual machines")
	flags.StringVar(&deployArgs.MasterFQDN, "master-fqdn", "", "main FQDN for master")
	flags.StringSliceVar(&deployArgs.MasterExtraFQDNs, "master-extra-fqdns", []string{}, "comma delimited list of SANs for the master")

	return deployCmd
}

func validateDeployArgs(deployArgs *util.DeployArguments) error {
	validateRootArgs(rootArgs)

	if deployArgs.DeploymentName == "" {
		deployArgs.DeploymentName = fmt.Sprintf("kube-%s", time.Now().Format("20060102-150405"))
		glog.Infof("--deployment-name is unset, generated a random deployment name: %q", deployArgs.DeploymentName)
	}

	if deployArgs.ResourceGroup == "" {
		deployArgs.ResourceGroup = deployArgs.DeploymentName
		glog.Infof("--resource-group is unset, derived one from --deployment-name: %q", deployArgs.ResourceGroup)
	}

	return nil
}

func deployRun(cmd *cobra.Command, args []string, deployArgs util.DeployArguments) error {
	d, err := util.NewDeployerFromCmd(rootArgs)
	if err != nil {
		return err
	}

	// Ensure the Resource Group exists
	_, err = d.EnsureResourceGroup(
		deployArgs.ResourceGroup,
		deployArgs.Location,
		true)
	if err != nil {
		panic(err)
	}

	// Create the Active Directory application
	appName := deployArgs.DeploymentName
	appURL := fmt.Sprintf("https://%s/", deployArgs.DeploymentName)
	applicationID, servicePrincipalObjectID, servicePrincipalClientSecret, err :=
		d.AdClient.CreateApp(appName, appURL)
	if err != nil {
		panic(err)
	}

	// Create the role assignment for the App/ServicePrincipal
	err = d.CreateRoleAssignment(rootArgs, deployArgs.ResourceGroup, servicePrincipalObjectID)
	if err != nil {
		panic(err)
	}

	// Create SSH key for deployment
	sshPrivateKey, sshPublicKeyString, err := util.GenerateSsh(path.Join(deployArgs.OutputDirectory, "private.key"))

	// (this is bound to the template which specifies subnet, for now
	// could declare them here and pass to template for better pairing
	masterIP := net.ParseIP("10.0.0.4")

	// Create PKI for deployment
	ca, apiserver, kubelet, kubeproxy, scheduler, replicationController, err :=
		util.CreateKubeCertificates(deployArgs.MasterFQDN, deployArgs.MasterExtraFQDNs, masterIP)
	if err != nil {
		return fmt.Errorf("error occurred while creating kube certificates")
	}

	// make a reserved ip for the kbue master TODO(colemick): for dns stability

	// Template Part 1: Generate dynamic config file (or put this in the config file and fill with template)
	// Template Part 2: blah

	_, _ = applicationID, servicePrincipalClientSecret
	_, _ = sshPrivateKey, sshPublicKeyString
	_, _, _, _, _, _ = ca, apiserver, kubelet, kubeproxy, scheduler, replicationController
	return nil
}

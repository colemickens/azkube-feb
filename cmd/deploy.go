package cmd

import (
	"crypto/rsa"
	"fmt"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	deployLongDescription = "creates a new kubernetes cluster in Azure"

	kubernetesStableReleaseURL = "https://github.com/kubernetes/kubernetes/releases/download/v1.1.8/kubernetes.tar.gz"
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
				log.Fatalf("Failed to validate arguments for `deploy` command: %q", err)
			}
			err = deployRun(cmd, args, deployArgs)
			if err != nil {
				log.Fatalf("Error occurred during deployment: %q", err)
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
	flags.StringVar(&deployArgs.KubernetesReleaseURL, "kubernetes-release-url", kubernetesStableReleaseURL, "size of the node virtual machines")
	flags.IntVar(&deployArgs.NodeCount, "node-count", 3, "initial number of node virtual machines")
	flags.StringVar(&deployArgs.Username, "username", "kube", "username to virtual machines")
	flags.StringVar(&deployArgs.MasterFQDN, "master-fqdn", "", "fqdn for master (used for PKI). calculated from cloudapp dns for master's public ip") // tODO is this wired up?
	flags.StringSliceVar(&deployArgs.MasterExtraFQDNs, "master-extra-fqdns", []string{}, "comma delimited list of SANs for the master")               // tODO is this wired up?

	return deployCmd
}

func validateDeployArgs(deployArgs *util.DeployArguments) error {
	validateRootArgs(rootArgs)

	// TODO: validate location + vmsizes, esp since used for masterfqdn

	if deployArgs.DeploymentName == "" {
		deployArgs.DeploymentName = fmt.Sprintf("kube-%s", time.Now().Format("20060102-150405"))
		log.Warnf("deployargs: --deployment-name is unset, generated a random deployment name: %q", deployArgs.DeploymentName)
	}

	if deployArgs.ResourceGroup == "" {
		deployArgs.ResourceGroup = deployArgs.DeploymentName
		log.Warnf("deployargs: --resource-group is unset, derived one from --deployment-name: %q", deployArgs.ResourceGroup)
	}

	if deployArgs.MasterFQDN == "" {
		deployArgs.MasterFQDN = fmt.Sprintf("%s.%s.cloudapp.azure.com", deployArgs.DeploymentName, deployArgs.Location)
		log.Warnf("deployargs: --master-fqdn is unset, derived from input: %q", deployArgs.MasterFQDN)
	}

	return nil
}

func deployRun(cmd *cobra.Command, args []string, deployArgs util.DeployArguments) error {
	d, err := util.NewDeployerFromCmd(rootArgs)
	if err != nil {
		return err
	}

	var (
		appName, appURL, applicationID, servicePrincipalObjectID, servicePrincipalClientSecret string

		sshPrivateKey      *rsa.PrivateKey
		sshPublicKeyString string

		ca, apiserver, client *util.PkiKeyCertPair
	)

	pkiLock := stepPki(d, deployArgs, ca, apiserver, client)
	sshLock := stepSsh(d, deployArgs, sshPrivateKey, &sshPublicKeyString)

	rgLock := stepRg(d, deployArgs)
	if err = <-rgLock; err != nil {
		return err
	}

	adLock := stepAd(d, deployArgs, &applicationID, &servicePrincipalObjectID, &servicePrincipalClientSecret)
	if err = <-adLock; err != nil {
		return err
	}
	if err = <-pkiLock; err != nil {
		return err
	}
	if err = <-sshLock; err != nil {
		return err
	}

	_, _ = appName, appURL

	deployLock := stepDeploy(d, deployArgs,
		applicationID, servicePrincipalObjectID, servicePrincipalClientSecret,
		sshPrivateKey, sshPublicKeyString,
		ca, apiserver, client)

	if err = <-deployLock; err != nil {
		return err
	}

	return nil
}

func stepRg(d *util.Deployer, deployArgs util.DeployArguments) chan error {
	var c chan error = make(chan error)

	go func() {
		var err error
		defer func() {
			c <- err
		}()

		_, err = d.EnsureResourceGroup(
			deployArgs.ResourceGroup,
			deployArgs.Location,
			true)
		if err != nil {
			return
		}
	}()

	return c
}

func stepAd(d *util.Deployer, deployArgs util.DeployArguments,
	applicationID, servicePrincipalObjectID, servicePrincipalClientSecret *string) chan error {
	var c chan error = make(chan error)

	go func() {
		var err error
		defer func() {
			c <- err
		}()

		appName := deployArgs.DeploymentName
		appURL := fmt.Sprintf("https://%s/", deployArgs.DeploymentName)
		*applicationID, *servicePrincipalObjectID, *servicePrincipalClientSecret, err =
			d.AdClient.CreateApp(appName, appURL)
		if err != nil {
			return
		}

		err = d.CreateRoleAssignment(rootArgs, deployArgs.ResourceGroup, *servicePrincipalObjectID)
		if err != nil {
			return
		}
	}()

	return c
}

func stepSsh(d *util.Deployer, deployArgs util.DeployArguments,
	sshPrivateKey *rsa.PrivateKey, sshPublicKeyString *string) chan error {
	var c chan error = make(chan error)

	go func() {
		var err error
		defer func() {
			c <- err
		}()
		sshPrivateKey, *sshPublicKeyString, err = util.GenerateSsh(path.Join(deployArgs.OutputDirectory, "private.key"))
		if err != nil {
			return
		}
	}()

	return c
}

func stepPki(d *util.Deployer, deployArgs util.DeployArguments,
	ca, apiserver, client *util.PkiKeyCertPair) chan error {
	var c chan error = make(chan error)

	go func() {
		var err error
		defer func() {
			c <- err
		}()
		ca, apiserver, client, err =
			util.CreateKubeCertificates(deployArgs.MasterFQDN, deployArgs.MasterExtraFQDNs)
		if err != nil {
			err = fmt.Errorf("error occurred while creating kube certificates")
			return
		}
	}()
	return c
}

func stepDeploy(d *util.Deployer, deployArgs util.DeployArguments,
	applicationID, servicePrincipalObjectID, servicePrincipalClientSecret string,
	sshPrivateKey *rsa.PrivateKey, sshPublicKeyString string,
	ca, apiserver, client *util.PkiKeyCertPair) chan error {
	var c chan error = make(chan error)

	go func() {
		var err error
		defer func() {
			c <- err
		}()

		// TODO(colemick, consider): make a reserved ip for the kbue master TODO(colemick): for dns stability
		flavorArgs := util.FlavorArguments{
			DeploymentName: deployArgs.DeploymentName,

			TenantID: rootArgs.TenantID,

			MasterSize:       deployArgs.MasterSize,
			NodeSize:         deployArgs.NodeSize,
			NodeCount:        deployArgs.NodeCount,
			Username:         deployArgs.Username,
			SshPublicKeyData: sshPublicKeyString,

			KubernetesReleaseURL: deployArgs.KubernetesReleaseURL,

			ServicePrincipalClientID:     applicationID,
			ServicePrincipalClientSecret: servicePrincipalClientSecret,

			MasterFQDN: deployArgs.MasterFQDN,

			CAKeyPair:        ca,
			ApiserverKeyPair: apiserver,
			ClientKeyPair:    client,
		}

		template, parameters, err := util.ProduceTemplateAndParameters(flavorArgs)
		if err != nil {
			return
		}

		_, err = d.DoDeployment(
			deployArgs.ResourceGroup,
			"myriad",
			template,
			parameters,
			true)
		if err != nil {
			return
		}

		_ = sshPrivateKey
		return
	}()
	return c
}

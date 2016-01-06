package cmd

import (
	"log"
	"reflect"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	deployMyriadLongDescription = "long desc"
)

func NewDeployMyriadCmd() *cobra.Command {
	var (
		statePath        string
		hyperkubeSpec    string
		serviceCidr      string
		podCidr          string
		dnsServiceIP     string
		username         string
		masterVmSize     string
		nodeVmSize       string
		initialNodeCount int
	)

	var deployMyriadCmd = &cobra.Command{
		Use:   "deploy-myriad",
		Short: "deploy the coreos kubernetes machines",
		Long:  deployMyriadLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting deploy-myriad command")

			var state *util.State
			var err error
			state, err = ReadAndValidateState(statePath,
				[]reflect.Type{
					reflect.TypeOf(state.Common),
					reflect.TypeOf(state.App),
					reflect.TypeOf(state.Ssh),
					reflect.TypeOf(state.Pki),
					reflect.TypeOf(state.Vault),
					reflect.TypeOf(state.Secrets),
				},
				[]reflect.Type{
					reflect.TypeOf(state.Myriad),
				},
			)

			RunDeployMyriadCmd(state, hyperkubeSpec, serviceCidr, podCidr, dnsServiceIP, username, masterVmSize, nodeVmSize, initialNodeCount)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished deploy-myriad command")
		},
	}

	deployMyriadCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")
	deployMyriadCmd.Flags().StringVarP(&hyperkubeSpec, "hyperkubeSpec", "h", "gcr.io/google_containers/hyperkube:v1.1.0", "hyperkube container spec")
	deployMyriadCmd.Flags().StringVarP(&serviceCidr, "serviceCidr", "c", "10.0.0.0/32", "the service cidr")
	deployMyriadCmd.Flags().StringVarP(&podCidr, "podCidr", "p", "10.0.0.0/32", "the pod cidr")
	deployMyriadCmd.Flags().StringVarP(&dnsServiceIP, "dnsServiceIp", "d", "10.0.0.10", "the dns service ip")
	deployMyriadCmd.Flags().StringVarP(&username, "username", "u", "azkube", "the username for the compute instances")
	deployMyriadCmd.Flags().StringVarP(&masterVmSize, "masterVmSize", "m", "A1_Standard", "the size of the VM for the master")
	deployMyriadCmd.Flags().StringVarP(&nodeVmSize, "nodeVmSize", "n", "A1_Standard", "the size of the VM for the nodes")
	deployMyriadCmd.Flags().IntVarP(&initialNodeCount, "initialNodeCount", "i", 3, "the initial number of nodes in the scale set")

	return deployMyriadCmd
}

func RunDeployMyriadCmd(state *util.State, hyperkubeSpec, serviceCidr, podCidr, dnsServiceIP, username, masterVmSize, nodeVmSize string, initialNodeCount int) {
	d, err := util.NewDeployerWithCertificate("a", "b", "c", "d", "e") // TODO (Colemickens): obviously
	if err != nil {
		panic(err)
	}

	configFileRendered := "" //anotehr template?

	cloudConfigTemplateInput := &util.CloudConfigTemplateInput{
		DeploymentName:         state.Common.DeploymentName,
		RuntimeConfigFile:      configFileRendered,
		HyperkubeContainerSpec: hyperkubeSpec,
		ServiceCidr:            serviceCidr,
		PodCidr:                podCidr,
		MasterIP:               state.Common.MasterIP.String(),
		MasterFQDN:             state.Common.MasterFQDN,
		DnsServiceIP:           dnsServiceIP,
	}

	masterCloudConfig, err := util.PopulateAndFlattenTemplate(util.MasterCloudConfigTemplate, cloudConfigTemplateInput)
	if err != nil {
		panic(err)
	}
	nodeCloudConfig, err := util.PopulateAndFlattenTemplate(util.NodeCloudConfigTemplate, cloudConfigTemplateInput)
	if err != nil {
		panic(err)
	}

	opensshPublicKeySlug, err := state.Ssh.OpenSshPublicKey()
	if err != nil {
		panic(err)
	}

	myriadTemplateInput := util.MyriadTemplateInput{
		DeploymentName:       state.Common.DeploymentName,
		MasterVmSize:         masterVmSize,
		NodeVmSize:           nodeVmSize,
		NodeVmssInitialCount: initialNodeCount,
		Username:             username,
		VaultName:            state.Vault.Name,
		PodCidr:              podCidr,
		ServiceCidr:          serviceCidr,
		SshPublicKeyData:     opensshPublicKeySlug,
		MasterCloudConfig:    masterCloudConfig,
		NodeCloudConfig:      nodeCloudConfig,
	}

	myriadTemplate, err := util.PopulateTemplate(util.MyriadTemplate, myriadTemplateInput)
	if err != nil {
		panic(err)
	}

	_, err = d.DoDeployment(*state.Common, "myriad", myriadTemplate, true)

	state.Myriad = &util.MyriadProperties{}
}

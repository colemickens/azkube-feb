package util

import (
	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	"github.com/Azure/azure-sdk-for-go/arm/resources"
)

type DeploymentConfig struct {
	Name             string
	ResourceGroup    string
	Location         string
	MasterVmSize     string
	NodeVmSize       string
	InitialNodeCount int

	MasterFqdn string
	Username   string

	VaultName string

	TenantID       string
	SubscriptionID string
	AppName        string
	AppURL         string

	DeployerObjectID     string
	DeployerClientID     string
	DeployerClientSecret string
	ClientNames          []string
}

type PkiProperties struct {
	ServicePrincipalFingerprint string
	ServicePrincipalPfx64       string
}

type SshProperties struct {
	OpenSshPublicKey string
	PrivateKeyPem    []byte
}

type AppProperties struct {
	AppURL                         string
	AppName                        string
	ApplicationID                  string
	ServicePrincipalCertificatePem string
	ServicePrincipalPrivateKeyPem  string
	ServicePrincipalObjectID       string
}

type SecretsProperties struct {
	ServicePrincipalSecretURL string
}

type NetworkProperties struct {
	VnetCidr   string
	SubnetCidr string

	PodCidr     string
	ServiceCidr string

	MasterPrivateIp string
	ApiServiceIp    string
	DnsServiceIp    string
}

type KubernetesProperties struct {
	HyperkubeContainerSpec string
}

type Deployer struct {
	Config DeploymentConfig
	State  DeploymentProperties

	DeploymentsClient resources.DeploymentsClient
	GroupsClient      resources.GroupsClient
	VaultClient       autorest.Client
	AdClient          autorest.Client
}

// appears in same order as in myriad variables
type DeploymentProperties struct {
	Pki *PkiProperties
	Ssh *SshProperties

	VaultName string

	App        *AppProperties
	Secrets    *SecretsProperties
	Network    *NetworkProperties
	Kubernetes *KubernetesProperties

	MyriadConfig *MyriadConfig
}

type MyriadConfig struct {
	MasterCloudConfig string
	NodeCloudConfig   string
}

package util

import (
	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	"github.com/Azure/azure-sdk-for-go/arm/resources"
)

type DeploymentConfig struct {
	Name                 string
	MasterVmSize         string
	NodeVmSize           string
	InitialNodeCount     int
	MasterFqdn           string
	Username             string
	TenantID             string
	SubscriptionID       string
	AppName              string
	AppURL               string
	ResourceGroup        string
	Location             string
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
	SshPublicKeyData64 string
}

type AppProperties struct {
	AppURL                   string
	AppName                  string
	ServicePrincipalObjectID string
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
}

// appears in same order as in myriad variables
type DeploymentProperties struct {
	Pki *PkiProperties
	Ssh *SshProperties

	App        *AppProperties
	Secrets    *SecretsProperties
	Network    *NetworkProperties
	Kubernetes *KubernetesProperties

	VaultConfig  *VaultConfig
	MyriadConfig *MyriadConfig
}

type VaultConfig struct {
	Name                     string
	ServicePrincipalObjectID string
	DeployerObjectID         string
}

type MyriadConfig struct {
	MasterCloudConfig string
	NodeCloudConfig   string
}

type ScaleConfig struct {
	NodeVmssName string
	NodeCount    int
}

package util

import (
	"crypto/rsa"
	"crypto/x509"

	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	"github.com/Azure/azure-sdk-for-go/arm/resources"
	kapi "k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api"
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
	CACertificate                   x509.Certificate
	ApiServerCertificate            x509.Certificate
	KubeletKubeconfig               kapi.Config
	KubeproxyKubeconfig             kapi.Config
	SchedulerKubeconfig             kapi.Config
	ReplicationControllerKubeconfig kapi.Config
}

type SshProperties struct {
	OpenSshPublicKey string
	PrivateKeyPem    []byte
}

type AppProperties struct {
	ID                          string
	Name                        string
	IdentifierURL               string
	ServicePrincipalCertificate x509.Certificate
	ServicePrincipalPrivateKey  rsa.PrivateKey
	ServicePrincipalObjectID    string
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

type VaultTemplateInput struct {
	VaultName                string
	TenantID                 string
	ServicePrincipalObjectId string
}

type CloudConfigTemplateInput struct {
	DeploymentName string

	ConfigFile             string
	HyperkubeContainerSpec string

	ServiceCidr string
	PodCidr     string

	MasterIP   string
	MasterFQDN string

	DnsServiceIP string
}

type MyriadTemplateInput struct {
	DeploymentName string

	MasterVmSize         string
	NodeVmSize           string
	NodeVmssInitialCount int
	Username             string

	VaultName                 string
	ServicePrincipalSecretURL string

	PodCidr     string
	ServiceCidr string

	SshPublicKeyData string

	MasterCloudConfig string
	MinionCloudConfig string
}

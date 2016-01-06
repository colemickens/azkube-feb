package util

import (
	"crypto/rsa"
	"crypto/x509"
	"net"

	"github.com/Azure/azure-sdk-for-go/arm/resources"
)

type CommonProperties struct {
	DeploymentName string
	ResourceGroup  string
	Location       string
	TenantID       string
	SubscriptionID string
	MasterFQDN     string
	MasterIP       net.IP // TODO(colemickens): populate this
}

type AppProperties struct {
	ApplicationID               string
	Name                        string
	IdentifierURL               string
	ServicePrincipalCertificate x509.Certificate
	ServicePrincipalPrivateKey  *rsa.PrivateKey
	ServicePrincipalObjectID    string
}

type PkiKeyCertPair struct {
	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey
}

type PkiProperties struct {
	CA                    *PkiKeyCertPair
	ApiServer             *PkiKeyCertPair
	Kubelet               *PkiKeyCertPair
	Kubeproxy             *PkiKeyCertPair
	Scheduler             *PkiKeyCertPair
	ReplicationController *PkiKeyCertPair
}

type SshProperties struct {
	PrivateKey *rsa.PrivateKey
}

type VaultProperties struct {
	Name string
}

type SecretsProperties struct {
	ServicePrincipalSecretURL string
}

type MyriadProperties struct {
	PodCidr     string
	ServiceCidr string

	MasterPrivateIp string
	ApiServiceIp    string
	DnsServiceIp    string

	HyperkubeContainerSpec string
}

// appears in same order as in variable in ARM template
type State struct {
	Common  *CommonProperties
	App     *AppProperties
	Pki     *PkiProperties
	Ssh     *SshProperties
	Vault   *VaultProperties
	Secrets *SecretsProperties
	Myriad  *MyriadProperties
}

type Deployer struct {
	DeploymentsClient resources.DeploymentsClient
	GroupsClient      resources.GroupsClient
	VaultClient       VaultClient
	//AdClient          AdClient
}

type VaultTemplateInput struct {
	VaultName                string
	TenantID                 string
	ServicePrincipalObjectID string
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
	NodeCloudConfig   string
}

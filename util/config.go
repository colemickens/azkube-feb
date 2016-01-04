package util

import (
	"crypto/rsa"
	"crypto/x509"

	//"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	"github.com/Azure/azure-sdk-for-go/arm/resources"
)

type CommonProperties struct {
	DeploymentName string
	Location       string
	TenantID       string
	SubscriptionID string
	ResourceGroup  string
}

type AppProperties struct {
	ApplicationID               string
	Name                        string
	IdentifierURL               string
	ServicePrincipalCertificate x509.Certificate
	ServicePrincipalPrivateKey  rsa.PrivateKey
	ServicePrincipalObjectID    string
}

type PkiProperties struct {
	CACertificate x509.Certificate

	ApiServerCertificate             x509.Certificate
	ApiServerPrivateKey              rsa.PrivateKey
	KubeletCertificate               x509.Certificate
	KubeletPrivateKey                rsa.PrivateKey
	KubeproxyCertificate             x509.Certificate
	KubeproxyPrivateKey              rsa.PrivateKey
	SchedulerCertificate             x509.Certificate
	SchedulerPrivateKey              rsa.PrivateKey
	ReplicationControllerCertificate x509.Certificate
	ReplicationControllerPrivateKey  rsa.PrivateKey
}

type SshProperties struct {
	OpenSshPublicKey string
	PrivateKeyPem    []byte
}

type VaultProperties struct {
	Name string
}

type SecretsProperties struct {
	ServicePrincipalSecretURL string
}

type MyriadProperties struct {
	VnetCidr   string
	SubnetCidr string

	PodCidr     string
	ServiceCidr string

	MasterPrivateIp string
	ApiServiceIp    string
	DnsServiceIp    string

	HyperkubeContainerSpec string
}

// appears in same order as in myriad variables
type State struct {
	Common           *CommonProperties
	App              *AppProperties
	Pki              *PkiProperties
	Ssh              *SshProperties
	Vault            *VaultProperties
	MyriadProperties *MyriadProperties
}

type Deployer struct {
	DeploymentsClient resources.DeploymentsClient
	GroupsClient      resources.GroupsClient
	VaultClient       VaultClient
	//	AdClient          AdClient
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

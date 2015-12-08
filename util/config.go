package util

type DeploymentConfig struct {
	Name                 string
	MasterVmSize         string
	NodeVmSize           string
	InitialNodeCount     int
	MasterFqdn           string
	Username             string
	TenantID             string
	SubscriptionID       string
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

type VaultProperties struct {
	Name string
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

type CloudConfigProperties struct {
	Master string
	Node   string
}

// appears in same order as in myriad variables
type DeploymentProperties struct {
	DeploymentConfig

	Pki PkiProperties
	Ssh SshProperties

	App        AppProperties
	Vault      VaultProperties
	Secrets    SecretsProperties
	Network    NetworkProperties
	Kubernetes KubernetesProperties

	CloudConfig CloudConfigProperties
}

type VaultConfig struct {
	Name string
}

type ScaleConfig struct {
	NodeVmssName string
	NodeCount    int
}

package util

import ()

type RootArguments struct {
	TenantID              string
	SubscriptionID        string
	AuthMethod            string
	ClientID              string
	ClientSecret          string
	ClientCertificatePath string
	PrivateKeyPath        string
}

type DeployArguments struct {
	OutputDirectory      string
	DeploymentName       string
	ResourceGroup        string
	Location             string
	MasterSize           string
	NodeSize             string
	NodeCount            int
	Username             string
	MasterFQDN           string
	MasterExtraFQDNs     []string
	KubernetesReleaseURL string
}

// part of interface to flavors.
// this is what flavors get.
// make a way to opt-in/opt-out?
type FlavorArguments struct {
	//TenantID       string (derivable)
	//SubscriptionID string (derivable from template func/vars)
	//ResourceGroup  string (derivable from template stuffs)
	//Location       string (derivable from template stuffs)

	MasterSize       string
	NodeSize         string
	NodeCount        int
	Username         string
	SshPublicKeyData string

	ServicePrincipalClientID     string
	ServicePrincipalClientSecret string

	MasterFQDN string

	CAKeyPair        *PkiKeyCertPair
	ApiserverKeyPair *PkiKeyCertPair
	ClientKeyPair    *PkiKeyCertPair
}

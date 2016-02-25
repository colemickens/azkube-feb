package util

import (
	"net"
)

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
	MasterFQDN           string
	MasterExtraFQDNs     []string
	MasterIP             net.IP
	KubernetesReleaseURL string
}

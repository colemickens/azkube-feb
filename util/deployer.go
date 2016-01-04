package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/azure-sdk-for-go/arm/resources"
)

func NewDeployerWithToken(subscriptionID, tenantID, clientID, clientSecret string) (deployer *Deployer, err error) {
	secret := &azure.ServicePrincipalTokenSecret{
		ClientSecret: clientSecret,
	}

	return newDeployer(subscriptionID, tenantID, clientID, secret)
}

func NewDeployerWithCertificate(subscriptionID, tenantID, appURL, certPath, keyPath string) (deployer *Deployer, err error) {
	certificateData, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Fatalln("failed", err)
	}

	block, _ := pem.Decode(certificateData)
	if block == nil {
		panic("failed to decode a pem block from certificate pem")
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}

	privateKey, err := parseRsaPrivateKey(keyPath)
	if err != nil {
		panic(err)
	}

	secret := &azure.ServicePrincipalCertificateSecret{
		Certificate: certificate,
		PrivateKey:  privateKey,
	}

	return newDeployer(subscriptionID, tenantID, appURL, secret)
}

func withSecret(tenantID, clientID, scope string, secret azure.ServicePrincipalSecret) (spt *azure.ServicePrincipalToken, err error) {
	spt, err = azure.NewServicePrincipalTokenWithSecret(
		clientID,
		tenantID,
		scope,
		secret,
	)
	if err != nil {
		return nil, err
	}

	err = spt.Refresh()
	if err != nil {
		return nil, err
	}

	return spt, nil
}

func newDeployer(subscriptionID, tenantID, clientID string, secret azure.ServicePrincipalSecret) (deployer *Deployer, err error) {
	deployer = &Deployer{}
	deployer.DeploymentsClient = resources.NewDeploymentsClient(subscriptionID)
	deployer.GroupsClient = resources.NewGroupsClient(subscriptionID)
	// deployer.AdClient = AdClient{ autorest.Client{} }
	deployer.VaultClient = VaultClient{autorest.Client{}}

	resourcesScopeSpt, err := withSecret(tenantID, clientID, azure.AzureResourceManagerScope, secret)
	if err != nil {
		return nil, err
	}
	vaultScopeSpt, err := withSecret(tenantID, clientID, AzureVaultScope, secret)
	if err != nil {
		return nil, err
	}
	//adScopeSpt, err := withSecret(tenantID, clientID, AzureActiveDirectoryScope, secret)
	//if err != nil {
	//	return nil, err
	//}

	deployer.DeploymentsClient.Authorizer = resourcesScopeSpt
	deployer.GroupsClient.Authorizer = resourcesScopeSpt
	deployer.VaultClient.Authorizer = vaultScopeSpt
	//deployer.AdClient.Authorizer = adScopeSpt

	return deployer, nil
}

func parseRsaPrivateKey(path string) (*rsa.PrivateKey, error) {
	privateKeyData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("failed", err)
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		panic("failed to decode a pem block from private key pem")
	}

	privatePkcs1Key, errPkcs1 := x509.ParsePKCS1PrivateKey(block.Bytes)
	if errPkcs1 == nil {
		return privatePkcs1Key, nil
	}

	privatePkcs8Key, errPkcs8 := x509.ParsePKCS8PrivateKey(block.Bytes)
	if errPkcs8 == nil {
		privatePkcs8RsaKey, ok := privatePkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("Pkcs8 contained non-RSA key. Expected RSA key.")
		}
		return privatePkcs8RsaKey, nil
	}

	return nil, fmt.Errorf("Failed to parse private key as Pkcs#1 or Pkcs#8. (%s). (%s).", errPkcs1, errPkcs8)
}

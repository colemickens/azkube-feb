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

const (
// TODO eliminate
// see vault.go
// AzureVaultScope = "https://vault.azure.com"
)

// TODO(colemickens): eliminate some code duplication with ServicePrincipalSecret

func NewDeployerWithSecret(subscriptionID, tenantID, clientID, clientSecret string) (deployer *Deployer, err error) {
	deployer = newDeployer(subscriptionID)

	resourcesScopeSpt, err := withSecretAndScope(tenantID, clientID, clientSecret, azure.AzureResourceManagerScope)
	if err != nil {
		return nil, err
	}
	vaultScopeSpt, err := withSecretAndScope(tenantID, clientID, clientSecret, AzureVaultScope)
	if err != nil {
		return nil, err
	}

	deployer.DeploymentsClient.Authorizer = resourcesScopeSpt
	deployer.GroupsClient.Authorizer = resourcesScopeSpt
	deployer.VaultClient.Authorizer = vaultScopeSpt

	return deployer, nil
}

func NewDeployerWithCertificate(subscriptionID, tenantID, appURL, certPath, keyPath string) (deployer *Deployer, err error) {
	deployer = newDeployer(subscriptionID)

	resourcesScopeSpt, err := withCertAndScope(tenantID, appURL, certPath, keyPath, azure.AzureResourceManagerScope)
	if err != nil {
		return nil, err
	}
	vaultScopeSpt, err := withCertAndScope(tenantID, appURL, certPath, keyPath, AzureVaultScope)
	if err != nil {
		return nil, err
	}

	deployer.DeploymentsClient.Authorizer = resourcesScopeSpt
	deployer.GroupsClient.Authorizer = resourcesScopeSpt
	deployer.VaultClient.Authorizer = vaultScopeSpt

	return deployer, nil
}

func newDeployer(subscriptionID string) (deployer *Deployer) {
	deployer = &Deployer{}
	deployer.DeploymentsClient = resources.NewDeploymentsClient(subscriptionID)
	deployer.GroupsClient = resources.NewGroupsClient(subscriptionID)
	deployer.VaultClient = autorest.Client{}
	return deployer
}

func withSecretAndScope(tenantID, clientID, clientSecret, scope string) (spt *azure.ServicePrincipalToken, err error) {
	spt, err = azure.NewServicePrincipalToken(
		clientID,
		tenantID,
		scope,
		clientSecret,
	)
	if err != nil {
		return nil, err
	}

	// TODO(colemickens): refresh token here so we blow up during init and not later

	return spt, nil
}

func withCertAndScope(tenantID, appURL, certPath, keyPath, scope string) (spt *azure.ServicePrincipalToken, err error) {
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

	spt, err = azure.NewServicePrincipalTokenFromCertificate(
		appURL,
		certificate,
		privateKey,
		tenantID,
		AzureVaultScope)

	// TODO(colemickens): refresh token here so we blow up during init and not later

	return spt, err
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

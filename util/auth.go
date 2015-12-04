package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/azure"
)

func GetAuthorizer(config DeployConfigOut, scope string) (servicePrincipalToken *azure.ServicePrincipalToken, err error) {
	spt, err := azure.NewServicePrincipalToken(
		config.ClientID,
		config.TenantID,
		scope,
		config.ClientSecret,
	)
	if err != nil {
		return nil, err
	}

	return spt, nil
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

func GetAuthorizerInAzure(config DeployConfigOut, scope string) (servicePrincipalToken *azure.ServicePrincipalToken, err error) {
	certificateData, err := ioutil.ReadFile("/var/lib/waagent/" + config.ServicePrincipalFingerprint + ".crt")
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

	privateKey, err := parseRsaPrivateKey("/var/lib/waagent/" + config.ServicePrincipalFingerprint + ".key")
	if err != nil {
		panic(err)
	}

	spt, err := azure.NewServicePrincipalTokenFromCertificate(
		config.AppURL,
		certificate,
		privateKey,
		config.TenantID,
		AzureVaultScope)
	return spt, err
}

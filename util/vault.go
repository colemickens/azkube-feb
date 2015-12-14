package util

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
)

const (
	AzureVaultApiVersion     = "2015-06-01"
	AzureVaultScope          = "https://vault.azure.net/"
	AzureVaultSecretTemplate = "https://{vault-name}.vault.azure.net/{secret-name}/{secret-version}"
)

type Secret struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

var cachedVaultClient *autorest.Client = nil

func (d *Deployer) PutSecret(vaultName, secretName, secretPath string) (secretURL string, err error) {
	secretID := secretName // at first it's just the name, hopefully later its name/version

	pathParams := map[string]interface{}{
		"secret-id": secretID,
	}

	var result struct {
		Id string `json:"id"`
	}

	baseURL := strings.Replace(AzureVaultSecretTemplate, "{vault-name}", vaultName, 1)

	req, err := autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(baseURL),
		autorest.WithPath("/secrets/{secret-id}"),
		autorest.WithPathParameters(pathParams))

	if err != nil {
		return "", err
	}

	resp, err := d.VaultClient.Send(req, http.StatusOK)
	if err != nil {
		return "", err
	}

	err = autorest.Respond(
		resp,
		d.VaultClient.ByInspecting(),
		autorest.WithErrorUnlessOK(),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	if err != nil {
		return "", err
	}

	log.Println("!!!! need to check output to find secret url")
	log.Println("secret id:", result.Id)

	return "", nil
}

func (d *Deployer) UploadSecrets(vaultName string) (secretsProperties *SecretsProperties, err error) {
	// TODO(colemickens): populate this from the same place that is consumed from
	secrets := map[string]string{
		"pki/ca.crt":                               "ca-crt",
		"pki/apiserver.crt":                        "apiserver-crt",
		"pki/apiserver.key":                        "apiserver-key",
		"pki/node-proxy-kubeconfig":                "node-proxy-kubeconfig",
		"pki/node-kubelet-kubeconfig":              "node-kubelet-kubeconfig",
		"pki/master-proxy-kubeconfig":              "master-proxy-kubeconfig",
		"pki/master-kubelet-kubeconfig":            "master-kubelet-kubeconfig",
		"pki/master-scheduler-kubeconfig":          "master-scheduler-kubeconfig",
		"pki/master-controller-manager-kubeconfig": "master-controller-manager-kubeconfig",
	}

	servicePrincipalSecretURL, err := d.PutSecret(
		vaultName,
		"servicePrincipal-pfx",
		"pki/servicePrincipal.pfx",
	)
	if err != nil {
		return nil, err
	}

	for secretPath, secretName := range secrets {
		_, err = d.PutSecret(vaultName, secretName, secretPath)
		if err != nil {
			return nil, err
		}
	}

	secretsProperties = &SecretsProperties{
		ServicePrincipalSecretURL: servicePrincipalSecretURL,
	}

	return secretsProperties, nil
}

func (d *Deployer) GetSecret(vaultName, secretName string) (*string, error) {
	p := map[string]interface{}{
		"secret-name":    secretName,
		"secret-version": "",
	}
	q := map[string]interface{}{
		"api-version": AzureVaultApiVersion,
	}

	secretURL := strings.Replace(AzureVaultSecretTemplate, "{vault-name}", vaultName, -1)

	req, err := autorest.Prepare(&http.Request{},
		autorest.AsGet(),
		autorest.WithBaseURL(secretURL),
		autorest.WithPathParameters(p),
		autorest.WithQueryParameters(q))

	if err != nil {
		panic(err)
	}

	resp, err := d.VaultClient.Send(req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var secret Secret

	err = autorest.Respond(
		resp,
		autorest.ByUnmarshallingJSON(&secret))
	if err != nil {
		return nil, err
	}

	secretValue, err := base64.StdEncoding.DecodeString(secret.Value)
	if err != nil {
		return nil, err
	}

	secretValueString := string(secretValue)

	return &secretValueString, nil
}

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
	AzureVaultScope          = "https://vault.azure.net"
	AzureVaultSecretTemplate = "https://{vault-name}.vault.azure.net/{secret-name}/{secret-version}"
)

type VaultClient struct {
	autorest.Client
}

type Secret struct {
	ID         *string          `json:"id,omitempty"`
	Value      string           `json:"value"`
	Attributes SecretAttributes `json:"attributes"`
}

type SecretAttributes struct {
	Enabled   bool    `json:"enabled"`
	NotBefore *string `json:"nbf"`
	Expires   *string `json:"exp"`
}

func (v *VaultClient) PutSecret(vaultName, secretName, secretValue string) (secretURL string, err error) {
	secretID := secretName // at first it's just the name, hopefully later its name/version

	pathParams := map[string]interface{}{
		"secret-id": secretID,
	}

	q := map[string]interface{}{
		"api-version": AzureVaultApiVersion,
	}

	var result struct {
		ID string `json:"id"`
	}

	secretValue64 := base64.URLEncoding.EncodeToString([]byte(secretValue))

	secret := Secret{
		//ID:    secretID,
		Value: secretValue64,
		Attributes: SecretAttributes{
			Enabled:   true,
			NotBefore: nil,
			Expires:   nil,
		},
	}

	baseURL := strings.Replace(AzureVaultSecretTemplate, "{vault-name}", vaultName, 1)

	req, err := autorest.Prepare(
		&http.Request{},
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(baseURL),
		autorest.WithPath("/secrets/{secret-id}"),
		autorest.WithPathParameters(pathParams),
		autorest.WithQueryParameters(q),
		autorest.WithJSON(secret),
	)

	if err != nil {
		return "", err
	}

	resp, err := v.Send(req, http.StatusOK)
	if err != nil {
		return "", err
	}

	err = autorest.Respond(
		resp,
		v.ByInspecting(),
		autorest.WithErrorUnlessOK(),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	if err != nil {
		return "", err
	}

	log.Println("!!!! need to check output to find secret url")
	log.Println("secret id:", result.ID)

	return "", nil
}

func (v *VaultClient) GetSecret(vaultName, secretName string) (*string, error) {
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

	resp, err := v.Send(req, http.StatusOK)
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

	secretValue, err := base64.URLEncoding.DecodeString(secret.Value)
	if err != nil {
		return nil, err
	}

	secretValueString := string(secretValue)

	return &secretValueString, nil
}

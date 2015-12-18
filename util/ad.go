package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
)

const (
	AzureActiveDirectoryScope = "https://graph.windows.net/"
	//AzureActiveDirectoryApiVersion = "1.6"
	AzureActiveDirectoryApiVersion = "1.5"
	//AzureActiveDirectoryApiVersion        = "1.42-previewInternal"
	AzureRoleManagementApiVersion         = "2015-07-01"
	AzureActiveDirectoryBaseURL           = "https://graph.windows.net/{tenant-id}"
	AzureManagementBaseURL                = "https://management.azure.net/{tenant-id}"
	AzureActiveDirectoryContributorRoleId = "b24988ac-6180-42a0-ab88-20f7382dd24c"
)

type AdApplication struct {
	ApplicationID string `json:"appId,omitempty"`    // readonly
	ObjectID      string `json:"objectId,omitempty"` // readonly

	AvailableToOtherTenants bool              `json:"availableToOtherTenants"`
	DisplayName             string            `json:"displayName,omitempty"`
	Homepage                string            `json:"homepage,omitempty"`
	IdentifierURLs          []string          `json:"identifierUrls,omitempty"`
	KeyCredentials          []AdKeyCredential `json:"keyCredentials,omitempty"`
}

type AdKeyCredential struct {
	KeyId     string `json:"keyId,omitempty"`
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	Type      string `json:"type,omitempty"`
	Usage     string `json:"usage,omitempty"`
	Value     string `json:"value,omitempty"`
}

type AdServicePrincipal struct {
	ObjectID string `json:"objectId,omitempty"` // readonly

	ApplicationID         string   `json:"appId,omitempty"`
	AccountEnabled        bool     `json:"accountEnabled,omitempty"`
	ServicePrincipalNames []string `json:"servicePrincipalNames,omitempty"`
}

type AdRoleAssignment struct {
	RoleDefinitionID string `json:"roleDefinitionId,omitempty"`
	PrincipalID      string `json:"principalId,omitempty"`
}

func (d *Deployer) CreateApp(appName, appURL string) (appProperties *AppProperties, err error) {
	appProperties = &AppProperties{}

	notBefore := time.Now()
	notAfter := time.Now().Add(5 * 365 * 24 * time.Hour)
	notAfter = time.Now().Add(10000 * 24 * time.Hour)

	// create the service principal's private key
	privateKey, err := rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		return nil, err
	}

	privateKeyPemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyPemBytes := pem.EncodeToMemory(privateKeyPemBlock)

	appProperties.ServicePrincipalPrivateKeyPem = string(privateKeyPemBytes)

	// create the cert and store it in state
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}
	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "Azkube",
			Organization: []string{"Azkube"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	certificateDer, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	certificatePemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificateDer})
	ioutil.WriteFile("last-cert.cer", certificatePemBytes, 0777) // TODO(colemickens): remove

	certificatePem := string(certificatePemBytes)
	log.Println(certificatePem)
	parts := strings.Split(certificatePem, "\n")
	certificatePem = strings.Join(parts[1:len(parts)-2], "")   // 500
	certificatePem = strings.Join(parts[1:len(parts)-2], "\n") // 500
	//certificatePemBytes, _ = json.Marshal(certificatePem)
	//certificatePem = string(certificatePemBytes)
	_ = json.Marshal

	appProperties.ServicePrincipalCertificatePem = certificatePem

	// create application
	applicationReq := AdApplication{
		AvailableToOtherTenants: false,
		DisplayName:             d.Config.AppName,
		Homepage:                d.Config.AppURL,
		IdentifierURLs:          []string{d.Config.AppURL},
		KeyCredentials: []AdKeyCredential{
			AdKeyCredential{
				KeyId:     uuid.New(),
				Type:      "AsymmetricX509Cert",
				Usage:     "Verify",
				StartDate: notBefore.Format(time.RFC3339),
				EndDate:   notAfter.Format(time.RFC3339),
				Value:     appProperties.ServicePrincipalCertificatePem,
			},
		},
	}

	p := map[string]interface{}{"tenant-id": d.Config.TenantID}
	q := map[string]interface{}{"api-version": AzureActiveDirectoryApiVersion}

	forceItIn := `{"availableToOtherTenants":false,"displayName":"test01113","homepage":"http://test000113","identifierUris":["http://test000113"],"keyCredentials":[{"startDate":"2015-12-18T08:25:21.694Z","endDate":"2043-05-05T08:02:54.000Z","value":"MIIDBDCCAeygAwIBAgIRAMVbOzlLG6kqcxqSBJNcHlowDQYJKoZIhvcNAQELBQAw\nIjEPMA0GA1UEChMGQXprdWJlMQ8wDQYDVQQDEwZBemt1YmUwHhcNMTUxMjE4MDgx\nODMxWhcNNDMwNTA1MDgxODMxWjAiMQ8wDQYDVQQKEwZBemt1YmUxDzANBgNVBAMT\nBkF6a3ViZTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKvd96W+mUO8\nNKbexxKLmrkOLVixzroDiDew3uIT94miwvDcnNCx79XXPOl57cyPwj8Hs5z47gqY\nuwlUayVSCLl81ixog3k0mwAh03O1W3pwdafpCC6qOOPKY4daEjS4raGj8pPvpEto\n2Tv5jGTpuKAeqmw/G0t3cGSs9ruOBO0WEtzd6of0zXTLEtt6SnKbeLrrU25g6qBg\nSPSPPtT88zNjhwA0q1FwSOEpbTjnWw1Ujw6RAk4xF+2wImeYAkwcix50zWAErLkb\ndrWszoEaX1H4mlhb80TOnnksfoQctDVfnS8IUxDIri9Dy15APKyoRG45l+ni6dS+\nDmc4Uv/VP9UCAwEAAaM1MDMwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsG\nAQUFBwMBMAwGA1UdEwEB/wQCMAAwDQYJKoZIhvcNAQELBQADggEBAIiuwS/ZkUiQ\ncgakyVLnzb/egvhlSB+ol72IyKiKa4PzgsX+JxhHC5+5p2g8IY5o3mfK2GxkiA23\nOdwxJDGepJuQrU3Aue2DC8U7PdCH/PweF+TWXU2DeyYp/8fhzD12gIS8TmJfgozM\nYiMaf9dyKFU/GrbosfjNS6GUrhbeZYKF3EfqwEe8Igv4JqQvQWpxxv8ckyisdBlq\nTw/GY8XWAfOnhYwVU0q1zS3aQpKahwIZBGqDNulXBHgRs2261UykSIUXyVc5Uawv\nx7uVdZgLjCfhR1XuBiBuTwyWb5xkbkP/lJLrl06UmJbuwcj1+pyXfnevCDbfnPNO\nGFtPzdmktNw=","keyId":"9dfcaeb7-3367-4387-ab92-52da1a789560","usage":"Verify","type":"AsymmetricX509Cert"}]}`

	var blahf map[string]interface{}
	json.Unmarshal([]byte(forceItIn), &blahf)
	_ = applicationReq

	req, err := autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(AzureActiveDirectoryBaseURL),
		autorest.WithPath("applications"),
		autorest.WithPathParameters(p),
		autorest.WithQueryParameters(q),
		autorest.WithJSON(blahf))

	log.Println(req)

	if err != nil {
		return nil, err
	}

	resp, err := d.AdClient.Send(req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var applicationResp AdApplication
	err = autorest.Respond(resp, autorest.ByUnmarshallingJSON(&applicationResp))
	if err != nil {
		return nil, err
	}

	appProperties.ApplicationID = applicationResp.ApplicationID

	log.Println("sleep 10")
	time.Sleep(10 * time.Second)

	// create service principal
	servicePrincipalReq := AdServicePrincipal{
		ApplicationID:         appProperties.ApplicationID,
		AccountEnabled:        true,
		ServicePrincipalNames: []string{d.Config.AppURL},
	}

	req, err = autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(AzureActiveDirectoryBaseURL),
		autorest.WithPath("servicePrincipals"),
		autorest.WithPathParameters(p),
		autorest.WithQueryParameters(q),
		autorest.WithJSON(servicePrincipalReq))
	if err != nil {
		return nil, err
	}

	resp, err = d.AdClient.Send(req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var servicePrincipalResp AdServicePrincipal
	err = autorest.Respond(resp, autorest.ByUnmarshallingJSON(&servicePrincipalResp))
	if err != nil {
		return nil, err
	}

	appProperties.ServicePrincipalObjectID = servicePrincipalResp.ObjectID

	log.Println("sleep 10")
	time.Sleep(10 * time.Second)

	// create role assignment for service principal

	roleDefinitionId := "/subscriptions/{subscription-id}/providers/Microsoft.Authorization/roleDefinitions/{role-definition-id}"
	roleDefinitionId = strings.Replace(roleDefinitionId, "{subscription-id}", d.Config.SubscriptionID, -1)
	roleDefinitionId = strings.Replace(roleDefinitionId, "{role-definition-id}", AzureActiveDirectoryContributorRoleId, -1)

	roleAssignmentReq := AdRoleAssignment{
		PrincipalID:      appProperties.ServicePrincipalObjectID,
		RoleDefinitionID: roleDefinitionId,
	}

	p_role := p
	p_role["role-assignment-name"] = "azkube_deployer_role"
	p_role["subscription-id"] = d.Config.SubscriptionID
	p_role["resource-group-name"] = d.Config.ResourceGroup
	q_role := map[string]interface{}{"api-version": AzureRoleManagementApiVersion}

	// TODO - Abandoning this until we find out if Authorization in azure/arm should include a client for this purpose
	// TODO(colemickens): Update azure-sdk-for-go to use new Authorization clients

	req, err = autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(AzureManagementBaseURL),
		autorest.WithPath("/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/Microsoft.Authorization/roleAssignments/{role-assignment-name}"),
		autorest.WithPathParameters(p_role),
		autorest.WithQueryParameters(q_role),
		autorest.WithJSON(roleAssignmentReq))
	if err != nil {
		return nil, err
	}

	resp, err = d.AdClient.Send(req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	return appProperties, nil
}

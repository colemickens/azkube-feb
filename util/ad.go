package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	"github.com/Azure/azure-sdk-for-go/arm/authorization"
	"golang.org/x/crypto/pkcs12"
)

const (
	AzureAdScope      = "https://graph.windows.net/"
	AzureAdApiVersion = "1.6"
	AzureAdBaseURL    = "https://graph.windows.net/{tenant-id}"

	AzureRoleManagementApiVersion = "2015-07-01"
	AzureManagementBaseURL        = "https://management.azure.com/{tenant-id}"

	AzureAdRoleReferenceTemplate = "/subscriptions/{subscription-id}/providers/Microsoft.Authorization/roleDefinitions/{role-definition-id}"
	AzureAdContributorRoleId     = "b24988ac-6180-42a0-ab88-20f7382dd24c"

	ServicePrincipalKeySize = 4096

	AzurePropagationWaitDelay = time.Second * 30 // TODO(Colemickens): poll instead of dumb sleep
)

type AdClient struct {
	autorest.Client
}

type AdApplication struct {
	ApplicationID string `json:"appId,omitempty"`    // readonly
	ObjectID      string `json:"objectId,omitempty"` // readonly

	AvailableToOtherTenants bool              `json:"availableToOtherTenants"`
	DisplayName             string            `json:"displayName,omitempty"`
	Homepage                string            `json:"homepage,omitempty"`
	IdentifierURIs          []string          `json:"identifierUris,omitempty"`
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

	ApplicationID  string `json:"appId,omitempty"`
	AccountEnabled bool   `json:"accountEnabled,omitempty"`
	//	ServicePrincipalNames []string `json:"servicePrincipalNames,omitempty"`
}

type AdRoleAssignment struct {
	RoleDefinitionID string `json:"roleDefinitionId,omitempty"`
	PrincipalID      string `json:"principalId,omitempty"`
}

func (app *AppProperties) ServicePrincipalPkcs12() ([]byte, error) {
	privateKey, err := PemToPrivateKey(app.ServicePrincipalPrivateKeyPem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pem into private key")
	}

	certificate, err := PemToCertificate(app.ServicePrincipalCertificatePem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pem into certificate")
	}

	pfxData, err := pkcs12.Encode(privateKey, certificate, nil, "")
	if err != nil {
		return nil, err
	}
	return pfxData, nil
}

func CreateServicePrincipalSecrets(notBefore, notAfter time.Time) (certDerBytes []byte, privateKey *rsa.PrivateKey, err error) {
	privateKey, err = rsa.GenerateKey(rand.Reader, ServicePrincipalKeySize)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}
	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "Azkube",
			Organization: []string{"Azkube"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		//ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDerBytes, err = x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	return certDerBytes, privateKey, nil
}

func (a *AdClient) CreateApp(common CommonProperties, appName, appURL string) (*AppProperties, error) {
	app := &AppProperties{}

	app.Name = appName
	app.IdentifierURL = appURL

	notBefore := time.Now()
	notAfter := time.Now().Add(5 * 365 * 24 * time.Hour)
	notAfter = time.Now().Add(10000 * 24 * time.Hour)

	cert, privateKey, err := CreateServicePrincipalSecrets(notBefore, notAfter)
	if err != nil {
		return nil, err
	}
	privateKeyPem := PrivateKeyToPem(privateKey)
	certificatePem := CertificateToPem(cert)

	app.ServicePrincipalPrivateKeyPem = string(privateKeyPem)
	app.ServicePrincipalCertificatePem = string(certificatePem)

	certificateDataParts := strings.Split(app.ServicePrincipalCertificatePem, "\n")
	certificateData := strings.Join(certificateDataParts[1:len(certificateDataParts)-2], "\n")

	startDate := notBefore.Format(time.RFC3339)
	endDate := notAfter.Format(time.RFC3339)

	////////////////////////////////////////////////////////////////////////////////////
	// create application
	applicationReq := AdApplication{
		AvailableToOtherTenants: false,
		DisplayName:             appName,
		Homepage:                appURL,
		IdentifierURIs:          []string{appURL},
		KeyCredentials: []AdKeyCredential{
			AdKeyCredential{
				KeyId:     uuid.New(),
				Type:      "AsymmetricX509Cert",
				Usage:     "Verify",
				StartDate: startDate,
				EndDate:   endDate,
				Value:     certificateData,
			},
		},
	}

	p := map[string]interface{}{"tenant-id": common.TenantID}
	q := map[string]interface{}{"api-version": AzureAdApiVersion}

	req, err := autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(AzureAdBaseURL),
		autorest.WithPath("applications"),
		autorest.WithPathParameters(p),
		autorest.WithQueryParameters(q),
		autorest.WithJSON(applicationReq))

	log.Println(req)

	if err != nil {
		return nil, err
	}

	resp, err := a.Send(req, http.StatusCreated)
	if err != nil {
		return nil, err
	}

	var applicationResp AdApplication
	err = autorest.Respond(resp, autorest.ByUnmarshallingJSON(&applicationResp))
	if err != nil {
		return nil, err
	}

	app.ApplicationID = applicationResp.ApplicationID

	time.Sleep(AzurePropagationWaitDelay)

	////////////////////////////////////////////////////////////////////////////////////
	// create service principal
	servicePrincipalReq := AdServicePrincipal{
		ApplicationID:  app.ApplicationID,
		AccountEnabled: true,
		//ServicePrincipalNames: []string{appURL},
	}

	req, err = autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(AzureAdBaseURL),
		autorest.WithPath("servicePrincipals"),
		autorest.WithPathParameters(p),
		autorest.WithQueryParameters(q),
		autorest.WithJSON(servicePrincipalReq))
	if err != nil {
		return nil, err
	}

	resp, err = a.Send(req, http.StatusCreated)
	if err != nil {
		return nil, err
	}

	var servicePrincipalResp AdServicePrincipal
	err = autorest.Respond(resp, autorest.ByUnmarshallingJSON(&servicePrincipalResp))
	if err != nil {
		return nil, err
	}

	app.ServicePrincipalObjectID = servicePrincipalResp.ObjectID

	time.Sleep(AzurePropagationWaitDelay)

	return app, nil
}

func (d *Deployer) CreateRoleAssignment(common CommonProperties, principalID string) error {
	roleAssignmentName := uuid.New()

	roleDefinitionId := strings.Replace(AzureAdRoleReferenceTemplate, "{subscription-id}", common.SubscriptionID, -1)
	roleDefinitionId = strings.Replace(roleDefinitionId, "{role-definition-id}", AzureAdContributorRoleId, -1)

	scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", common.SubscriptionID, common.ResourceGroup)

	roleAssignmentParameters := authorization.RoleAssignmentCreateParameters{
		Properties: &authorization.RoleAssignmentProperties{
			RoleDefinitionID: &roleDefinitionId,
			PrincipalID:      &principalID,
		},
	}

	_, err := d.RoleAssignmentsClient.Create(
		scope,
		roleAssignmentName,
		roleAssignmentParameters,
	)
	if err != nil {
		return err
	}

	return nil
}

#!/usr/bin/env bash

# this creates a Contributor service principal

# this creates a special Owner service principal, but also 
# grants it the Company Administrator role, which is seemingly necessary
# for this SP to be able to provision new applications

# source:

# https://support.microsoft.com/en-us/kb/3004133
# https://stackoverflow.com/questions/16093584/windows-azure-graph-api-to-add-an-application

# create a private key for the SP
openssl genrsa -out sp.key 4096

# create a certificate for the SP
openssl req -x509 -new -nodes -key sp.key -sub "/CN=sp" -days 10000 -out sp.crt

# pull the key value out into the format that Azure Xplat CLI expects
key_value="$(tail -n+2 "sp.crt" | head -n-1)"

# create the Azure AD application
app_output="$(azure ad app create \
	--name "AzureKubernetesAdministrator" \
	--home-page "http://azurekubernetesadministrator" \
	--identifier-urls "http://azurekubernetesadministrator" \
	--key-usage "Verify" \
	--end-date "2020-01-01" \
	--key-value "${key_value}")"

# capture the application id
app_id="$(echo "${app_output}" | grep "Application Id" | awk '{ print $4 }')"

# create the service principal
azure ad sp create "${app_id}"

# assign the Owner role to the SP
azure role assignment create \
	--subscription "${subscription_id}" \
	--roleName "Owner" \
	--spn "http://azurekubernetesadministrator"

## ### INSERT STACK OVER FLOW-Y STUFF HERE

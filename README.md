# azkube

## Overview
Tool used to deploy and bootstrap a Kubernetes cluster in Azure.


## Ideas

Use confd for managing upgrades?


## Requirements

1. Install either [Widows Azure PowerShell Tools](https://github.com/Azure/azure-powershell) or the [azure-xplat-cli](https://github.com/Azure/azure-xplat-cli).

2. Docker


## Prepare

You will need to create a ServicePrincipal with the Owner role in your subscription. While these instructions are written with the [azure-xplat-cli](https://github.com/Azure/azure-xplat-cli) in mind, the Service Principal can be created with the PowerShell tooling as well.

This Owner account is used by `azkube` to provision a new ServicePrincipal with Contributor access for each cluster.
It is also used to do the initial deployment.

This Contributor account is used by the Kubernetes cluster to provision and manage Azure resources as needed.

1. Create a new Azure Active Directory account:

	`azure ad app create --name "azkube" --password "somepassword"`

2. Create a service principal for the application

	`azure sp create {app-id}`

	Replace `{app-id}` with the ApplicationID reported by the previous step.

3. Assign the Owner role to the new Service Principal.

	`azure ad role assignment create --spn "http://azkube" --roleName "Owner" --subscription "{your-subscription-id}`



## Deploy

`docker run -it colemickens:azkube /azkube deploy -config config.json`

Wait for the deployment to finish!

## Todo
- ?

## Future Ideas
- Terraform?
- Confd for upgrading cluster?

## Issues brought up:

1. MASTER_IP could be fragile? (I'm not sure it is in Azure...)

2. Porting to ubuntu could be difficult because of docker-bootstrap instance


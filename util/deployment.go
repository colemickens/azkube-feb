package util

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/azure-sdk-for-go/arm/resources"
)

var cachedDeploymentsClient *resources.DeploymentsClient = nil

func getDeploymentsClient(config DeployConfigOut) (client *resources.DeploymentsClient, err error) {
	if cachedDeploymentsClient != nil {
		return cachedDeploymentsClient, nil
	}

	dpc := resources.NewDeploymentsClient(config.SubscriptionID)
	dpc.Authorizer, err = GetAuthorizer(config, azure.AzureResourceManagerScope)
	if err != nil {
		return nil, err
	}

	cachedDeploymentsClient = &dpc
	return cachedDeploymentsClient, nil
}

func DoDeployment(config DeployConfigOut, name string, template, params map[string]interface{}, waitDeployment bool) (response *resources.DeploymentExtended, err error) {
	dpc, err := getDeploymentsClient(config)
	if err != nil {
		panic(err)
	}

	deployment := resources.Deployment{
		Properties: &resources.DeploymentProperties{
			Template:   &template,
			Parameters: &params,
			Mode:       resources.Incremental,
		},
	}

	deploymentResponse, err := dpc.CreateOrUpdate(
		config.ResourceGroup,
		config.ResourceGroup+"-"+name+"-deploy",
		deployment)
	if err != nil {
		panic(err)
	}

	if waitDeployment {
		deploymentName := *deploymentResponse.Name
		response, err = WaitDeployment(config, deploymentName)
		return response, err
	}

	return &deploymentResponse, err
}

func WaitDeployment(config DeployConfigOut, deploymentName string) (*resources.DeploymentExtended, error) {
	dpc, err := getDeploymentsClient(config)
	if err != nil {
		return nil, err
	}

	var response resources.DeploymentExtended
	for {
		response, err = dpc.Get(config.ResourceGroup, deploymentName)
		if err != nil {
			return &response, err
		}

		state := response.Properties.ProvisioningState

		if state == nil {
			continue
		}

		if *state == "Succeeded" {
			log.Println(deploymentName + " deployment succeeded!")
			return &response, nil
		} else if *state == "Failed" {
			return &response, fmt.Errorf(deploymentName + " deployment failed!")
		} else {
			return &response, fmt.Errorf(deploymentName+" deployment went into unknown state: %s", *state)
		}
	}

	return &response, nil
}

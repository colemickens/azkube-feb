package util

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

func (d *Deployer) DoDeployment(resourceGroupName, deploymentName string, template map[string]interface{}, parameters map[string]interface{}, waitDeployment bool) (response *resources.DeploymentExtended, err error) {
	deployment := resources.Deployment{
		Properties: &resources.DeploymentProperties{
			Template:   &template,
			Parameters: &parameters,
			Mode:       resources.Incremental,
		},
	}

	deploymentResponse, err := d.DeploymentsClient.CreateOrUpdate(
		resourceGroupName,
		resourceGroupName+"-"+deploymentName+"-deploy",
		deployment)
	if err != nil {
		panic(err)
	}

	if waitDeployment {
		deploymentName := *deploymentResponse.Name
		response, err = d.WaitDeployment(resourceGroupName, deploymentName)
		return response, err
	}

	return &deploymentResponse, err
}

func (d *Deployer) WaitDeployment(resourceGroup, deploymentName string) (*resources.DeploymentExtended, error) {
	var err error
	var response resources.DeploymentExtended
	for {
		time.Sleep(5 * time.Second)

		response, err = d.DeploymentsClient.Get(resourceGroup, deploymentName)
		if err != nil {
			return &response, err
		}

		state := response.Properties.ProvisioningState

		if state == nil || *state == "Accepted" || *state == "Running" {
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

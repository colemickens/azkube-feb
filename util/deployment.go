package util

import (
	"fmt"
	"log"

	// "github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/azure-sdk-for-go/arm/resources"
)

func (d *Deployer) DoDeployment(commonProperties CommonProperties, name string, template map[string]interface{}, waitDeployment bool) (response *resources.DeploymentExtended, err error) {
	deployment := resources.Deployment{
		Properties: &resources.DeploymentProperties{
			Template: &template,
			Mode:     resources.Incremental,
		},
	}

	deploymentResponse, err := d.DeploymentsClient.CreateOrUpdate(
		commonProperties.ResourceGroup,
		commonProperties.ResourceGroup+"-"+name+"-deploy",
		deployment)
	if err != nil {
		panic(err)
	}

	if waitDeployment {
		// TODO(colemickens): assert this name is the same?
		// here we use the returned deploymentName but in groups we use original resGroup name?
		deploymentName := *deploymentResponse.Name
		response, err = d.WaitDeployment(deploymentName)
		return response, err
	}

	return &deploymentResponse, err
}

func (d *Deployer) WaitDeployment(deploymentName string) (*resources.DeploymentExtended, error) {
	var err error
	var response resources.DeploymentExtended
	for {
		response, err = d.DeploymentsClient.Get(deploymentName, deploymentName)
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

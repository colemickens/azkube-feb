package util

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/azure-sdk-for-go/arm/resources"
)

func EnsureResourceGroup(config DeploymentProperties, waitDeployment bool) (resourceGroup *resources.ResourceGroup, err error) {
	groupsClient := resources.NewGroupsClient(config.SubscriptionID)
	groupsClient.Authorizer, err = GetAuthorizer(config, azure.AzureResourceManagerScope)
	if err != nil {
		return nil, err
	}

	response, err := groupsClient.CreateOrUpdate(config.ResourceGroup, resources.ResourceGroup{
		Name:     &config.ResourceGroup,
		Location: &config.Location,
	})
	if err != nil {
		return &response, err
	}

	if waitDeployment {
		return WaitResourceGroup(config)
	}

	return &response, nil
}

func WaitResourceGroup(config DeploymentProperties) (resourceGroup *resources.ResourceGroup, err error) {
	groupsClient := resources.NewGroupsClient(config.SubscriptionID)
	groupsClient.Authorizer, err = GetAuthorizer(config, azure.AzureResourceManagerScope)
	if err != nil {
		return nil, err
	}

	var response resources.ResourceGroup
	for {
		response, err = groupsClient.Get(config.ResourceGroup)
		if err != nil {
			return &response, err
		}

		state := response.Properties.ProvisioningState

		if state == nil {
			continue
		}

		if *state == "Succeeded" {
			log.Println("group deployment succeeded!")
			return &response, nil
		} else if *state == "Failed" {
			return &response, fmt.Errorf("group deployment failed!")
		} else {
			return &response, fmt.Errorf("group deployment went into unknown state: %s", *state)
		}
	}

	return nil, nil
}

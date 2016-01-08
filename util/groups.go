package util

import (
	"fmt"
	"log"

	//"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

func (d *Deployer) EnsureResourceGroup(name, location string, waitDeployment bool) (resourceGroup *resources.ResourceGroup, err error) {
	response, err := d.GroupsClient.CreateOrUpdate(name, resources.ResourceGroup{
		Name:     &name,
		Location: &location,
	})
	if err != nil {
		return &response, err
	}

	if waitDeployment {
		return d.WaitResourceGroup(name)
	}

	return &response, nil
}

func (d *Deployer) WaitResourceGroup(groupName string) (resourceGroup *resources.ResourceGroup, err error) {
	var response resources.ResourceGroup
	for {
		response, err = d.GroupsClient.Get(groupName)
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

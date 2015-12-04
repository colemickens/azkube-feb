package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var variableRegex regexp.Regexp
var parameterRegex regexp.Regexp

func init() {
	variableRegex = *regexp.MustCompile("")
	parameterRegex = *regexp.MustCompile("")
}

func formatCloudConfig(filepath string) (string, error) {
	cloudConfigBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	cloudConfig := string(cloudConfigBytes)

	if strings.ContainsAny(cloudConfig, "'") {
		panic("can't have single quotes in the cloud config files")
	}

	data, err := json.Marshal(&cloudConfig)
	if err != nil {
		panic(err)
	}

	parameterRegex.ReplaceAllString(cloudConfig, "', parameters('$1'), '")
	variableRegex.ReplaceAllString(cloudConfig, "', variables('$1'), '")

	result := "[concat('" + string(data) + "')]"

	return result, nil
}

func InsertCloudConfig(c string) string {
	masterCloudConfig, err := formatCloudConfig("templates/coreos/master-cloudconfig.in.yaml")
	if err != nil {
		panic(err)
	}
	nodeCloudConfig, err := formatCloudConfig("templates/coreos/node-cloudconfig.in.yaml")
	if err != nil {
		panic(err)
	}

	template := c
	template = strings.Replace(template, "{{MASTER_CLOUDCONFIG}}", masterCloudConfig, -1)
	template = strings.Replace(template, "{{NODE_CLOUDCONFIG}}", nodeCloudConfig, -1)

	return template
}

func LoadAndFormat(name string, config DeployConfigOut, preprocessor func(c string) string) (template, parameters map[string]interface{}, err error) {
	myriadTemplateBytes, err := ioutil.ReadFile("templates/coreos/" + name + ".in.json")
	if err != nil {
		return nil, nil, err
	}

	myriadTemplate := string(myriadTemplateBytes)
	if preprocessor != nil {
		myriadTemplate = preprocessor(myriadTemplate)
	}

	myriadParamtersBytes, err := ioutil.ReadFile("templates/coreos/" + name + "-paramters.in.json")
	if err != nil {
		return nil, nil, err
	}

	params := string(myriadParamtersBytes)
	params = strings.Replace(params, "{{MASTER_VM_SIZE}}", config.MasterVmSize, -1)
	params = strings.Replace(params, "{{NODE_VM_SIZE}}", config.NodeVmSize, -1)
	params = strings.Replace(params, "{{NODE_COUNT}}", fmt.Sprintf("%d", config.NodeCount), -1)
	params = strings.Replace(params, "{{USERNAME}}", config.Username, -1)
	params = strings.Replace(params, "{{MASTER_FQDN}}", config.MasterFqdn, -1)
	params = strings.Replace(params, "{{SSH_PUBLIC_KEY_DATA}}", config.SshPublicKeyData, -1)
	params = strings.Replace(params, "{{TENANT_ID}}", config.TenantID, -1)
	params = strings.Replace(params, "{{APP_URL}}", config.AppURL, -1)
	params = strings.Replace(params, "{{DEPLOYER_OBJECT_ID}}", config.DeployerObjectID, -1)
	params = strings.Replace(params, "{{SERVICE_PRINCIPAL_OBJECT_ID}}", config.ServicePrincipalObjectID, -1)
	params = strings.Replace(params, "{{SERVICE_PRINCIPAL_SECRET_URL}}", config.ServicePrincipalSecretURL, -1)
	params = strings.Replace(params, "{{VAULT_NAME}}", config.VaultName, -1)

	err = json.Unmarshal([]byte(myriadTemplate), &template)
	err = json.Unmarshal([]byte(params), &parameters)

	return template, parameters, nil
}

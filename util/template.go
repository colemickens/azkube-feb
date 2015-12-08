package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"text/template"
)

var (
	masterCloudConfigTemplate *template.Template
	nodeCloudConfigTemplate   *template.Template
	vaultTemplate             *template.Template
	myriadTemplate            *template.Template
	scaleTemplate             *template.Template
)

func init() {
	x := func(y, z string) *template.Template {
		bytes, err := ioutil.ReadFile(y)
		if err != nil {
			panic(err)
		}
		contents := string(bytes)

		return template.Must(template.New(z).Parse(contents))
	}

	masterCloudConfigTemplate = x("templates/coreos/master.cloudconfig.in.yaml", "masterCloudConfigTemplate")
	nodeCloudConfigTemplate = x("templates/coreos/node.cloudconfig.in.yaml", "nodeCloudConfigTemplate")
	vaultTemplate = x("templates/vault/vault.in.json", "vaultTemplate")
	myriadTemplate = x("templates/coreos/myriad.in.json", "myriadTemplate")
	scaleTemplate = x("templates/scale/scale.in.json", "scaleTemplate")
}

func formatCloudConfig(filepath string) (string, error) {
	cloudConfigBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	cloudConfig := string(cloudConfigBytes)

	data, err := json.Marshal(&cloudConfig)
	if err != nil {
		panic(err)
	}

	return string(data), nil
}

func CreateMyriadTemplate(config DeploymentProperties) map[string]interface{} {
	var masterBuf bytes.Buffer
	err := masterCloudConfigTemplate.Execute(&masterBuf, config)
	if err != nil {
		panic(err)
	}
	masterBytes, err := json.Marshal(masterBuf.String())
	if err != nil {
		panic(err)
	}
	config.CloudConfig.Master = string(masterBytes)

	var nodeBuf bytes.Buffer
	err = nodeCloudConfigTemplate.Execute(&nodeBuf, config)
	if err != nil {
		panic(err)
	}
	nodeBytes, err := json.Marshal(nodeBuf.String())
	if err != nil {
		panic(err)
	}
	config.CloudConfig.Node = string(nodeBytes)

	var myriadBuf bytes.Buffer
	var myriadMap map[string]interface{}
	err = myriadTemplate.Execute(&myriadBuf, config)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(myriadBuf.Bytes(), &myriadMap)
	if err != nil {
		panic(err)
	}

	return myriadMap
}

func CreateVaultTemplate(config DeploymentProperties) map[string]interface{} {
	vaultBuf := bytes.Buffer{}
	err := vaultTemplate.Execute(&vaultBuf, config)
	if err != nil {
		panic(err)
	}

	var vaultTemplate map[string]interface{}
	err = json.Unmarshal(vaultBuf.Bytes(), &vaultTemplate)
	if err != nil {
		panic(err)
	}

	return vaultTemplate
}

func CreateScaleTemplate(config DeploymentConfig) map[string]interface{} {
	scaleBuf := bytes.Buffer{}
	err := scaleTemplate.Execute(&scaleBuf, config)
	if err != nil {
		panic(err)
	}

	var scaleTemplate map[string]interface{}
	err = json.Unmarshal(scaleBuf.Bytes(), &scaleTemplate)
	if err != nil {
		panic(err)
	}

	return scaleTemplate
}

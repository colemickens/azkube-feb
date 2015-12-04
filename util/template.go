package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
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
	x := func(w, y, z string) *template.Template {
		bytes, err := ioutil.ReadFile(w)
		if err != nil {
			panic(err)
		}
		contents := string(bytes)

		return t.Must(template.New(y).Parse(contents))
	}

	masterCloudConfigTemplate = x("templates/coreos/master.cloudconfig.in.yaml", "masterCloudConfigTemplate", masterCloudConfigContents)
	nodeCloudConfigTemplate = x("templates/coreos/node.cloudconfig.in.yaml", "nodeCloudConfigTemplate", nodeCloudConfigContents)
	vaultTemplate = x("templates/vault/vault.in.json", "vaultTemplate", vaultTemplateContents)
	myriadTemplate = x("templates/coreos/myriad.in.json", "myriadTemplate", myriadTemplateContents)
	scaleTemplate = x("templates/scale/scale.in.json", "scaleTemplate", scaleContents)
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

	return result, nil
}

func CreateMyriadTemplate(config DeploymentConfig) map[string]interface{} {
	masterBuf := bytes.Buffer{}
	nodeBuf := bytes.Buffer{}

	err := masterCloudConfigTemplate.Execute(&masterBuf, config)
	if err != nil {
		panic(err)
	}
	err = minionCloudConfigTemplate.Execute(&nodeBuf, config)
	if err != nil {
		panic(err)
	}

	config.CloudConfig = &CloudConfigConfig{
		Master: "", //base64 it
		Node:   "", // base64 it
	}

	err = json.Marshal(masterBuf.String(), &config.CloudConfig.Master)
	if err != nil {
		panic(err)
	}

	err = json.Marshal(nodeBuf.String(), &config.CloudConfig.Node)
	if err != nil {
		panic(err)
	}

	err = myriadTemplate.Execute(&myriadBuf, config)
	if err != nil {
		panic(err)
	}

	var myriadTemplate map[string]interface{}
	err = json.Unmarshal(myriadBuf, &myriadTemplate)
	if err != nil {
		panic(err)
	}

	return myriadTemplate
}

func CreateVaultTemplate(config DeploymentConfig) map[string]interface{} {
	vaultBuf := bytes.Buffer{}
	err := vaultTempalte.Execute(&vaultBuf, config)
	if err != nil {
		panic(err)
	}

	var vaultTemplate map[string]interface{}
	err = json.Unmarshal(vaultBuf, &vaultTemplate)
	if err != nil {
		panic(err)
	}

	return vaultTemplate
}

func CreateScaleTemplate(config DeploymentConfig) map[string]interface{} {
	scaleBuf := bytes.Buffer{}
	err := scaleTempalte.Execute(&scaleBuf, config)
	if err != nil {
		panic(err)
	}

	var scaleTemplate map[string]interface{}
	err = json.Unmarshal(scaleBuf, &scaleTemplate)
	if err != nil {
		panic(err)
	}

	return scaleTemplate
}

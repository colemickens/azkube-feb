package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"text/template"
)

var (
	MasterCloudConfigTemplate *template.Template
	NodeCloudConfigTemplate   *template.Template
	VaultTemplate             *template.Template
	MyriadTemplate            *template.Template
	ScaleTemplate             *template.Template
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

	MasterCloudConfigTemplate = x("templates/coreos/master.cloudconfig.in.yaml", "masterCloudConfigTemplate")
	NodeCloudConfigTemplate = x("templates/coreos/node.cloudconfig.in.yaml", "nodeCloudConfigTemplate")
	VaultTemplate = x("templates/vault/vault.in.json", "vaultTemplate")
	MyriadTemplate = x("templates/coreos/myriad.in.json", "myriadTemplate")
	ScaleTemplate = x("templates/scale/scale.in.json", "scaleTemplate")
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

func (d *Deployer) LoadMyriadCloudConfigs() (myriadConfig *MyriadConfig, err error) {
	var masterBuf bytes.Buffer
	err = MasterCloudConfigTemplate.Execute(&masterBuf, d.State)
	if err != nil {
		return nil, err
	}
	masterBytes, err := json.Marshal(masterBuf.String())
	if err != nil {
		panic(err)
	}
	myriadConfig.MasterCloudConfig = string(masterBytes)

	var nodeBuf bytes.Buffer
	err = NodeCloudConfigTemplate.Execute(&nodeBuf, d.State)
	if err != nil {
		panic(err)
	}
	nodeBytes, err := json.Marshal(nodeBuf.String())
	if err != nil {
		panic(err)
	}
	myriadConfig.NodeCloudConfig = string(nodeBytes)

	return myriadConfig, nil
}

func (d *Deployer) PopulateTemplate(t *template.Template) (template map[string]interface{}, err error) {
	var myriadBuf bytes.Buffer
	var myriadMap map[string]interface{}

	err = t.Execute(&myriadBuf, d.State)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(myriadBuf.Bytes(), &myriadMap)
	if err != nil {
		return nil, err
	}

	return myriadMap, nil
}

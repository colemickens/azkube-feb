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

	MasterCloudConfigTemplate = x("templates/coreos/master-cloudconfig.in.yml", "masterCloudConfigTemplate")
	NodeCloudConfigTemplate = x("templates/coreos/node-cloudconfig.in.yml", "nodeCloudConfigTemplate")
	VaultTemplate = x("templates/vault/vault.in.json", "vaultTemplate")
	MyriadTemplate = x("templates/coreos/myriad.in.json", "myriadTemplate")
	ScaleTemplate = x("templates/scale/scale.in.json", "scaleTemplate")
}

func PopulateAndFlattenTemplate(t *template.Template, state interface{}) (string, error) {
	filled, err := PopulateTemplate(t, state)
	if err != nil {
		return "", nil
	}

	data, err := json.Marshal(&filled)
	if err != nil {
		return "", nil
	}

	return string(data), nil
}

func PopulateTemplate(t *template.Template, state interface{}) (template map[string]interface{}, err error) {
	var myriadBuf bytes.Buffer
	var myriadMap map[string]interface{}

	err = t.Execute(&myriadBuf, state)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(myriadBuf.Bytes(), &myriadMap)
	if err != nil {
		return nil, err
	}

	return myriadMap, nil
}

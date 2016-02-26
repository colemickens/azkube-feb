package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"
	"text/template"

	log "github.com/Sirupsen/logrus"
)

var (
	myriadTemplate           *template.Template
	myriadParametersTemplate *template.Template
	masterScript             string
	nodeScript               string

	variableRegex  = regexp.MustCompile(`[[[(a-zA-Z)]]]`)
	parameterRegex = regexp.MustCompile(`{{{(a-zA-Z)}}}`)
)

func init() {
	mustRead := func(f string) string {
		bytes, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}
		contents := string(bytes)
		return contents
	}

	masterScript = mustRead("templates/coreos/master-cloudconfig.in.yml")
	nodeScript = mustRead("templates/coreos/node-cloudconfig.in.yml")

	var err error
	myriadTemplate, err =
		template.New("myriadTemplate").Parse(mustRead("templates/coreos/azdeploy.in.json"))
	myriadParametersTemplate, err =
		template.New("myriadParameters").Parse(mustRead("templates/coreos/parameters.in.json"))
	if err != nil {
		panic(err)
	}
}

func ProduceTemplateAndParameters(flavorArgs FlavorArguments) (template, parameters map[string]interface{}, err error) {
	log.Info("template: preparing master script")
	masterScript, _ := prepareScript(masterScript)

	log.Info("template: preparing node script")
	nodeScript, _ := prepareScript(nodeScript)

	// TODO: consider, does this "list" become part of the flavor interface?
	myriadParameters, err := populateTemplate(
		myriadParametersTemplate,
		flavorArgs)
	if err != nil {
		return nil, nil, err
	}

	myriadTemplate, err := populateTemplate(
		myriadTemplate,
		struct{ MasterScript, NodeScript string }{masterScript, nodeScript})
	if err != nil {
		return nil, nil, err
	}

	// TODO: persist this to disk

	// smoosh them into maps so we can return

	return myriadParameters, myriadTemplate, nil
}

func prepareScript(script string) (string, error) {
	if strings.Contains(script, "'") {
		panic("NO SINGLE QUOTES") // TODO(colemick): nicer...
	}

	script = variableRegex.ReplaceAllString(script, `', variables('$1'), '`)
	script = parameterRegex.ReplaceAllString(script, `', parameters('$1'), '`)

	script = `[concat('` + script + `')]`

	bytes, err := json.Marshal(script)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func populateTemplate(t *template.Template, state interface{}) (template map[string]interface{}, err error) {
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

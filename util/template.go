package util

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	if err != nil {
		panic(err)
	}

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

	log.Info("template: populating myriad parameters")
	// TODO: consider, does this "list" become part of the flavor interface?
	myriadParameters, err := populateTemplate(
		myriadParametersTemplate,
		flavorArgs)
	if err != nil {
		return nil, nil, err
	}

	log.Info("template: populating myriad template")
	myriadTemplate, err := populateTemplate(
		myriadTemplate,
		struct{ MasterScript, NodeScript string }{masterScript, nodeScript})
	if err != nil {
		return nil, nil, err
	}

	// TODO: persist this to disk
	log.Info("template: all done")

	return myriadTemplate, myriadParameters, nil
}

func prepareScript(script string) (string, error) {
	if strings.Contains(script, "'") {
		panic("NO SINGLE QUOTES") // TODO(colemick): nicer...
	}

	script = variableRegex.ReplaceAllString(script, `', variables('$1'), '`)
	script = parameterRegex.ReplaceAllString(script, `', parameters('$1'), '`)

	// this is commented out because we escape in the template syntax
	// but the single quotes here need to NOT be escaped
	//script = `[concat('` + script + `')]`

	/*bytes, err := json.Marshal(script)
	if err != nil {
		return "", err
	}

	return string(bytes), nil*/
	return script, nil
}

func populateTemplate(t *template.Template, state interface{}) (template map[string]interface{}, err error) {
	var myriadBuf bytes.Buffer
	var myriadMap map[string]interface{}

	if t == nil {
		log.Fatal("Nil!!!")
	}

	err = t.Execute(&myriadBuf, state)
	if err != nil {
		return nil, fmt.Errorf("template: failed to execute template: %q", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("/home/cole/%s.txt", t.Name()), myriadBuf.Bytes(), 0666)
	if err != nil {
		panic(err)
	}
	log.Info("looggeeeeeeeeeeeed it")

	err = json.Unmarshal(myriadBuf.Bytes(), &myriadMap)
	if err != nil {
		return nil, fmt.Errorf("template: failed to unmarshal into map: %q", err)
	}

	return myriadMap, nil
}

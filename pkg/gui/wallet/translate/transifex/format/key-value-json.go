package format

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type KeyValueJson struct{}

func (f *KeyValueJson) Init(initParams map[string]interface{}) {
	// nothing to do at the moment
}
func (f KeyValueJson) Ext() string { return "json" }

func (f KeyValueJson) Clean(content []byte) ([]byte, string, error) {
	var data map[string]string
	jsonErr := json.Unmarshal(content, &data)
	if jsonErr != nil {
		return nil, "", fmt.Errorf("Not valid json: %s", jsonErr)
	}
	for key, value := range data {
		if key == "" {
			delete(data, key)
		}
		if value == "" {
			data[key] = " "
		}
		content, jsonErr = json.Marshal(data)
		if jsonErr != nil {
			panic("An error occurred when encoding json after updating json so that transifex can use it")
		}
	}

	return content, "KEYVALUEJSON", nil

}
func (f KeyValueJson) Write(rootDir, langCode, srcLang, filename, translation string, fileLocator FileLocator) error {
	path := fileLocator.Find(rootDir, langCode, filename, "json")
	fmt.Println("Updating translations file: " + path)
	return ioutil.WriteFile(path, []byte(translation), 0644)
}

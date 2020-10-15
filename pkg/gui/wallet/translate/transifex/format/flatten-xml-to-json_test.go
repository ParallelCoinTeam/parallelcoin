package format

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"
	tu "testutil"
	"text/template"
)

const (
	xmlData            = `<a><b at="atv" at2="atv2"><l>value1</l><d at="atv" at2="atv2">dv</d></b><c>value2</c><c>value3</c></a>`
	expectedUpdatedXml = `<a><b at="atv" at2="atv2"><l>newvalue1</l><d at="atv" at2="atv2">newdv</d></b><c>newvalue2</c><c>newvalue3</c></a>`
)

func Test_Clean(t *testing.T) {
	format := FlattenXmlToJson{}
	cleaned, i18, err := format.Clean([]byte(xmlData))

	if err != nil {
		t.Error(err)
	}
	if i18 != "KEYVALUEJSON" {
		t.Error("Expecte i18 to be KEYVALUEJSON")
	}
	var prop map[string]string
	json.Unmarshal(cleaned, &prop)

	if len(prop) != 4 {
		t.Errorf("Expected 4 items but found %d in %s", len(prop), prop)
	}

	for key, item := range prop {
		switch key {
		case "b[at=atv and at2=atv2] l":
			tu.AssertEquals("itemValue0", "value1", item, t)
		case "b[at=atv and at2=atv2] d[at=atv and at2=atv2]":
			tu.AssertEquals("itemValue1", "dv", item, t)
		case "c":
			tu.AssertEquals("itemValue2", "value2", item, t)
		case "c<2>":
			tu.AssertEquals("itemValue2", "value3", item, t)
		default:
			t.Errorf("Unexpected key: %s in %s", key, prop)
		}
	}
}

type templateFiller struct {
	Key   string
	Value string
}

func Test_Clean_Write_KeyText(t *testing.T) {
	format := FlattenXmlToJson{}
	format.Init(map[string]interface{}{"key": "translationKey"})

	if format.key != "translationKey" {
		t.Error("Incorrect formatKey: " + format.key)
	}

	keyXmlTemplateText := `
<strings>
	<list id="list1">
		<el>
			<translationKey>{{.Key}} 1</translationKey>
			<value>{{.Value}} 1<value>
		</el>
		<el>
			<translationKey>{{.Key}} 2</translationKey>
			<value>{{.Value}} 2<value>
		</el>
	</list>
	<list id="list2"/>
		<el>
			<translationKey>{{.Key}} 3</translationKey>
			<value>{{.Value}} 3<value>
		</el>
	</list>	
</strings>`
	tmpl, err := template.New("xmlTemplate").Parse(keyXmlTemplateText)

	xmlData := new(bytes.Buffer)
	tmpl.Execute(xmlData, templateFiller{Key: "key", Value: "value"})
	cleaned, i18, err := format.Clean(xmlData.Bytes())

	if err != nil {
		t.Error(err)
	}
	if i18 != "KEYVALUEJSON" {
		t.Error("Expecte i18 to be KEYVALUEJSON")
	}
	var prop map[string]string
	json.Unmarshal(cleaned, &prop)

	if len(prop) != 3 {
		t.Errorf("Expected 4 items but found %d in %s", len(prop), prop)
	}

	for key, item := range prop {
		switch key {
		case "list[id=list1] el <key 1>":
			tu.AssertEquals("", "value1", item, t)
		case "list[id=list1] el <key 2>":
			tu.AssertEquals("", "value2", item, t)
		case "list[id=list2] el <key 3>":
			tu.AssertEquals("", "value3", item, t)
		default:
			t.Errorf("Unexpected key: %s in %s", key, prop)
		}
	}

	beforeWriteData := new(bytes.Buffer)
	tmpl.Execute(beforeWriteData, templateFiller{Key: "before", Value: "before"})
	tmpDir, _ := ioutil.TempDir("", "test")
	fileToUpdate := filepath.Join(tmpDir, "en-name.xml")
	if err := ioutil.WriteFile(fileToUpdate, beforeWriteData.Bytes(), 644); err != nil {
		panic(err)
	}

	format.Write(tmpDir, "en", "fr", "name", string(cleaned), FileLocators["LANG-NAME"])

	updatedXmlData, _ := ioutil.ReadFile(filepath.Join(tmpDir, "fr-name.xml"))

	tu.AssertEquals("", xmlData.String(), string(updatedXmlData), t)
}

func Test_WriteUpdateExisting(t *testing.T) {
	translation := make(map[string]string)
	translation["b[at=atv and at2=atv2] l"] = "newvalue1"
	translation["b[at=atv and at2=atv2] d[at=atv and at2=atv2]"] = "newdv"
	translation["c"] = "newvalue2"
	translation["c<2>"] = "newvalue3"

	translationAsBytes, _ := json.Marshal(translation)

	tmpDir, _ := ioutil.TempDir("", "test")
	fileToUpdate := filepath.Join(tmpDir, "en-name.xml")
	if err := ioutil.WriteFile(fileToUpdate, []byte(xmlData), 644); err != nil {
		panic(err)
	}

	format := FlattenXmlToJson{}
	err := format.Write(tmpDir, "en", "en", "name", string(translationAsBytes), FileLocators["LANG-NAME"])
	if err != nil {
		t.Errorf("Unexpected error occurred. %s", err)
	}

	updatedXmlData, _ := ioutil.ReadFile(fileToUpdate)

	tu.AssertEquals("", expectedUpdatedXml, string(updatedXmlData), t)
}

func Test_WriteCreateNew(t *testing.T) {
	translation := make(map[string]string)
	translation["b[at=atv and at2=atv2] l"] = "newvalue1"
	translation["b[at=atv and at2=atv2] d[at=atv and at2=atv2]"] = "newdv"
	translation["c"] = "newvalue2"
	translation["c<2>"] = "newvalue3"

	translationAsBytes, _ := json.Marshal(translation)

	tmpDir, _ := ioutil.TempDir("", "test")
	fileToUpdate := filepath.Join(tmpDir, "en-name.xml")
	if err := ioutil.WriteFile(fileToUpdate, []byte(xmlData), 644); err != nil {
		panic(err)
	}

	format := FlattenXmlToJson{}
	err := format.Write(tmpDir, "fr", "en", "name", string(translationAsBytes), FileLocators["LANG-NAME"])
	if err != nil {
		t.Errorf("Unexpected error occurred. %s", err)
	}

	updatedXmlData, _ := ioutil.ReadFile(filepath.Join(tmpDir, "fr-name.xml"))

	tu.AssertEquals("", expectedUpdatedXml, string(updatedXmlData), t)
}

func Test_Write_MissingKey(t *testing.T) {
	translation := make(map[string]string)
	translation["a b[at=atv and at2=atv2] l"] = "newvalue1"

	translationAsBytes, _ := json.Marshal(translation)

	tmpDir, _ := ioutil.TempDir("", "test")
	fileToUpdate := filepath.Join(tmpDir, "en-name.xml")
	if err := ioutil.WriteFile(fileToUpdate, []byte(xmlData), 644); err != nil {
		panic(err)
	}

	err := FlattenXmlToJson{}.Write(tmpDir, "fr", "en", "name", string(translationAsBytes), FileLocators["LANG-NAME"])
	if err == nil {
		t.Errorf("Expected an error")
	}
}

package format

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	tu "testutil"
)

func Test_LangNameLocator_Find(t *testing.T) {
	l := FileLocators["LANG-NAME"]
	tu.AssertEquals("", filepath.Join("xyz", "fr-name.json"), l.Find("xyz", "fr", "name", "json"), t)
	tu.AssertEquals("", filepath.Join("xyz", "en-name.json"), l.Find("xyz", "en", "name", "json"), t)
	tu.AssertEquals("", filepath.Join("xyz", "abc", "en-name.xml"), l.Find("xyz/abc", "en", "name", "xml"), t)
	tu.AssertEquals("", "en-name.xml", l.Find("", "en", "name", "xml"), t)

	l = FileLocators["3-CHAR-LANG-NAME"]
	tu.AssertEquals("", filepath.Join("xyz", "fre-name.json"), l.Find("xyz", "fr", "name", "json"), t)
	tu.AssertEquals("", filepath.Join("xyz", "eng-name.json"), l.Find("xyz", "en", "name", "json"), t)
	tu.AssertEquals("", filepath.Join("xyz", "abc", "eng-name.xml"), l.Find("xyz/abc", "en", "name", "xml"), t)
	tu.AssertEquals("", "eng-name.xml", l.Find("", "en", "name", "xml"), t)

}

func Test_LangNameLocator_List(t *testing.T) {
	l := FileLocators["LANG-NAME"]
	root, _ := ioutil.TempDir("", "root")

	os.Create(filepath.Join(root, "en-name.json"))
	os.Create(filepath.Join(root, "it-name.xsd"))
	os.Create(filepath.Join(root, "fr-name.json"))

	translations, err := l.List(root, "name", "json")
	if err != nil {
		t.Errorf(err.Error())
	}

	index := map[string]bool{}

	for k, _ := range translations {
		index[k] = true
	}

	if 2 != len(translations) {
		t.Errorf("Wrong size of translations: %d: %v", 2, translations)
	}
	if !index["en"] {
		t.Errorf("Did not find en: %v", translations)
	}
	if !index["fr"] {
		t.Errorf("Did not find fr: %v", translations)
	}

	if filepath.Base(translations["en"]) != "en-name.json" {
		t.Errorf("Did not find en-name.json: %v", translations)
	}

	if filepath.Base(translations["fr"]) != "fr-name.json" {
		t.Errorf("Did not find fr-name.json: %v", translations)
	}

}

func Test_LangNameLocator_List_3charlang(t *testing.T) {
	l := FileLocators["3-CHAR-LANG-NAME"]
	root, _ := ioutil.TempDir("", "root")

	os.Create(filepath.Join(root, "eng-name.json"))
	os.Create(filepath.Join(root, "ita-name.xsd"))
	os.Create(filepath.Join(root, "fre-name.json"))

	translations, err := l.List(root, "name", "json")
	if err != nil {
		t.Errorf(err.Error())
	}

	index := map[string]bool{}

	for k, _ := range translations {
		index[k] = true
	}

	if 2 != len(translations) {
		t.Errorf("Wrong size of translations: %d: %v", 2, translations)
	}
	if !index["en"] {
		t.Errorf("Did not find en: %v", translations)
	}
	if !index["fr"] {
		t.Errorf("Did not find fr: %v", translations)
	}

	if filepath.Base(translations["en"]) != "eng-name.json" {
		t.Errorf("Did not find eng-name.json: %v", translations)
	}

	if filepath.Base(translations["fr"]) != "fre-name.json" {
		t.Errorf("Did not find fre-name.json: %v", translations)
	}

}

func Test_LocDirLocator_Find(t *testing.T) {
	l := FileLocators["LOC-DIR"]
	tu.AssertEquals("", filepath.Join("xyz", "fr", "name.json"), l.Find("xyz", "fr", "name", "json"), t)
	tu.AssertEquals("", filepath.Join("xyz", "en", "name.json"), l.Find("xyz", "en", "name", "json"), t)
	tu.AssertEquals("", filepath.Join("xyz", "abc", "en", "name.xml"), l.Find("xyz/abc", "en", "name", "xml"), t)
	tu.AssertEquals("", filepath.Join("en", "name.xml"), l.Find("", "en", "name", "xml"), t)

	l = FileLocators["3-CHAR-LOC-DIR"]
	tu.AssertEquals("", filepath.Join("xyz", "fre", "name.json"), l.Find("xyz", "fr", "name", "json"), t)
	tu.AssertEquals("", filepath.Join("xyz", "eng", "name.json"), l.Find("xyz", "en", "name", "json"), t)
	tu.AssertEquals("", filepath.Join("xyz", "abc", "eng", "name.xml"), l.Find("xyz/abc", "en", "name", "xml"), t)
	tu.AssertEquals("", filepath.Join("eng", "name.xml"), l.Find("", "en", "name", "xml"), t)
}

func Test_LocDirLocator_List(t *testing.T) {
	l := FileLocators["LOC-DIR"]

	root, _ := ioutil.TempDir("", "root")

	createFile(root, "en", "name.json", t)
	createFile(root, "it", "name.xsd", t)
	createFile(root, "fr", "name.json", t)

	translations, err := l.List(root, "name", "json")
	if err != nil {
		t.Errorf(err.Error())
	}

	index := map[string]bool{}

	for k, _ := range translations {
		index[k] = true
	}

	if 2 != len(translations) {
		t.Errorf("Wrong size of translations: %d: %v", 2, translations)
	}
	if !index["en"] {
		t.Errorf("Did not find en: %v", translations)
	}
	if !index["fr"] {
		t.Errorf("Did not find fr: %v", translations)
	}

	if filepath.Base(translations["en"]) != "name.json" {
		t.Errorf("Did not find name.json: %v", translations)
	}

	if filepath.Base(translations["fr"]) != "name.json" {
		t.Errorf("Did not find name.json: %v", translations)
	}
}

func Test_LocDirLocator_List_3charlang(t *testing.T) {
	l := FileLocators["3-CHAR-LOC-DIR"]
	root, _ := ioutil.TempDir("", "root")

	createFile(root, "eng", "name.json", t)
	createFile(root, "ita", "name.xsd", t)
	createFile(root, "fre", "name.json", t)

	translations, err := l.List(root, "name", "json")
	if err != nil {
		t.Errorf(err.Error())
	}

	index := map[string]bool{}

	for k, _ := range translations {
		index[k] = true
	}

	if 2 != len(translations) {
		t.Errorf("Wrong size of translations: %d: %v", 2, translations)
	}
	if !index["en"] {
		t.Errorf("Did not find eng: %v", translations)
	}
	if !index["fr"] {
		t.Errorf("Did not find fre: %v", translations)
	}

	if filepath.Base(translations["en"]) != "name.json" {
		t.Errorf("Did not find name.json: %v", translations)
	}

	if filepath.Base(translations["fr"]) != "name.json" {
		t.Errorf("Did not find name.json: %v", translations)
	}
}

func createFile(root, loc, name string, t *testing.T) {
	locDir := filepath.Join(root, loc)
	os.Mkdir(locDir, 644)
	if _, err := os.Create(filepath.Join(locDir, name)); err != nil {
		t.Errorf(err.Error())
	}

}

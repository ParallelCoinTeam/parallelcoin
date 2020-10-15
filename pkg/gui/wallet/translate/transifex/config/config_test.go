package config

import (
	tu "github.com/p9c/pod/pkg/gui/wallet/translate/testutil"
	"path/filepath"
	"testing"
)

func Test_ReadConfig_LangDir(t *testing.T) {
	configText := `
[{
	"type": "FLATTENXMLTOJSON",
    "structure": "3-CHAR-LOC-DIR",
    "resources": [{
      "dir": "loc",
      "fname": "string",
      "name": "Iso19139 Strings",
      "slug": "iso19139-strings-xml",
      "priority": "0",
      "categories": ["schemaplugin", "iso19139"]
    },{
      "dir": "loc",
      "fname": "string",	
      "name": "Iso19139 Strings2",
      "slug": "iso19139-strings2-xml",
      "priority": "0",
      "categories": [],
      "ExtraParams": {
      	"key": "xpath/to/key"
      }
    }]
},{    
    "type": "KEYVALUEJSON",
    "structure": "LANG-NAME",
    "resources": [{
      "dir": "js",
      "fname": "core",
      "name": "Angular UI Common Strings",
      "slug": "core",
      "priority": "0",
      "categories": ["Angular_UI"]
    }]
}]`

	root := tu.CreateFileTree(
		tu.Dir("xyz",
			tu.FileAndData("config.json", []byte(configText)),
			tu.Dir("loc",
				tu.Dir("eng", tu.File("string.xml")),
				tu.Dir("fre", tu.File("string.xml")),
				tu.Dir("ger", tu.File("xyz"))),
			tu.Dir("js",
				tu.File("en-core.json"),
				tu.File("it-core.json"),
				tu.File("fr-core.json"))))

	files, err := ReadConfig(filepath.Join(root, "config.json"), root, "en")

	if err != nil {
		t.Errorf("Error reading config. %v", err)
		panic(err)
	}

	t.Errorf("%v", files)
	if len(files) != 3 {
		t.Errorf("Expected 3 file: %v", files)
	}

	if len(files[0].Translations) != 2 {
		t.Errorf("Expected 2 translations: %v", files[0].Translations)
	}

	if v, ok := files[0].Translations["en"]; !ok {
		t.Errorf("Expected en translation: %v", files[0].Translations)
	} else if v != filepath.Join(root, "loc", "eng", "string.xml") {
		t.Errorf("Path of en translation was unexpected: %s", v)
	}

	if v, ok := files[0].Translations["fr"]; !ok {
		t.Errorf("Expected fr translation: %v", files[0].Translations)
	} else if v != filepath.Join(root, "loc", "fre", "string.xml") {
		t.Errorf("Path of fr translation was unexpected: %s", v)
	}

	if len(files[2].Translations) != 3 {
		t.Errorf("Not 3 translations for second translations: %v", files)
	}

	if v, ok := files[2].Translations["fr"]; !ok {
		t.Errorf("Expected fr translation: %v", files[2].Translations)
	} else if v != filepath.Join(root, "js", "fr-core.json") {
		t.Errorf("Path of fr translation was unexpected: %s", v)
	}

	if v, ok := files[2].Translations["en"]; !ok {
		t.Errorf("Expected en translation: %v", files[2].Translations)
	} else if v != filepath.Join(root, "js", "en-core.json") {
		t.Errorf("Path of en translation was unexpected: %s", v)
	}

	if v, ok := files[2].Translations["it"]; !ok {
		t.Errorf("Expected it translation: %v", files[2].Translations)
	} else if v != filepath.Join(root, "js", "it-core.json") {
		t.Errorf("Path of it translation was unexpected: %s", v)
	}

}

func Test_UnmappedLang(t *testing.T) {
	configText := `
[{
	"type": "FLATTENXMLTOJSON",
    "structure": "3-CHAR-LOC-DIR",
    "resources": [{
      "dir": "loc",
      "fname": "string",
      "name": "Iso19139 Strings",
      "slug": "iso19139-strings-xml",
      "priority": "0",
      "categories": ["schemaplugin", "iso19139"]
    }]
}]`
	root := tu.CreateFileTree(
		tu.Dir("xyz",
			tu.FileAndData("config.json", []byte(configText)),
			tu.Dir("loc",
				tu.Dir("eng", tu.File("string.xml")),
				tu.Dir("fre", tu.File("string.xml")),
				tu.Dir("ger", tu.File("xyz")),
				tu.Dir("xxz", tu.File("string.xml")))))

	defer func() {
		if r := recover(); r != nil {
			// panic is expected
		}
	}()

	ReadConfig(filepath.Join(root, "config.json"), root, "en")

	t.Errorf("A panic should have occurred because xxz is not a valid lang code")
}

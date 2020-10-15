package format

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// A format which takes nested xml that has the translations as the text of the leaf nodes
// since the xml can be nested and the nodes can have attributes, the path to the leaf node is encoded
// as the key of the Json.
// the format that is uploaded to transifex will be a valid Json key value formatted file
type FlattenXmlToJson struct {
	key string
}

func (f *FlattenXmlToJson) Init(initParams map[string]interface{}) {
	if keyInterface, has := initParams["key"]; has {
		f.key = keyInterface.(string)
	}
}
func (f *FlattenXmlToJson) Ext() string { return "xml" }

func (f *FlattenXmlToJson) Clean(content []byte) ([]byte, string, error) {

	r := strings.NewReader(string(content))
	parser := xml.NewDecoder(r)

	contentJson := make(map[string]string)
	rawKeys := make(map[string]int)

	key := []string{}
	for {
		token, err := parser.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, "", err
		}
		switch t := token.(type) {
		case xml.StartElement:
			name := nodeName(t)
			if name == "" {
				panic("A node name is an empty string: \n" + string(content))
			}
			key = append(key, name)
		case xml.EndElement:
			key = key[:len(key)-1]
		case xml.CharData:
			if len(key) > 1 {
				fkey, add := finalKey(key, rawKeys, xml.CharData(t))
				if add {
					contentJson[fkey] = string(xml.CharData(t))
				}
			}
		}
	}

	var err error
	content, err = json.Marshal(contentJson)
	if err != nil {
		return nil, "", err
	}
	return content, "KEYVALUEJSON", nil
}

type Node struct {
	name     xml.Name
	attrs    []xml.Attr
	value    xml.CharData
	key      string
	children map[string]Node
}

func (f FlattenXmlToJson) Write(rootDir, langCode, srcLang, filename, translation string, fileLocator FileLocator) error {
	path := fileLocator.Find(rootDir, langCode, filename, "xml")
	fmt.Println("Updating translations file: " + path)

	var translationJson map[string]string
	if err := json.Unmarshal([]byte(translation), &translationJson); err != nil {
		return err
	}
	var out = bytes.Buffer{}
	encoder := xml.NewEncoder(&out)

	srcFile := fileLocator.Find(rootDir, srcLang, filename, "xml")

	srcContent, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer srcContent.Close()

	parser := xml.NewDecoder(srcContent)

	key := []string{}
	rawKeys := map[string]int{}
	for {
		token, err := parser.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		switch t := token.(type) {
		case xml.StartElement:
			key = append(key, nodeName(t))
			encoder.EncodeToken(t)
		case xml.EndElement:
			key = key[:len(key)-1]
			encoder.EncodeToken(t)
		case xml.CharData:
			fkey, add := finalKey(key, rawKeys, xml.CharData(t))
			if add {
				newData := translationJson[fkey]
				delete(translationJson, fkey)

				encoder.EncodeToken(xml.CharData(newData))
			}
		}
	}

	if len(translationJson) > 0 {
		extraKeys := []string{}
		for key, _ := range translationJson {
			extraKeys = append(extraKeys, key)
		}
		return fmt.Errorf("One or more translations did not have a matching key in the xml file: [%s]", strings.Join(extraKeys, ", "))
	}

	if err = encoder.Flush(); err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(path), 644); err != nil {
		return err
	}
	return ioutil.WriteFile(path, out.Bytes(), 644)
}

func nodeName(t xml.StartElement) string {
	elmt := xml.StartElement(t)
	name := elmt.Name.Local
	nodeRep := bytes.Buffer{}
	nodeRep.WriteString(name)
	if len(elmt.Attr) > 0 {
		nodeRep.WriteString("[")
		for i, att := range elmt.Attr {
			if i > 0 {
				nodeRep.WriteString(" and ")
			}
			nodeRep.WriteString(att.Name.Local)
			nodeRep.WriteString("=")
			nodeRep.WriteString(att.Value)
		}
		nodeRep.WriteString("]")
	}

	return strings.TrimSpace(nodeRep.String())
}

// returns the final key and true or "", false.  If the second return value is false
// then it has been determined that the node is not a translation node
func finalKey(key []string, rawKeys map[string]int, t xml.CharData) (string, bool) {
	fkey := strings.TrimSpace(strings.Join(key[1:], " "))
	text := string(xml.CharData(t))
	if fkey != "" && strings.TrimSpace(text) != "" {
		if _, has := rawKeys[fkey]; has {
			i := rawKeys[fkey]
			rawKeys[fkey] = i + 1
			return fmt.Sprintf("%s<%d>", fkey, i), true
		} else {
			rawKeys[fkey] = 2
			return fkey, true
		}
	}

	return "", false
}

package transifex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ListResources(t *testing.T) {
	const data = `[
    {
        "source_language_code": "en",
        "name": "Angular UI Common Strings",
        "i18n_type": "KEYVALUEJSON",
        "priority": "1",
        "slug": "core",
        "categories": null
    }
]`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, data)
	}))
	defer ts.Close()

	var transifexAPI = NewTransifexAPI("project", "", "")
	transifexAPI.ApiUrl = ts.URL + "/"
	resources, err := transifexAPI.ListResources()
	if err != nil {
		t.Error("Failed to list sources", err)
	}

	if len(resources) != 1 {
		t.Error("Expected 1 resource", resources)
	}

	t.Log("ListResources passed", resources)
}

func Test_CreateResource(t *testing.T) {
	var requestJson map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(bytes, &requestJson)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "{}")
	}))
	defer ts.Close()

	var transifexAPI = NewTransifexAPI("project", "", "")
	transifexAPI.ApiUrl = ts.URL + "/"

	resource := UploadResourceRequest{BaseResource{"slug", "name", "KEYVALUEJSON", "hi", ""}, "data", "true"}
	err := transifexAPI.CreateResource(resource)
	if err != nil {
		t.Error("Failed to list sources", err)
	}

	if c := requestJson["content"]; c != "data" {
		t.Error("incorrect content: " + c.(string))
	}
}

package transifex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

const (
	KeyValueJson string = "KEYVALUEJSON"
)

type TransifexAPI struct {
	ApiUrl, Project, username, password string
	client                              *http.Client
	Debug                               bool
}

type BaseResource struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	I18nType string `json:"i18n_type"`
	Priority string `json:"priority"`
	Category string `json:"category"`
}
type Resource struct {
	BaseResource
	SourceLanguage string `json:"source_language_code"`
}
type UploadResourceRequest struct {
	BaseResource
	Content             string `json:"content"`
	Accept_translations string `json:"accept_translations"`
}
type Language struct {
	Coordinators []string `json:"coordinators"`
	LanguageCode string   `json:"language_code"`
	Translators  []string `json:"translators"`
	Reviewers    []string `json:"reviewers"`
}

func NewTransifexAPI(project, username, password string) TransifexAPI {
	return TransifexAPI{"https://www.transifex.com/api/2", project, username, password, &http.Client{}, false}
}

func (t TransifexAPI) ListResources() ([]Resource, error) {
	resp, err := t.execRequest("GET", t.resourcesUrl(true), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, err
	}
	var resources []Resource
	jsonErr := json.Unmarshal(data, &resources)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return resources, nil
}

func (t TransifexAPI) CreateResource(newResource UploadResourceRequest) error {
	data, marshalErr := json.Marshal(newResource)
	if marshalErr != nil {
		return marshalErr
	}
	fmt.Println("\n\n", string(data))
	resp, err := t.execRequest("POST", t.resourcesUrl(false), bytes.NewReader(data))
	if err != nil {
		return err
	}

	checkData, checkErr := t.checkValidJsonResponse(resp, fmt.Sprintf("Failed to create resource: %s\n", newResource.Slug))
	switch checkData.(type) {
	case []interface{}:
		checkDataArray := checkData.([]interface{})
		fmt.Printf(`Create %s Summary:

Strings Added: %v
Strings updated: %v
Strings deleted: %v

`, newResource.Slug, checkDataArray[0], checkDataArray[1], checkDataArray[2])
	}
	return checkErr
}

func (t TransifexAPI) UpdateResourceContent(slug, content string) error {
	data, marshalErr := json.Marshal(map[string]string{"slug": slug, "content": content})
	if marshalErr != nil {
		return marshalErr
	}

	resp, err := t.execRequest("PUT", t.resourceUrl(slug, true)+"content/", bytes.NewReader(data))
	if err != nil {
		return err
	}

	checkData, checkErr := t.checkValidJsonResponse(resp, fmt.Sprintf("Error updating content of %s", slug))
	dataMap := checkData.(map[string]interface{})
	fmt.Printf(`Update %s Source Language Content Summary:

Strings Added: %v
Strings updated: %v
Strings deleted: %v

`, slug, dataMap["strings_added"], dataMap["strings_updated"], dataMap["strings_delete"])
	return checkErr
}

func (t TransifexAPI) ValidateConfiguration() error {
	msg := "Error occurred when checking credentials. Please check credentials and network connection"
	if _, err := t.SourceLanguage(); err != nil {
		return fmt.Errorf(msg)
	}
	return nil
}

func (t TransifexAPI) UploadTranslationFile(slug, langCode, content string) error {
	data, marshalErr := json.Marshal(map[string]string{"content": content})
	if marshalErr != nil {
		return marshalErr
	}

	resp, err := t.execRequest("PUT", t.resourceUrl(slug, true)+"translation/"+langCode+"/", bytes.NewReader(data))
	if err != nil {
		return err
	}

	checkData, checkErr := t.checkValidJsonResponse(resp, fmt.Sprintf("Error adding %s translations for %s", langCode, slug))
	dataMap := checkData.(map[string]interface{})
	fmt.Printf(`Update %s %s Translation summary:

Strings Added: %v
Strings updated: %v
Strings deleted: %v

`, slug, langCode, dataMap["strings_added"], dataMap["strings_updated"], dataMap["strings_delete"])
	return checkErr
}

func (t TransifexAPI) SourceLanguage() (string, error) {
	if t.Debug {
		fmt.Println("Executing transifex.SourceLanguage")
	}

	jsonData, err := t.getJson(t.ApiUrl+"/project/"+t.Project, "Error loading SourceLanguage")
	if err != nil {
		return "", err
	}
	lang, has := jsonData.(map[string]interface{})["source_language_code"]
	if !has {
		return "", fmt.Errorf("An error occurred while reading response. Expected a 'source_language_code' json field:\n%s", jsonData)
	}
	sourceLang := lang.(string)
	if strings.TrimSpace(sourceLang) == "" {
		return "", fmt.Errorf("No source language found.  This is probably a bug")
	}

	if t.Debug {
		fmt.Printf("Source language read: %s\n", sourceLang)
	}

	return sourceLang, nil
}
func (t TransifexAPI) Languages() ([]Language, error) {
	resp, err := t.execRequest("GET", fmt.Sprintf("%s/project/%s/languages", t.ApiUrl, t.Project), nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	var jsonData []Language
	if err = json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}

	return jsonData, nil
}
func (t TransifexAPI) DownloadTranslations(slug string) (map[string]string, error) {
	sourceLang, err := t.SourceLanguage()
	if err != nil {
		return nil, err
	}
	fullLangs, langErr := t.Languages()
	if langErr != nil {
		return nil, langErr
	}
	langs := make([]string, len(fullLangs)+1)
	langs[0] = sourceLang
	for i, l := range fullLangs {
		langs[i+1] = l.LanguageCode
	}

	translations := make(map[string]string, len(langs))
	for _, lang := range langs {
		url := fmt.Sprintf("%s/project/%s/resource/%s/translation/%s", t.ApiUrl, t.Project, slug, lang)
		data, err2 := t.getJson(url, "Error downloing translations file")
		if err2 != nil {
			return nil, err2
		}

		translations[lang] = data.(map[string]interface{})["content"].(string)
	}
	return translations, nil
}

func (t TransifexAPI) getJson(url string, errMsg string) (interface{}, error) {
	resp, err := t.execRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	jsonData, checkErr := t.checkValidJsonResponse(resp, errMsg)
	if checkErr != nil {
		return nil, checkErr
	}

	return jsonData, nil
}
func (t TransifexAPI) execRequest(method string, url string, requestData io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, requestData)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(t.username, t.password)
	if requestData != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	if t.Debug {
		fmt.Printf("\nExecuting http %s request: '%s'\n\n", method, url)
		dump, _ := httputil.DumpRequest(request, true)
		fmt.Println(string(dump))
	}

	resp, finalErr := t.client.Do(request)

	if t.Debug {
		dump, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(dump))

	}

	if resp.StatusCode > 400 {
		return nil, fmt.Errorf("Response Code: %v\nResponse Status: %s", resp.StatusCode, resp.Status)
	}

	return resp, finalErr
}

func (t TransifexAPI) resourcesUrl(endSlash bool) string {
	url := fmt.Sprintf("%s/project/%s/resources", t.ApiUrl, t.Project)
	if endSlash {
		return url + "/"
	}
	return url
}
func (t TransifexAPI) resourceUrl(slug string, endSlash bool) string {
	url := fmt.Sprintf("%s/project/%s/resource/%s", t.ApiUrl, t.Project, slug)
	if endSlash {
		return url + "/"
	}
	return url
}

func (t TransifexAPI) checkValidJsonResponse(resp *http.Response, errorMsg string) (interface{}, error) {

	responseData, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		return nil, readErr
	}

	var jsonData interface{}
	if err := json.Unmarshal(responseData, &jsonData); err != nil {
		return nil, fmt.Errorf(errorMsg + "\n\nError:\n" + string(responseData))
	}

	if t.Debug {

		fmt.Printf("Response: %s", jsonData)
		switch jsonData.(type) {
		case map[string]interface{}:
			for key, val := range jsonData.(map[string]interface{}) {
				fmt.Println(key, ":", val)
			}
		case []interface{}:
			for _, val := range jsonData.([]interface{}) {
				fmt.Println(val)
			}
		default:
			fmt.Printf("Response: %s", jsonData)
		}
	}

	fmt.Print("\n")
	return jsonData, nil
}

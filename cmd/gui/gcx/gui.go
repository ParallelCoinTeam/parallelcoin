package gcx

import (
	"github.com/p9c/pod/pkg/gui/webview"
	"net/http"
)

type GUI struct {
	Boot *Boot            `json:"boot"`
	Cf   *Configuration   `json:"cf"`
	Fs   *http.FileSystem `json:"fs"`
	Wv   webview.WebView  `json:"wv"`
}

type Boot struct {
	IsBoot     bool `json:"boot"`
	IsFirstRun bool `json:"firstrun"`
	IsBootMenu bool `json:"menu"`
	IsBootLogo bool `json:"logo"`
	IsLoading  bool `json:"loading"`
}

type Configuration struct {
	Assets   string `json:"assets"`
	Theme    bool   `json:"theme"`
	IsDev    bool   `json:"dev"`
}

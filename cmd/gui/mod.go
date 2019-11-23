package gui

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/webview"
	"net/http"
)

type Bios struct {
	cx         *conte.Xt
	Fs         *http.FileSystem `json:"fs"`
	Wv         webview.WebView  `json:"wv"`
	Theme      bool             `json:"theme"`
	IsBoot     bool             `json:"boot"`
	IsFirstRun *bool             `json:"firstrun"`
	IsBootMenu bool             `json:"menu"`
	IsBootLogo bool             `json:"logo"`
	IsLoading  bool             `json:"loading"`
	IsDev      bool             `json:"dev"`
	IsScreen   string           `json:"screen"`
}

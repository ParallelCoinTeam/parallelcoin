package model

import (
	"gioui.org/app"
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/controller"
	"github.com/p9c/pod/cmd/gui/theme"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/pod"
)

type DuOScomponent struct {
	Name       string
	Version    string
	Model      interface{}
	View       func()
	Controller func()
}

type DuoUI struct {
	Window     *app.Window
	Context    *layout.Context
	Theme      *theme.DuoUItheme
	Pages      *DuoUIpages
	Navigation map[string]*theme.DuoUIthemeNav
	//Configuration *DuoUIconfiguration
	IsReady bool
}

type DuoUIpages struct {
	CurrentPage *theme.DuoUIpage
	Controller  map[string]*controller.DuoUIpage
	Theme       map[string]*theme.DuoUIpage
}
type DuoUIlog struct {
	LogMessages []log.Entry
	LogChan     chan log.Entry
	StopLogger  chan struct{}
}

//type DuoUIconfiguration struct {
//	Abbrevation        string
//	PrimaryTextColor   color.RGBA
//	SecondaryTextColor color.RGBA
//	PrimaryBgColor     color.RGBA
//	SecondaryBgColor   color.RGBA
//	Navigations        map[string]*view.DuoUIthemeNav
//}

type DuoUIconfTabs struct {
	Current  string
	TabsList map[string]*controller.Button
}

//type DuoUIalert struct {
//	Time      time.Time   `json:"time"`
//	Title     string      `json:"title"`
//	Message   interface{} `json:"message"`
//	AlertType string      `json:"type"`
//}

type DuoUIsettings struct {
	Abbrevation string
	Tabs        *DuoUIconfTabs
	Daemon      *DaemonConfig `json:"daemon"`
}

type DaemonConfig struct {
	Config  *pod.Config `json:"config"`
	Schema  pod.Schema  `json:"schema"`
	Widgets map[string]interface{}
}

type DuoUIblock struct {
	Height        int64   `json:"height"`
	BlockHash     string  `json:"hash"`
	PowAlgoID     uint32  `json:"pow"`
	Difficulty    float64 `json:"diff"`
	Amount        float64 `json:"amount"`
	TxNum         int     `json:"txnum"`
	Confirmations int64
	Time          string `json:"time"`
	Link          *controller.Button
}

type DuoUItoast struct {
	Title   string
	Message string
}

package model

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/app"
	"github.com/p9c/pod/pkg/gui/layout"
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
	Window  *app.Window
	Context *layout.Context
	Theme   *theme.DuoUItheme
	Pages   map[string]*theme.DuoUIpage
	Navigation map[string]*theme.DuoUIthemeNav
	//Configuration *DuoUIconfiguration
	Quit    chan struct{}
	Ready   chan struct{}
	IsReady bool
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
	Height     int64   `json:"height"`
	BlockHash  string  `json:"hash"`
	PowAlgoID  uint32  `json:"pow"`
	Difficulty float64 `json:"diff"`
	Amount     float64 `json:"amount"`
	TxNum      int     `json:"txnum"`
	Time       string  `json:"time"`
}

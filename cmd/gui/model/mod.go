package model

import (
	"sync"

	"gioui.org/app"
	"gioui.org/layout"
	"go.uber.org/atomic"

	"github.com/p9c/gel"
	"github.com/p9c/gelook"
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
	Theme      *gelook.DuoUItheme
	Pages      *DuoUIpages
	Navigation map[string]*gelook.DuoUIthemeNav
	// Configuration *DuoUIconfiguration
	IsReady bool
}

type DuoUIpages struct {
	CurrentPage *gelook.DuoUIpage
	Controller  map[string]*gel.DuoUIpage
	Theme       map[string]*gelook.DuoUIpage
}
type DuoUIlog struct {
	Mx          sync.Mutex
	LogMessages atomic.Value // []log.Entry
	LogChan     chan log.Entry
	StopLogger  chan struct{}
}

// type DuoUIconfiguration struct {
//	Abbrevation        string
//	PrimaryTextColor   color.RGBA
//	SecondaryTextColor color.RGBA
//	PrimaryBgColor     color.RGBA
//	SecondaryBgColor   color.RGBA
//	Navigations        map[string]*view.DuoUIthemeNav
// }

type DuoUIconfTabs struct {
	Current  string
	TabsList map[string]*gel.Button
}

// type DuoUIalert struct {
//	Time      time.Time   `json:"time"`
//	Title     string      `json:"title"`
//	Message   interface{} `json:"message"`
//	AlertType string      `json:"type"`
// }

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
	Time          int64 `json:"time"`
	Link          *gel.Button
}

type DuoUItoast struct {
	Title   string
	Message string
}

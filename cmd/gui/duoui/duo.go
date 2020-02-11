package duoui

import (
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
	"github.com/p9c/pod/pkg/log"
	"image/color"

	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/fonts"
	"github.com/p9c/pod/pkg/gui/app"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
)

type DuoUI struct {
	m *models.DuoUI
}

func startLogger() (stopLogger chan struct{}) {
	log.L.LogChan = logChan
	log.L.SetLevel("Info", false)
	go func() {
	out:
		for {
			select {
			case n := <-log.L.LogChan:
				logMessages = append(logMessages, n)
			case <-stopLogger:
				defer func() {
					stopLogger = make(chan struct{})
				}()
				logMessages = []log.Entry{}
				log.L.LogChan = nil
				break out
			}
		}
	}()
	return
}

func DuOuI(rc *rcd.RcVar, cx *conte.Xt) (duo *models.DuoUI, err error) {

	duo = &models.DuoUI{
		CurrentPage: "overview",
		DuoUIwindow: app.NewWindow(
			app.Size(unit.Dp(1024), unit.Dp(640)),
			app.Title("ParallelCoin"),
		),
		Quit:  make(chan struct{}),
		Ready: make(chan struct{}, 1),
	}
	fonts.Register()
	duo.DuoUIcontext = layout.NewContext(duo.DuoUIwindow.Queue())

	startLogger()

	navigations := make(map[string]*parallel.DuoUIthemeNav)
	//navigations["mainMenu"] = mainMenu()
	duo.DuoUIconfiguration = &models.DuoUIconfiguration{
		Abbrevation:        "DUO",
		PrimaryBgColor:     color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30},
		SecondaryBgColor:   color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf},
		PrimaryTextColor:   color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf},
		SecondaryTextColor: color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30},
		Navigations:        navigations,
	}
	// Icons
	rc.Settings.Daemon = rcd.GetCoreSettings(cx)
	// Settings tabs
	confTabs := make(map[string]*widget.Button)
	settingsFields := make(map[string]interface{})
	for _, group := range rc.Settings.Daemon.Schema.Groups {
		confTabs[group.Legend] = new(widget.Button)
		for _, field := range group.Fields {
			switch field.Type {
			case "array":
				settingsFields[field.Name] = new(widget.Button)
			case "input":
				settingsFields[field.Name] = &widget.Editor{
					SingleLine: true,
					Submit:     true,
				}
			case "switch":
				settingsFields[field.Name] = new(widget.CheckBox)
			case "radio":
				settingsFields[field.Name] = new(widget.Enum)
			default:
				settingsFields[field.Name] = new(widget.Button)
			}
		}
	}
	duo.DuoUIconfiguration.Tabs = models.DuoUIconfTabs{
		Current:  "wallet",
		TabsList: confTabs,
	}
	duo.DuoUIconfiguration.Settings.Daemon.Widgets = settingsFields
	duo.DuoUItheme = parallel.NewDuoUItheme()
	return
}

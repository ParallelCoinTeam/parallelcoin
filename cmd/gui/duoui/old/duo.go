package duoui

import (
	"sync"

	"github.com/stalker-loki/pod/cmd/gui/component"
	"github.com/stalker-loki/pod/cmd/gui/pages"
	"github.com/stalker-loki/pod/cmd/gui/rcd"
	"github.com/stalker-loki/pod/pkg/gui/clipboard"
	"github.com/stalker-loki/pod/pkg/gui/gelook"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/stalker-loki/pod/cmd/gui/model"
)

var clipboardStarted bool
var clipboardMu sync.Mutex

type DuoUI struct {
	ly *model.DuoUI
	rc *rcd.RcVar
}

func DuOuI(rc *rcd.RcVar) (duo *model.DuoUI, err error) {
	duo = &model.DuoUI{
		Window: app.NewWindow(
			app.Size(unit.Dp(1024), unit.Dp(640)),
			app.Title("ParallelCoin"),
		),
	}
	duo.Context = layout.NewContext(duo.Window.Queue())
	// rc.StartLogger()
	// sys.Components["logger"].View()
	// d.sys.Components["logger"].View
	duo.Navigation = &model.DuoUInav{
		Items: make(map[string]*gelook.DuoUIthemeNav),
	}
	// navigations["mainMenu"] = mainMenu()
	// Icons
	// rc.Settings.Daemon = rcd.GetCoreSettings()
	duo.Theme = gelook.NewDuoUItheme()
	// duo.Pages = components.LoadPages(duo.Context, duo.Theme, rc)
	duo.Pages = &model.DuoUIpages{
		Controller: nil,
		Theme:      pages.LoadPages(rc, duo.Context, duo.Theme),
	}
	component.SetPage(rc, duo.Pages.Theme["OVERVIEW"])
	clipboardMu.Lock()
	if !clipboardStarted {
		clipboardStarted = true
		clipboard.Start()
	}
	clipboardMu.Unlock()
	return
}

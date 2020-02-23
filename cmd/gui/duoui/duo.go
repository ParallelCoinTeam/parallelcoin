package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/clipboard"
	"sync"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/pkg/fonts"
)

var clipboardStarted bool
var clipboardMu sync.Mutex

type DuoUI struct {
	ly *model.DuoUI
	rc *rcd.RcVar
}

func DuOuI() (duo *model.DuoUI, err error) {

	duo = &model.DuoUI{
		Window: app.NewWindow(
			app.Size(unit.Dp(1024), unit.Dp(640)),
			app.Title("ParallelCoin"),
		),
	}
	fonts.Register()
	duo.Context = layout.NewContext(duo.Window.Queue())

	//rc.StartLogger()
	//sys.Components["logger"].View()

	//d.sys.Components["logger"].View

	duo.Navigation = make(map[string]*theme.DuoUIthemeNav)
	//navigations["mainMenu"] = mainMenu()

	// Icons
	//rc.Settings.Daemon = rcd.GetCoreSettings()

	duo.Theme = theme.NewDuoUItheme()
	//duo.Pages = components.LoadPages(duo.Context, duo.Theme, rc)

	clipboardMu.Lock()
	if !clipboardStarted {
		clipboardStarted = true
		clipboard.Start()
	}
	clipboardMu.Unlock()

	return
}

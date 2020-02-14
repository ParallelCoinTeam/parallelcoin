package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"

	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/pkg/fonts"
	"github.com/p9c/pod/pkg/gui/app"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
)

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
		Quit:  make(chan struct{}),
		Ready: make(chan struct{}, 1),
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
	return
}

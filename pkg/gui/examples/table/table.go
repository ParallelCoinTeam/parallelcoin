package main

import (
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/util/logi"
	qu "github.com/p9c/pod/pkg/util/qu"
	
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
)

type App struct {
	th   *gui.Theme
	quit qu.C
}

func main() {
	logi.L.SetLevel("trace", false, "pod")
	quit := qu.T()
	th := gui.NewTheme(p9fonts.Collection(), quit)
	model := App{
		th: th,
	}
	go func() {
		if e := gui.NewWindow(th).
			Size(64, 32).
			Title("table example").
			Open().
			Run(
				model.mainWidget, func(l.Context) {}, func() {
					quit.Q()
				}, quit,
		); err.Chk(e) {
		}
	}()
	<-quit
}

func (m *App) mainWidget(gtx l.Context) l.Dimensions {
	th := m.th
	gtx.Constraints.Max = gtx.Constraints.Min
	dims := th.Flex().AlignStart().Rigid(
		m.th.Body1("test").Fn,
	).Fn(gtx)
	Infos(dims)
	return dims
}

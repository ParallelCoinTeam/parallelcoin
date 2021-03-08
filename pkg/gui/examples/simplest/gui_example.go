package main

import (
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui"
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
)

type App struct {
	th   *gui.Theme
	quit qu.C
}

func main() {
	quit := qu.T()
	th := gui.NewTheme(p9fonts.Collection(), quit)
	minerModel := App{
		th: th,
	}
	go func() {
		if e := gui.NewWindow(th).
			Size(64, 32).
			Title("nothing to see here").
			Open().
			Run(
				minerModel.mainWidget, func(l.Context) {}, func() {
					quit.Q()
				}, quit,
			); dbg.Chk(e) {
		}
	}()
	<-quit
}

func (m *App) mainWidget(gtx l.Context) l.Dimensions {
	th := m.th
	return th.Flex().Rigid(
		gui.EmptyMaxWidth(),
	).Fn(gtx)
}

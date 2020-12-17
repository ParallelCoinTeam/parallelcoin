package main

import (
	l "gioui.org/layout"
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
)

type App struct {
	th   *p9.Theme
	quit qu.C
}

func main() {
	quit := qu.T()
	th := p9.NewTheme(p9fonts.Collection(), quit)
	minerModel := App{
		th: th,
	}
	go func() {
		if err := f.NewWindow(th).
			Size(64, 32).
			Title("nothing to see here").
			Open().
			Run(
				minerModel.mainWidget, func(l.Context) {}, func() {
					quit.Q()
				}, quit,
		); Check(err) {
		}
	}()
	<-quit
}

func (m *App) mainWidget(gtx l.Context) l.Dimensions {
	th := m.th
	return th.Flex().Rigid(
		p9.EmptyMaxWidth(),
	).Fn(gtx)
}

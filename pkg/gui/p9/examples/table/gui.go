package main

import (
	l "gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
)

type App struct {
	th   *p9.Theme
	quit chan struct{}
}

func main() {
	quit := make(chan struct{})
	th := p9.NewTheme(p9fonts.Collection(), quit)
	minerModel := App{
		th: th,
	}
	go func() {
		if err := f.NewWindow().
			Size(800, 600).
			Title("table example").
			Open().
			Run(minerModel.mainWidget, func(l.Context) {}, func() {
				close(quit)
			}, quit); Check(err) {
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

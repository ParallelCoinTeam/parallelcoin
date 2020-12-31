package main

import (
	qu "github.com/p9c/pod/pkg/util/quit"
	
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
)

type Model struct {
	th *p9.Theme
}

func main() {
	quit := qu.T()
	th := p9.NewTheme(p9fonts.Collection(), quit)
	minerModel := Model{
		th: th,
	}
	go func() {
		if err := f.NewWindow(th).
			Size(64, 32).
			Title("example").
			Open().
			Run(
				minerModel.mainWidget,
				func(l.Context) {},
				func() {
					quit.Q()
				}, quit,
			); Check(err) {
		}
	}()
	<-quit
}

func (m *Model) mainWidget(gtx l.Context) l.Dimensions {
	return m.th.Flex().Flexed(1,
		m.th.Fill("red", m.th.Flex().AlignMiddle().SpaceAround().
			Flexed(1,
				m.th.H6("example").Fn,
			).Fn, l.Center).Fn,
	).Fn(gtx)
}

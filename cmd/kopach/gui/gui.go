package gui

import (
	"gioui.org/app"
	"gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/util/interrupt"
)

type MinerModel struct {
	*p9.Theme
}

func Run(quit chan struct{}) {
	th := p9.NewTheme(p9fonts.Collection(), quit)
	minerModel := MinerModel{
		Theme: th,
	}
	go func() {
		if err := f.Window().
			Size(640, 480).
			Title("parallelcoin kopach miner control gui").
			Open().
			Run(
				minerModel.Widget,
				func() {
					Debug("quitting miner")
					close(quit)
					interrupt.Request()
				}); Check(err) {
		}
	}()
	app.Main()
}

func (m *MinerModel) Widget(gtx layout.Context) layout.Dimensions {
	return m.Fill("DocBg").Embed(
		m.Flex().Vertical().Rigid(
			m.Label().Text("this is a test").Fn,
		).Fn,
	).Fn(gtx)
}

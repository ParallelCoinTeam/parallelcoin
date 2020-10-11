package gui

import (
	"os"

	"gioui.org/app"
	"gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
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
			Run(func(ctx *layout.Context) {
				minerModel.Widget(*ctx)
			}, func() {
				close(quit)
				os.Exit(0)
			}); Check(err) {
		}
	}()
	app.Main()
}

func (m *MinerModel) Widget(gtx layout.Context) {
	m.Fill("DocBg").Embed(
		m.Flex().Vertical().Rigid(
			m.Label().Text("this is a test").Fn,
		).Fn,
	).Fn(gtx)
}

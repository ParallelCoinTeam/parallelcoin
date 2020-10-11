package gui

import (
	"gioui.org/app"
	"gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	icons "github.com/p9c/pod/pkg/gui/ico/svg"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/util/interrupt"
)

type MinerModel struct {
	*p9.Theme
	DarkTheme  bool
	logoButton *p9.Clickable
}

func Run(quit chan struct{}) {
	th := p9.NewTheme(p9fonts.Collection(), quit)
	minerModel := MinerModel{
		Theme:     th,
		DarkTheme: false,
		logoButton: th.Clickable().SetClick(func() {
			Debug("clicked logo button")
		}),
	}
	minerModel.SetTheme(false)
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

func (m *MinerModel) Widget(gtx layout.Context) {
	m.Flex().Vertical().Rigid(
		m.Fill("PanelBg").Embed(
			// m.Inset(0.25).Embed(
			m.Flex().Rigid(
				m.Inset(0.25).Embed(
				m.IconButton(m.logoButton.SetClick(
					func() {
						Info("clicked logo button")
						m.FlipTheme()
					})).
					Color("PanelBg").
					Background("PanelText").
					Scale(p9.Scales["H4"]).
					Icon(icons.ParallelCoin).
					Fn,
				).Fn,
			).Rigid(
				m.Inset(0.5).Embed(
					// m.Fill("Primary").Embed(
					m.H5("kopach miner control").
						Color("PanelText").
						Fn,
					// ).Fn,
				).Fn,
			).Fn,
			// ).Fn,
		).Fn,
	).Flexed(1,
		m.Fill("DocBg").Embed(
			m.Body1("Body1").Color("DocText").Fn,
		).Fn,
	).Fn(gtx)
}

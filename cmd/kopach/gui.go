package kopach

import (
	"fmt"
	"image"
	"runtime"
	"time"

	"gioui.org/app"
	"gioui.org/layout"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/kopach/gui"
	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	icons "github.com/p9c/pod/pkg/gui/ico/svg"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/util/interrupt"
)

var maxThreads = float32(runtime.NumCPU())

type MinerModel struct {
	*p9.Theme
	Cx         *conte.Xt
	worker     *Worker
	DarkTheme  bool
	logoButton *p9.Clickable
	mineToggle *p9.Bool
	cores      *p9.Float
	nCores     int
	solData    []SolutionData
	solButtons []*p9.Clickable
}

func Run(w *Worker, cx *conte.Xt) {
	th := p9.NewTheme(p9fonts.Collection(), w.quit)
	Debug(*cx.Config.Generate, *cx.Config.GenThreads)
	minerModel := MinerModel{
		worker:    w,
		Theme:     th,
		DarkTheme: false,
		logoButton: th.Clickable().SetClick(func() {
			gui.Debug("clicked logo button")
		}),
		mineToggle: th.Bool(*cx.Config.Generate),
		cores:      th.Float().SetValue(float32(*cx.Config.GenThreads)),
		solButtons: make([]*p9.Clickable, 201),
	}
	for i := 0; i < 201; i++ {
		minerModel.solButtons[i] = th.Clickable()
	}
	minerModel.SetTheme(minerModel.DarkTheme)
	go func() {
		if err := f.Window().
			Size(640, 480).
			Title("parallelcoin kopach miner control gui").
			Open().
			Run(
				minerModel.Widget,
				func() {
					gui.Debug("quitting miner")
					close(w.quit)
					interrupt.Request()
				}); gui.Check(err) {
		}
	}()
	app.Main()
}

func (m *MinerModel) Widget(gtx layout.Context) {
	m.Flex().Vertical().Rigid(
		m.Fill("PanelBg").Embed(
			m.Flex().Rigid(
				m.Inset(0.25).Embed(
					m.IconButton(m.logoButton.SetClick(
						func() {
							gui.Info("clicked logo button")
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
					m.H5("kopach miner control").
						Color("PanelText").
						Fn,
				).Fn,
			).Fn,
		).Fn,
	).Flexed(1,
		m.Fill("DocBg").Embed(
			m.Inset(0.5).Embed(
				m.Flex().Vertical().Flexed(1,
					// m.Inset(0.5).Embed(
					m.Flex().Vertical().Rigid(
						m.H5("miner settings").Fn,
					).Rigid(
						m.Flex().Flexed(0.5,
							m.Body1("enable mining").
								Color("DocText").
								Fn,
						).Flexed(0.5,
							m.Switch(m.mineToggle.SetOnChange(
								func(b bool) {
									if b {
										gui.Debug("start mining")
										m.worker.StartChan <- struct{}{}
									} else {
										gui.Debug("stop mining")
										m.worker.StopChan <- struct{}{}
									}
								})).
								Fn,
						).Fn,
						// ).Fn,
					).Rigid(
						m.Flex().Rigid(
							// m.Inset(0.5).Embed(
							m.Flex().Flexed(0.5,
								m.Body1("number of mining threads"+
									fmt.Sprintf("%3v", int(m.cores.Value()+0.5))).
									Fn,
							).Rigid(
								m.Caption("0").
									Color("Primary").
									Fn,
							).Flexed(0.5,
								m.Slider().
									Float(m.cores.SetHook(func(fl float32) {
										iFl := int(fl + 0.5)
										if m.nCores != iFl {
											Debug("cores value changed", iFl)
										}
										m.nCores = iFl
										m.cores.SetValue(float32(iFl))
										// TODO: then restart the threads
										m.worker.SetThreads <- m.nCores
									})).
									Min(0).Max(maxThreads).
									Fn,
							).Rigid(
								m.Caption(fmt.Sprint(int(maxThreads))).
									Color("Primary").
									Fn,
							).Fn,
							// ).Fn,
						).Fn,
					).Rigid(
						m.Flex().Vertical().Rigid(
							func(ctx layout.Context) layout.Dimensions {
								return layout.Dimensions{
									Size: image.Point{
										X: int(m.TextSize.Scale(2).V),
										Y: int(m.TextSize.Scale(2).V),
									},
									Baseline: 0,
								}
							},
						).Rigid(
							m.H5("found blocks").Fn,
						).Flexed(1,
							m.List().Vertical().
								Length(len(m.worker.solutions)).
								ListElement(func(gtx layout.Context, index int) layout.Dimensions {
									if m.worker.solutionsUpdated.Load() {
										// regenerate m.solData
										m.solData = make([]SolutionData, len(m.worker.solutions))
										for i := range m.worker.solutions {
											m.solData[i] = m.worker.solutions[i]
										}
										m.worker.solutionsUpdated.Store(false)
									}
									// display from m.solData
									return m.Flex().Rigid(
										m.Button(
											m.solButtons[index].SetClick(func() {
												Debug("clicked for block", m.solData[index].height)
											})).Text(fmt.Sprint(m.solData[index].height)).Fn,
									).Flexed(1,
										m.Inset(0.5).Embed(
											m.Flex().Vertical().Rigid(
												m.Body1(fmt.Sprint(m.solData[index].block.BlockHash())).
													Font("go regular").
													TextScale(0.75).Fn,
											).Rigid(
												m.Body1(fmt.Sprint(m.solData[index].time.Format(time.RFC3339))).Fn,
											).Fn,
										).Fn,
									).Fn(gtx)
								}).
								Fn,
						).Fn,
					).Fn,
				).Fn,
			).Fn,
		).Fn,
	).Fn(gtx)
}

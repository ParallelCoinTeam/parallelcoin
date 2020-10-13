package kopach

import (
	"fmt"
	"image"
	"runtime"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/text"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	icons "github.com/p9c/pod/pkg/gui/ico/svg"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/util/interrupt"
)

var maxThreads = float32(runtime.NumCPU())

type MinerModel struct {
	*p9.Theme
	Cx            *conte.Xt
	worker        *Worker
	DarkTheme     bool
	logoButton    *p9.Clickable
	mineToggle    *p9.Bool
	cores         *p9.Float
	nCores        int
	solButtons    []*p9.Clickable
	lists         map[string]*layout.List
	solutionCount int
}

func Run(w *Worker, cx *conte.Xt) {
	th := p9.NewTheme(p9fonts.Collection(), w.quit)
	Debug(*cx.Config.Generate, *cx.Config.GenThreads)
	solButtons := make([]*p9.Clickable, 201)
	for i := range solButtons {
		solButtons[i] = th.Clickable()
	}
	lists := map[string]*layout.List{
		"found": {
			Axis:      layout.Vertical,
			Alignment: layout.Start,
		},
	}
	minerModel := MinerModel{
		worker:    w,
		Theme:     th,
		DarkTheme: false,
		logoButton: th.Clickable().SetClick(func() {
			Debug("clicked logo button")
		}),
		mineToggle: th.Bool(*cx.Config.Generate),
		cores:      th.Float().SetValue(float32(*cx.Config.GenThreads)),
		solButtons: solButtons,
		lists:      lists,
	}
	for i := 0; i < 201; i++ {
		minerModel.solButtons[i] = th.Clickable()
	}
	minerModel.SetTheme(minerModel.DarkTheme)
	win := f.Window()
	go func() {
		if err := win.
			Size(640, 480).
			Title("kopach").
			Open().
			Run(
				minerModel.Widget,
				func() {
					Debug("quitting miner")
					close(w.quit)
					interrupt.Request()
				}); Check(err) {
		}
	}()
	go func() {
		for {
			select {
			case <-minerModel.worker.Update:
				win.Window.Invalidate()
			}
		}
	}()
	app.Main()
}

func (m *MinerModel) Widget(gtx layout.Context) {
	counter := 0
	m.Flex().Vertical().Rigid(
		m.Fill("PanelBg").Embed(
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
					m.H5("kopach miner control").
						Color("PanelText").
						Fn,
				).Fn,
			).Fn,
		).Fn,
	).Flexed(1,
		m.Fill("DocBg").Embed(
			m.Inset(0.5).Embed(
				m.Flex().Vertical().Rigid(
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
										Debug("start mining")
										m.worker.StartChan <- struct{}{}
									} else {
										Debug("stop mining")
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
							m.Fill("PanelBg").Embed(
								m.Inset(0.25).Embed(
									m.Flex().Flexed(1,
										func(gtx layout.Context) layout.Dimensions {
											cs := gtx.Constraints
											cs.Min = cs.Max
											return m.lists["found"].
												Layout(gtx, m.worker.solutionCount,
													func(c layout.Context, i int) layout.Dimensions {
														counter++
														return m.Flex().Rigid(
															m.Inset(0.25).Embed(
																m.Button(
																	m.solButtons[i].SetClick(func() {
																		Debug("clicked for block", m.worker.solutions[i].height)
																	})).Text(fmt.Sprint(m.worker.solutions[i].height)).Fn,
															).Fn,
														).Flexed(1,
															m.Inset(0.25).Embed(
																m.Flex().Vertical().Rigid(
																	m.Flex().Rigid(
																		// m.Inset(0.25).Embed(
																		m.Body1(m.worker.solutions[i].algo).Font("bariol bold").Fn,
																		// ).Fn,
																	).Flexed(1,
																		// m.Inset(0.25).Embed(
																		m.Body1(fmt.Sprint(
																			m.worker.solutions[i].time.Format(time.RFC3339))).
																			Alignment(text.End).Fn,
																	).Fn,
																	// ).Fn,
																).Rigid(
																	m.Body1(m.worker.solutions[i].hash).
																		Font("go regular").
																		TextScale(0.75).
																		Alignment(text.End).Fn,
																).Fn,
															).Fn,
														).Fn(c)
													})
										},
									).Fn,
								).Fn,
							).Fn,
						).Fn,
					).Fn,
				).Fn,
			).Fn,
		).Fn,
	).Fn(gtx)
}

package kopach

import (
	"fmt"
	"image"
	"runtime"
	"time"
	
	l "gioui.org/layout"
	"gioui.org/text"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/gui"
	icons "github.com/p9c/pod/pkg/gui/ico/svg"
)

type MinerModel struct {
	*gui.Window
	Cx                     *conte.Xt
	worker                 *Worker
	DarkTheme              bool
	logoButton             *gui.Clickable
	mineToggle             *gui.Bool
	nCores                 int
	solButtons             []*gui.Clickable
	lists                  map[string]*gui.List
	solutionCount          int
	modalWidget            l.Widget
	modalOn                bool
	modalScrim, modalClose *gui.Clickable
	password               *gui.Password
	threadSlider           *gui.IntSlider
}


func (m *MinerModel) Widget(gtx l.Context) l.Dimensions {
	return m.Stack().Stacked(
		m.Flex().Flexed(
			1,
			m.VFlex().
				Rigid(m.Header).Flexed(
				1,
				m.Fill("DocBg", 0, 0, 0, m.Inset(
					0.5,
					m.VFlex().
						Rigid(m.H5("miner settings").Fn).
						Rigid(m.RunControl).
						Rigid(m.SetThreads).
						Rigid(m.PreSharedKey).
						Rigid(m.VSpacer).
						Rigid(m.H5("found blocks").Fn).
						Rigid(
							m.Fill("PanelBg", l.Center, 0, 0, m.FoundBlocks).Fn,
						).Fn,
				).Fn).Fn,
			).Fn,
		).Fn,
	).
		Stacked(
			func(gtx l.Context) l.Dimensions {
				if m.modalOn {
					return m.Fill("scrim", l.Center, 0, 0, m.VFlex().
						Flexed(
							0.1,
							m.Flex().Rigid(
								func(gtx l.Context) l.Dimensions {
									return l.Dimensions{
										Size: image.Point{
											X: gtx.Constraints.Max.X,
											Y: gtx.Constraints.Max.Y,
										},
										Baseline: 0,
									}
								},
							).Fn,
						).AlignMiddle().
						Rigid(m.modalWidget).
						Flexed(
							0.1,
							m.Flex().Rigid(
								func(gtx l.Context) l.Dimensions {
									return l.Dimensions{
										Size: image.Point{
											X: gtx.Constraints.Max.X,
											Y: gtx.Constraints.Max.Y,
										},
										Baseline: 0,
									}
								},
							).Fn,
						).Fn).Fn(gtx)
				} else {
					return l.Dimensions{}
				}
			},
		).
		Fn(gtx)
}

func (m *MinerModel) FillSpace(gtx l.Context) l.Dimensions {
	return l.Dimensions{
		Size: image.Point{
			X: gtx.Constraints.Min.X,
			Y: gtx.Constraints.Min.Y,
		},
		Baseline: 0,
	}
}

func (m *MinerModel) VSpacer(gtx l.Context) l.Dimensions {
	return l.Dimensions{
		Size: image.Point{
			X: int(m.TextSize.Scale(2).V),
			Y: int(m.TextSize.Scale(2).V),
		},
		Baseline: 0,
	}
}

func (m *MinerModel) Header(gtx l.Context) l.Dimensions {
	return m.Fill("Primary", l.Center, 0, 0, m.Flex().Rigid(
		m.Inset(
			0.25,
			m.IconButton(m.logoButton).
				Color("Light").
				Background("").
				Icon(m.Icon().Color("Light").Scale(gui.Scales["H5"]).Src(&icons.ParallelCoin)).
				Fn,
		).Fn,
	).Rigid(
		m.Inset(
			0.5,
			m.H5("kopach").
				Color("Light").
				Fn,
		).Fn,
	).Flexed(
		1,
		m.Inset(
			0.5,
			m.Body1(fmt.Sprintf("%d hash/s", int(m.worker.hashrate))).
				Color("DocBg").
				Alignment(text.End).
				Fn,
		).Fn,
	).Fn).Fn(gtx)
}

func (m *MinerModel) RunControl(gtx l.Context) l.Dimensions {
	return m.Inset(
		0.25,
		m.Flex().Flexed(
			0.5,
			m.Body1("enable mining").
				Color("DocText").
				Fn,
		).Flexed(
			0.5,
			m.Switch(
				m.mineToggle.SetOnChange(
					func(b bool) {
						if b {
							dbg.Ln("start mining")
							m.worker.StartChan <- struct{}{}
						} else {
							dbg.Ln("stop mining")
							m.worker.StopChan <- struct{}{}
						}
					},
				),
			).
				Fn,
		).Fn,
	).Fn(gtx)
}

func (m *MinerModel) SetThreads(gtx l.Context) l.Dimensions {
	return m.Flex().Rigid(
		m.Inset(
			0.25,
			m.Flex().
				Flexed(
					0.5,
					m.Body1(
						"number of mining threads"+
							fmt.Sprintf("%3v", m.threadSlider.GetValue()),
					).
						Fn,
				).
				Flexed(
					0.5,
					m.threadSlider.Fn,
				).
				Fn,
		).Fn,
	).Fn(gtx)
}

func (m *MinerModel) PreSharedKey(gtx l.Context) l.Dimensions {
	return m.Inset(
		0.25,
		m.Flex().Flexed(
			0.5,
			m.Body1("cluster preshared key").
				Color("DocText").
				Fn,
		).Flexed(
			0.5,
			m.password.Fn,
		).Fn,
	).Fn(gtx)
}

func (m *MinerModel) BlockInfoModalCloser(gtx l.Context) l.Dimensions {
	return m.Button(
		m.modalScrim.SetClick(
			func() {
				m.modalOn = false
			},
		),
	).Background("Primary").Text("close").Fn(gtx)
}

var currentBlock SolutionData

func (m *MinerModel) BlockDetails(gtx l.Context) l.Dimensions {
	return m.Fill("DocBg", l.Center, 0, 0, m.VFlex().AlignMiddle().Rigid(
		m.Inset(
			0.5,
			m.H5("Block Information").Alignment(text.Middle).Color("DocText").Fn,
		).Fn,
	).Rigid(
		m.Inset(
			0.5,
			m.Flex().Rigid(
				m.VFlex().
					Rigid(m.H6("Height").Font("bariol bold").Fn).
					Rigid(m.H6("PoW Hash").Font("bariol bold").Fn).
					Rigid(m.H6("Algorithm").Font("bariol bold").Fn).
					Rigid(m.H6("Version").Font("bariol bold").Fn).
					Rigid(m.H6("Index Hash").Font("bariol bold").Fn).
					Rigid(m.H6("Prev Block").Font("bariol bold").Fn).
					Rigid(m.H6("Merkle Root").Font("bariol bold").Fn).
					Rigid(m.H6("Timestamp").Font("bariol bold").Fn).
					Rigid(m.H6("Bits").Font("bariol bold").Fn).
					Rigid(m.H6("Nonce").Font("bariol bold").Fn).
					Fn,
			).Rigid(
				m.VFlex().
					Rigid(
						m.Flex().AlignBaseline().
							Rigid(m.H6(" ").Font("bariol bold").Fn).
							Rigid(m.Body1(fmt.Sprintf("%d", currentBlock.height)).Fn).
							Fn,
					).
					Rigid(
						m.Flex().AlignBaseline().
							Rigid(m.H6(" ").Font("bariol bold").Fn).
							Rigid(
								m.Caption(fmt.Sprintf("%s", currentBlock.hash)).Font("go regular").Fn,
							).Fn,
					).
					Rigid(
						m.Flex().AlignBaseline().
							Rigid(m.H6(" ").Font("bariol bold").Fn).
							Rigid(m.Body1(currentBlock.algo).Fn).
							Fn,
					).
					Rigid(
						m.Flex().AlignBaseline().
							Rigid(m.H6(" ").Font("bariol bold").Fn).
							Rigid(m.Body1(fmt.Sprintf("%d", currentBlock.version)).Fn).
							Fn,
					).
					Rigid(
						m.Flex().AlignBaseline().
							Rigid(m.H6(" ").Font("bariol bold").Fn).
							Rigid(
								m.Caption(fmt.Sprintf("%s", currentBlock.indexHash)).
									Font("go regular").Fn,
							).
							Fn,
					).
					Rigid(
						m.Flex().AlignBaseline().
							Rigid(m.H6(" ").Font("bariol bold").Fn).
							Rigid(
								m.Caption(fmt.Sprintf("%s", currentBlock.prevBlock)).
									Font("go regular").
									Fn,
							).Fn,
					).
					Rigid(
						m.Flex().AlignBaseline().
							Rigid(m.H6(" ").Font("bariol bold").Fn).
							Rigid(
								m.Caption(fmt.Sprintf("%s", currentBlock.merkleRoot)).
									Font("go regular").
									Fn,
							).Fn,
					).
					Rigid(
						m.Flex().AlignBaseline().
							Rigid(m.H6(" ").Font("bariol bold").Fn).
							Rigid(m.Body1(currentBlock.timestamp.Format(time.RFC3339)).Fn).Fn,
					).
					Rigid(
						m.Flex().
							AlignBaseline().
							Rigid(m.H6(" ").Font("bariol bold").Fn).
							Rigid(m.Body1(fmt.Sprintf("%x", currentBlock.bits)).Fn).Fn,
					).
					Rigid(
						m.Flex().AlignBaseline().
							Rigid(m.H6(" ").Font("bariol bold").Fn).
							Rigid(m.Body1(fmt.Sprintf("%d", currentBlock.nonce)).Fn).Fn,
					).Fn,
			).Fn,
		).Fn,
	).Rigid(
		m.Inset(
			0.5,
			m.BlockInfoModalCloser,
		).Fn,
	).Fn).Fn(gtx)
}

func (m *MinerModel) FoundBlocks(gtx l.Context) l.Dimensions {
	var widgets []l.Widget
	for x := range m.worker.solutions {
		i := x
		widgets = append(
			widgets, func(gtx l.Context) l.Dimensions {
				return m.Flex().
					Rigid(
						m.Button(
							m.solButtons[i].SetClick(
								func() {
									currentBlock = m.worker.solutions[i]
									dbg.Ln("clicked for block", currentBlock.height)
									m.modalWidget = m.BlockDetails
									m.modalOn = true
								},
							),
						).Color("DocBg").
							Text(fmt.Sprint(m.worker.solutions[i].height)).
							Inset(0.5).Fn,
					).Flexed(
					1,
					m.Inset(
						0.25,
						m.VFlex().
							Rigid(
								m.Flex().
									Rigid(
										m.Body1(m.worker.solutions[i].algo).Font("plan9").Fn,
									).
									Flexed(
										1,
										m.VFlex().
											Rigid(
												m.Body1(m.worker.solutions[i].hash).
													Font("go regular").
													TextScale(0.75).
													Alignment(text.End).
													Fn,
											).
											Rigid(
												m.Caption(
													fmt.Sprint(
														m.worker.solutions[i].time.Format(time.RFC3339),
													),
												).
													Alignment(text.End).
													Fn,
											).
											Fn,
									).Fn,
							).Fn,
					).Fn,
				).Fn(gtx)
			},
		)
	}
	return m.Inset(
		0.25,
		// m.Flex().Flexed(1,
		func(gtx l.Context) l.Dimensions {
			
			// dbg.S(widgets)
			return m.lists["found"].
				End().
				// ScrollWidth(int(m.Theme.TextSize.V * 3)).
				Vertical().
				Length(len(widgets)).
				ScrollToEnd().
				DisableScroll(false).
				ListElement(
					func(gtx l.Context, index int) l.Dimensions {
						return widgets[index](gtx)
					},
				).Fn(gtx)
			// Slice(gtx, widgets...)(gtx)
		},
		// ).Fn,
	).Fn(gtx)
}

// +build !headless

package monitor

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/pkg/util/interrupt"
)

type State struct {
	Ctx                       *conte.Xt
	Gtx                       *layout.Context
	Rc                        *rcd.RcVar
	Theme                     *gelook.DuoUItheme
	MainList                  *layout.List
	ModesList                 *layout.List
	CloseButton               *gel.Button
	LogoButton                *gel.Button
	RunMenuButton             *gel.Button
	StopMenuButton            *gel.Button
	PauseMenuButton           *gel.Button
	RestartMenuButton         *gel.Button
	SettingsFoldButton        *gel.Button
	RunModeFoldButton         *gel.Button
	BuildFoldButton           *gel.Button
	ModesButtons              map[string]*gel.Button
	RunMode                   string
	RunModeOpen               bool
	SettingsOpen              bool
	BuildOpen                 bool
	Running                   bool
	Pausing                   bool
	LightTheme                bool
	WindowWidth, WindowHeight int
}

func NewMonitor(cx *conte.Xt, gtx *layout.Context,
	rc *rcd.RcVar) *State {
	return &State{
		Ctx:   cx,
		Gtx:   gtx,
		Rc:    rc,
		Theme: gelook.NewDuoUItheme(),
		MainList: &layout.List{
			Axis: layout.Vertical,
		},
		ModesList: &layout.List{
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		},
		CloseButton:        new(gel.Button),
		LogoButton:         new(gel.Button),
		RunMenuButton:      new(gel.Button),
		StopMenuButton:     new(gel.Button),
		PauseMenuButton:    new(gel.Button),
		RestartMenuButton:  new(gel.Button),
		SettingsFoldButton: new(gel.Button),
		RunModeFoldButton:  new(gel.Button),
		BuildFoldButton:    new(gel.Button),
		ModesButtons: map[string]*gel.Button{
			"node":   new(gel.Button),
			"wallet": new(gel.Button),
			"shell":  new(gel.Button),
			"gui":    new(gel.Button),
		},
		RunMode:      "node",
		Running:      false,
		Pausing:      false,
		LightTheme:   true,
		WindowWidth:  0,
		WindowHeight: 0,
	}
}

func Run(cx *conte.Xt, rc *rcd.RcVar) (err error) {
	w := app.NewWindow(
		app.Size(unit.Dp(1600), unit.Dp(900)),
		app.Title("ParallelCoin Pod Monitor"),
	)
	gtx := layout.NewContext(w.Queue())
	mon := NewMonitor(cx, gtx, rc)
	go func() {
		L.Debug("starting up GUI event loop")
	out:
		for {
			select {
			case <-cx.KillAll:
				L.Debug("kill signal received")
				break out
			case e := <-w.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					L.Debug("destroy event received")
					close(cx.KillAll)
				case system.FrameEvent:
					gtx.Reset(e.Config, e.Size)
					cs := gtx.Constraints
					mon.WindowWidth, mon.WindowHeight =
						cs.Width.Max, cs.Height.Max
					TopLevelLayout(mon)()
					e.Frame(gtx.Ops)
				}
			}
		}
		L.Debug("gui shut down")
		os.Exit(0)
	}()
	// w.Invalidate()
	interrupt.AddHandler(func() {
		close(cx.KillAll)
	})
	app.Main()
	return
}

func TopLevelLayout(m *State) func() {
	return func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(m.Gtx,
			DuoUIheader(m),
			Body(m),
			BottomBar(m),
		)

	}
}

func Body(m *State) layout.FlexChild {
	return layout.Flexed(1, func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(m.Gtx, layout.Flexed(1, func() {
			cs := m.Gtx.Constraints
			gelook.DuoUIdrawRectangle(m.Gtx,
				cs.Width.Max, cs.Height.Max, m.Theme.Colors["DocBg"],
				[4]float32{0, 0, 0, 0},
				[4]float32{0, 0, 0, 0},
			)
		}),
		)
	})
}

func DuoUIheader(m *State) layout.FlexChild {
	return layout.Rigid(func() {
		layout.Flex{
			Axis:      layout.Horizontal,
			Spacing:   layout.SpaceBetween,
			Alignment: layout.Middle,
		}.Layout(m.Gtx,
			layout.Rigid(func() {
				cs := m.Gtx.Constraints
				gelook.DuoUIdrawRectangle(m.Gtx, cs.Width.Max,
					cs.Height.Max, m.Theme.Colors["PanelBg"],
					[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				var (
					textSize, iconSize               = 64, 64
					width, height                    = 72, 72
					paddingV, paddingH               = 8, 8
					insetSize, textInsetSize float32 = 16, 24
					closeInsetSize           float32 = 4
				)
				if m.WindowWidth < 1024 || m.WindowHeight < 1280 {
					textSize, iconSize = 24, 24
					width, height = 32, 32
					paddingV, paddingH = 8, 8
					insetSize = 10
					textInsetSize = 16
					closeInsetSize = 4
				}
				layout.Flex{Axis: layout.Horizontal}.Layout(m.Gtx,
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(insetSize)).Layout(m.Gtx,
							func() {
								var logoMeniItem gelook.DuoUIbutton
								logoMeniItem = m.Theme.DuoUIbutton(
									"", "",
									"", m.Theme.Colors["PanelBg"],
									"", "",
									"logo", m.Theme.Colors["PanelText"],
									textSize, iconSize,
									width, height,
									paddingV, paddingH)
								for m.LogoButton.Clicked(m.Gtx) {
									FlipTheme(m)
								}
								logoMeniItem.IconLayout(m.Gtx, m.LogoButton)
							},
						)
					}),
					layout.Flexed(1, func() {
						layout.UniformInset(unit.Dp(textInsetSize)).Layout(m.
							Gtx,
							func() {
								t := m.Theme.DuoUIlabel(unit.Dp(float32(
									textSize)),
									"monitor")
								t.Color = m.Theme.Colors["PanelText"]
								t.Layout(m.Gtx)
							},
						)
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(closeInsetSize*2)).Layout(
							m.Gtx,
							func() {
								t := m.Theme.DuoUIlabel(unit.Dp(float32(
									24)),
									fmt.Sprintf("%dx%d",
										m.Gtx.Constraints.Width.Max,
										m.Gtx.Constraints.Height.Max))
								t.Color = m.Theme.Colors["PanelBg"]
								t.Font.Typeface = m.Theme.Fonts["Primary"]
								t.Layout(m.Gtx)
							})
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(closeInsetSize)).Layout(
							m.Gtx, func() {
								m.Theme.DuoUIbutton("", "settings",
									m.Theme.Colors["PanelText"],
									"", "",
									m.Theme.Colors["PanelBg"], "closeIcon",
									m.Theme.Colors["PanelText"],
									0, 32, 41, 41,
									0, 0).IconLayout(m.Gtx, m.CloseButton)
								for m.CloseButton.Clicked(m.Gtx) {
									L.Debug("close button clicked")
									close(m.Ctx.KillAll)
								}
							})
					}),
				)
			}),
		)
	})
}

func FlipTheme(m *State) {
	if m.LightTheme {
		m.Theme.Colors["PanelText"] = m.Theme.Colors["Light"]
		m.Theme.Colors["PanelBg"] = m.Theme.Colors["Dark"]
		m.Theme.Colors["DocText"] = m.Theme.Colors["White"]
		m.Theme.Colors["DocBg"] = m.Theme.Colors["Black"]
	} else {
		m.Theme.Colors["PanelText"] = m.Theme.Colors["Light"]
		m.Theme.Colors["PanelBg"] = m.Theme.Colors["Dark"]
		m.Theme.Colors["DocText"] = m.Theme.Colors["Dark"]
		m.Theme.Colors["DocBg"] = m.Theme.Colors["Light"]
	}
	m.LightTheme = !m.LightTheme
}

// +build !headless

package monitor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	log "github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/util/interrupt"
)

type State struct {
	Ctx                       *conte.Xt
	Gtx                       *layout.Context
	Rc                        *rcd.RcVar
	Theme                     *gelook.DuoUItheme
	Config                    *Config
	MainList                  *layout.List
	ModesList                 *layout.List
	CloseButton               *gel.Button
	LogoButton                *gel.Button
	RunMenuButton             *gel.Button
	StopMenuButton            *gel.Button
	PauseMenuButton           *gel.Button
	RestartMenuButton         *gel.Button
	RunModeFoldButton         *gel.Button
	SettingsFoldButton        *gel.Button
	SettingsCloseButton       *gel.Button
	SettingsTitleCloseButton  *gel.Button
	BuildFoldButton           *gel.Button
	BuildCloseButton          *gel.Button
	BuildTitleCloseButton     *gel.Button
	ModesButtons              map[string]*gel.Button
	GroupsList                *layout.List
	Running                   bool
	Pausing                   bool
	WindowWidth, WindowHeight int
	Loggers                   *Node
}

const ConfigFileName = "monitor.json"

func (m *State) LoadConfig() {
	m.Config.Width, m.Config.Height = 800, 600
	filename := filepath.Join(*m.Ctx.Config.DataDir, ConfigFileName)
	if apputil.FileExists(filename) {
		b, err := ioutil.ReadFile(filename)
		if err == nil {
			err = json.Unmarshal(b, m.Config)
			if err != nil {
				L.Error("error unmarshalling config", err)
				os.Exit(1)
			}
		} else {
			L.Fatal("unexpected error reading configuration file:", err)
			os.Exit(1)
		}
	}
	m.SetTheme(m.Config.DarkTheme)
}

func (m *State) SaveConfig() {
	m.Config.Width, m.Config.Height = m.WindowWidth, m.WindowHeight
	filename := filepath.Join(*m.Ctx.Config.DataDir, ConfigFileName)
	if yp, e := json.MarshalIndent(m.Config, "", "  "); e == nil {
		apputil.EnsureDir(filename)
		if e := ioutil.WriteFile(filename, yp, 0600); e != nil {
			L.Error(e)
		}
	}
}

type Config struct {
	Width, Height int
	RunMode       string
	RunModeOpen   bool
	SettingsOpen  bool
	BuildOpen     bool
	DarkTheme     bool
}

func NewMonitor(cx *conte.Xt, gtx *layout.Context, rc *rcd.RcVar) *State {
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
		CloseButton:              new(gel.Button),
		LogoButton:               new(gel.Button),
		RunMenuButton:            new(gel.Button),
		StopMenuButton:           new(gel.Button),
		PauseMenuButton:          new(gel.Button),
		RestartMenuButton:        new(gel.Button),
		SettingsFoldButton:       new(gel.Button),
		RunModeFoldButton:        new(gel.Button),
		BuildFoldButton:          new(gel.Button),
		BuildCloseButton:         new(gel.Button),
		BuildTitleCloseButton:    new(gel.Button),
		SettingsCloseButton:      new(gel.Button),
		SettingsTitleCloseButton: new(gel.Button),
		ModesButtons: map[string]*gel.Button{
			"node":   new(gel.Button),
			"wallet": new(gel.Button),
			"shell":  new(gel.Button),
			"gui":    new(gel.Button),
		},
		Config: &Config{
			RunMode:   "node",
			DarkTheme: true,
		},
		Running:      false,
		Pausing:      false,
		WindowWidth:  0,
		WindowHeight: 0,
		GroupsList: &layout.List{
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		},
	}
}

func Run(cx *conte.Xt, rc *rcd.RcVar) (err error) {
	mon := NewMonitor(cx, nil, rc)
	var lgs []string
	for i := range log.Loggers {
		lgs = append(lgs, i)
	}
	mon.Loggers = GetTree(lgs)
	mon.LoadConfig()
	w := app.NewWindow(
		app.Size(unit.Dp(float32(mon.Config.Width)),
			unit.Dp(float32(mon.Config.Height))),
		app.Title("ParallelCoin Pod Monitor"),
	)
	mon.Gtx = layout.NewContext(w.Queue())
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
					mon.SaveConfig()
					close(cx.KillAll)
				case system.FrameEvent:
					mon.Gtx.Reset(e.Config, e.Size)
					cs := mon.Gtx.Constraints
					mon.WindowWidth, mon.WindowHeight =
						cs.Width.Max, cs.Height.Max
					mon.TopLevelLayout()
					e.Frame(mon.Gtx.Ops)
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

func (m *State) TopLevelLayout() {
	m.FlexV(
		m.DuoUIheader(),
		m.Body(),
		m.BottomBar(),
	)
}

func (m *State) DuoUIheader() layout.FlexChild {
	return Rigid(func() {
		m.FlexH(Rigid(func() {
			cs := m.Gtx.Constraints
			m.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
			var (
				textSize, iconSize       = 64, 64
				width, height            = 72, 72
				paddingV, paddingH       = 8, 8
				insetSize, textInsetSize = 16, 24
				closeInsetSize           = 4
			)
			if m.WindowWidth < 1024 || m.WindowHeight < 1280 {
				textSize, iconSize = 24, 24
				width, height = 32, 32
				paddingV, paddingH = 8, 8
				insetSize = 10
				textInsetSize = 16
				closeInsetSize = 4
			}
			m.FlexH(Rigid(func() {
				m.Inset(insetSize,
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
							m.FlipTheme()
							m.SaveConfig()
						}
						logoMeniItem.IconLayout(m.Gtx, m.LogoButton)
					},
				)
			}), Rigid(func() {
				m.Inset(textInsetSize, func() {
					t := m.Theme.DuoUIlabel(unit.Dp(float32(
						textSize)),
						"monitor")
					t.Color = m.Theme.Colors["PanelText"]
					t.Layout(m.Gtx)
				},
				)
			}), Spacer(), Rigid(func() {
				m.Inset(closeInsetSize*2, func() {
					t := m.Theme.DuoUIlabel(unit.Dp(float32(24)),
						fmt.Sprintf("%dx%d",
							m.WindowWidth,
							m.WindowHeight))
					t.Color = m.Theme.Colors["PanelText"]
					t.Font.Typeface = m.Theme.Fonts["Primary"]
					t.Layout(m.Gtx)
				})
			}), Rigid(func() {
				m.Inset(closeInsetSize, func() {
					m.IconButton("closeIcon", "PanelText",
						"PanelBg", m.CloseButton)
					for m.CloseButton.Clicked(m.Gtx) {
						L.Debug("close button clicked")
						m.SaveConfig()
						close(m.Ctx.KillAll)
					}
				})
			}),
			)
		}),
		)
	})
}

func (m *State) FlipTheme() {
	m.SetTheme(Toggle(&m.Config.DarkTheme))
}

func (m *State) SetTheme(dark bool) {
	if dark {
		m.Theme.Colors["DocText"] = m.Theme.Colors["Dark"]
		m.Theme.Colors["DocBg"] = m.Theme.Colors["Light"]
		m.Theme.Colors["PanelText"] = m.Theme.Colors["Dark"]
		m.Theme.Colors["PanelBg"] = m.Theme.Colors["White"]
		// m.Theme.Colors["Primary"] = m.Theme.Colors["Gray"]
		// m.Theme.Colors["Secondary"] = m.Theme.Colors["White"]
	} else {
		m.Theme.Colors["DocText"] = m.Theme.Colors["Light"]
		m.Theme.Colors["DocBg"] = m.Theme.Colors["Black"]
		m.Theme.Colors["PanelText"] = m.Theme.Colors["Light"]
		m.Theme.Colors["PanelBg"] = m.Theme.Colors["Dark"]
		// m.Theme.Colors["Primary"] = m.Theme.Colors["Dark"]
		// m.Theme.Colors["Secondary"] = m.Theme.Colors["Black"]
	}
}

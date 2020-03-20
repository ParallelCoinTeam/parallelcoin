package monitor

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"gioui.org/layout"
	"go.uber.org/atomic"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

const ConfigFileName = "monitor.json"

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
	KillMenuButton            *gel.Button
	RunModeFoldButton         *gel.Button
	SettingsFoldButton        *gel.Button
	SettingsCloseButton       *gel.Button
	SettingsTitleCloseButton  *gel.Button
	BuildFoldButton           *gel.Button
	BuildCloseButton          *gel.Button
	BuildTitleCloseButton     *gel.Button
	ModesButtons              map[string]*gel.Button
	GroupsList                *layout.List
	Running                   atomic.Bool
	Pausing                   atomic.Bool
	WindowWidth, WindowHeight int
	Loggers                   *Node
	SettingsFields            *layout.List
	RunningInRepo             bool
	RunningInRepoButton       *gel.Button
	RunFromProfileButton      *gel.Button
	HasGo                     bool
	HasOtherGo                bool
	UseBuiltinGoButton        *gel.Button
	InstallNewGoButton        *gel.Button
	CannotRun                 bool
	RunCommandChan            chan string
}

func NewMonitor(cx *conte.Xt, gtx *layout.Context, rc *rcd.RcVar) (s *State) {
	s = &State{
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
		KillMenuButton:           new(gel.Button),
		SettingsFoldButton:       new(gel.Button),
		RunModeFoldButton:        new(gel.Button),
		BuildFoldButton:          new(gel.Button),
		BuildCloseButton:         new(gel.Button),
		BuildTitleCloseButton:    new(gel.Button),
		SettingsCloseButton:      new(gel.Button),
		SettingsTitleCloseButton: new(gel.Button),
		ModesButtons: map[string]*gel.Button{
			"node":    new(gel.Button),
			"wallet":  new(gel.Button),
			"shell":   new(gel.Button),
			"gui":     new(gel.Button),
			"monitor": new(gel.Button),
		},
		Config:       &Config{},
		WindowWidth:  0,
		WindowHeight: 0,
		GroupsList: &layout.List{
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		},
		SettingsFields: &layout.List{
			Axis: layout.Vertical,
		},
		RunningInRepoButton:  new(gel.Button),
		RunFromProfileButton: new(gel.Button),
		UseBuiltinGoButton:   new(gel.Button),
		InstallNewGoButton:   new(gel.Button),
		RunCommandChan:       make(chan string),
	}
	s.Running.Store(false)
	s.Pausing.Store(false)
	s.Config.RunMode.Store("node")
	s.Config.DarkTheme.Store(true)
	return
}

type Config struct {
	Width, Height atomic.Int32
	RunMode       atomic.String
	RunModeOpen   atomic.Bool
	SettingsOpen  atomic.Bool
	BuildOpen     atomic.Bool
	DarkTheme     atomic.Bool
	RunInRepo     atomic.Bool
	UseBuiltinGo  atomic.Bool
}

func (c *Config) GetUnsafeConfig() (out *UnsafeConfig) {
	out = &UnsafeConfig{
		Width:        c.Width.Load(),
		Height:       c.Height.Load(),
		RunMode:      c.RunMode.Load(),
		RunModeOpen:  c.RunModeOpen.Load(),
		SettingsOpen: c.SettingsOpen.Load(),
		BuildOpen:    c.BuildOpen.Load(),
		DarkTheme:    c.DarkTheme.Load(),
		RunInRepo:    c.RunInRepo.Load(),
		UseBuiltinGo: c.UseBuiltinGo.Load(),
	}
	return
}

type UnsafeConfig struct {
	Width, Height int32
	RunMode       string
	RunModeOpen   bool
	SettingsOpen  bool
	BuildOpen     bool
	DarkTheme     bool
	RunInRepo     bool
	UseBuiltinGo  bool
}

func (u *UnsafeConfig) LoadInto(c *Config) {
	c.Width.Store(u.Width)
	c.Height.Store(u.Height)
	c.RunMode.Store(u.RunMode)
	c.RunModeOpen.Store(u.RunModeOpen)
	c.SettingsOpen.Store(u.SettingsOpen)
	c.BuildOpen.Store(u.BuildOpen)
	c.DarkTheme.Store(u.DarkTheme)
	c.RunInRepo.Store(u.RunInRepo)
	c.UseBuiltinGo.Store(u.UseBuiltinGo)
}

func (s *State) LoadConfig() {
	L.Debug("loading config")
	var err error
	u := new(UnsafeConfig)
	u.Width, u.Height = 800, 600
	u.RunMode = "node"
	L.Debugs(u)
	filename := filepath.Join(*s.Ctx.Config.DataDir, ConfigFileName)
	if apputil.FileExists(filename) {
		L.Debug("config file exists")
		var b []byte
		if b, err = ioutil.ReadFile(filename); !L.Check(err) {
			L.Warn(string(b))
			if err = json.Unmarshal(b, u); L.Check(err) {
				u.LoadInto(s.Config)
				L.Debugs(s.Config)
				s.SaveConfig()
			}
			u.LoadInto(s.Config)
			L.Debugs(s.Config)
		}
	} else {
		L.Warn("creating new configuration")
		u.LoadInto(s.Config)
		L.Debugs(s.Config)
		s.SaveConfig()
	}
	s.SetTheme(u.DarkTheme)
}

func (s *State) SaveConfig() {
	// L.Debug("saving config")
	filename := filepath.Join(*s.Ctx.Config.DataDir, ConfigFileName)
	u := s.Config.GetUnsafeConfig()
	// L.Debugs(u)
	if yp, e := json.MarshalIndent(u, "", "  "); !L.Check(e) {
		L.Debug(string(yp))
		apputil.EnsureDir(filename)
		if e := ioutil.WriteFile(filename, yp, 0600); L.Check(e) {
			// panic(e)
		}
		u.LoadInto(s.Config)
		// b, err := ioutil.ReadFile(filename)
		// if string(b) != string(yp) {
		// 	L.Fatal(err)
		// 	panic(err)
		// }
	}
}

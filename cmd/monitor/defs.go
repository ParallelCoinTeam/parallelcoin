package monitor

import (
	"encoding/json"
	"gioui.org/app"
	"io/ioutil"
	"path/filepath"

	"gioui.org/layout"
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
	W                         *app.Window
	Rc                        *rcd.RcVar
	Theme                     *gelook.DuoUItheme
	Config                    *Config
	MainList                  *layout.List
	ModesList                 *layout.List
	CloseButton               *gel.Button
	RestartButton             *gel.Button
	LogoButton                *gel.Button
	RunMenuButton             *gel.Button
	StopMenuButton            *gel.Button
	PauseMenuButton           *gel.Button
	RestartMenuButton         *gel.Button
	KillMenuButton            *gel.Button
	RunModeFoldButton         *gel.Button
	SettingsFoldButton        *gel.Button
	SettingsCloseButton       *gel.Button
	SettingsZoomButton        *gel.Button
	SettingsTitleCloseButton  *gel.Button
	BuildFoldButton           *gel.Button
	BuildCloseButton          *gel.Button
	BuildZoomButton           *gel.Button
	BuildTitleCloseButton     *gel.Button
	FilterButton              *gel.Button
	FilterHeaderButton        *gel.Button
	FilterAllButton           *gel.Button
	FilterHideButton          *gel.Button
	FilterShowButton          *gel.Button
	FilterNoneButton          *gel.Button
	ModesButtons              map[string]*gel.Button
	GroupsList                *layout.List
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
	FilterButtons             []*gel.Button
	FilterList                layout.List
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
		RestartButton:            new(gel.Button),
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
		BuildZoomButton:          new(gel.Button),
		BuildTitleCloseButton:    new(gel.Button),
		SettingsCloseButton:      new(gel.Button),
		SettingsZoomButton:       new(gel.Button),
		SettingsTitleCloseButton: new(gel.Button),
		FilterButton:             new(gel.Button),
		FilterHeaderButton:       new(gel.Button),
		FilterAllButton:          new(gel.Button),
		FilterNoneButton:         new(gel.Button),
		FilterHideButton:         new(gel.Button),
		FilterShowButton:         new(gel.Button),
		ModesButtons: map[string]*gel.Button{
			"node":   new(gel.Button),
			"wallet": new(gel.Button),
			"shell":  new(gel.Button),
			"gui":    new(gel.Button),
			"mon":    new(gel.Button),
		},
		Config:       &Config{FilterNodes: make(map[string]*Node)},
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
	s.Config.RunMode = "node"
	s.Config.DarkTheme = true
	return
}

type TreeNode struct {
	Closed, Hidden bool
}

type Config struct {
	Width          int
	Height         int
	RunMode        string
	RunModeOpen    bool
	RunModeZoomed  bool
	SettingsOpen   bool
	SettingsZoomed bool
	SettingsTab    string
	BuildOpen      bool
	BuildZoomed    bool
	DarkTheme      bool
	RunInRepo      bool
	UseBuiltinGo   bool
	Running        bool
	Pausing        bool
	FilterOpen     bool
	FilterNodes    map[string]*Node
}

func (s *State) LoadConfig() {
	L.Debug("loading config")
	var err error
	cnf := &Config{}
	//u.Width, u.Height = 800, 600
	//u.RunMode = "node"
	//L.Debugs(u)
	filename := filepath.Join(*s.Ctx.Config.DataDir, ConfigFileName)
	if apputil.FileExists(filename) {
		//L.Debug("config file exists")
		var b []byte
		if b, err = ioutil.ReadFile(filename); !L.Check(err) {
			//L.Warn(string(b))
			//L.Debug(s.Config)
			if err = json.Unmarshal(b, cnf); L.Check(err) {
				//s.SaveConfig()
			}
			for i := range cnf.FilterNodes {
				s.Config.FilterNodes[i].Hidden = cnf.FilterNodes[i].Hidden
				s.Config.FilterNodes[i].Closed = cnf.FilterNodes[i].Closed
			}
			s.Config.Width = cnf.Width
			s.Config.Height = cnf.Height
			s.Config.RunMode = cnf.RunMode
			s.Config.RunModeOpen = cnf.RunModeOpen
			s.Config.RunModeZoomed = cnf.RunModeZoomed
			s.Config.SettingsOpen = cnf.SettingsOpen
			s.Config.SettingsZoomed = cnf.SettingsZoomed
			s.Config.SettingsTab = cnf.SettingsTab
			s.Config.BuildOpen = cnf.BuildOpen
			s.Config.BuildZoomed = cnf.BuildZoomed
			s.Config.DarkTheme = cnf.DarkTheme
			s.Config.RunInRepo = cnf.RunInRepo
			s.Config.UseBuiltinGo = cnf.UseBuiltinGo
			s.Config.Running = cnf.Running
			s.Config.Pausing = cnf.Pausing
			s.Config.FilterOpen = cnf.FilterOpen
			//L.Debugs(s.Config.FilterNodes)
			//for i := range s.Config.FilterNodes {
			//	L.Debug(s.Config.FilterNodes[i])
			//}
		}
	} else {
		L.Warn("creating new configuration")
		s.Config.UseBuiltinGo = s.HasGo
		s.Config.RunInRepo = s.RunningInRepo
		//L.Debugs(s.Config)
		s.SaveConfig()
	}
	if s.Config.Width < 1 || s.Config.Height < 1 {
		s.Config.Width = 800
		s.Config.Height = 600
	}
	if s.Config.SettingsTab == "" {
		s.Config.SettingsTab = "config"
	}
	s.Rc.Settings.Tabs.Current = s.Config.SettingsTab
	s.SetTheme(s.Config.DarkTheme)
}

func (s *State) SaveConfig() {
	//L.Debug("saving config")
	//for i, v := range s.Config.FilterNodes {
	//	L.Debug(i, v)
	//}
	//L.Debugs(s.Config.FilterNodes)
	s.Config.Width = s.WindowWidth
	s.Config.Height = s.WindowHeight
	filename := filepath.Join(*s.Ctx.Config.DataDir, ConfigFileName)
	//L.Debugs(u)
	if yp, e := json.MarshalIndent(s.Config, "", "  "); !L.Check(e) {
		//L.Debug(string(yp))
		apputil.EnsureDir(filename)
		if e := ioutil.WriteFile(filename, yp, 0600); L.Check(e) {
			panic(e)
		}
	}
}

package monitor

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"gioui.org/layout"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/data/ring"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
	"github.com/p9c/pod/pkg/util/logi"
)

const ConfigFileName = "monitor.json"

// State stores the state of the monitor
type State struct {
	gui.State
	Ctx                       *conte.Xt
	Worker                    *worker.Worker
	Config                    *Config
	Buttons                   map[string]*gel.Button
	FilterLevelsButtons       []gel.Button
	FilterButtons             []gel.Button
	Lists                     map[string]*layout.List
	ModesButtons              map[string]*gel.Button
	CheckBoxes                map[string]*gel.CheckBox
	CommandEditor             gel.Editor
	WindowWidth, WindowHeight int
	Loggers                   *Node
	RunningInRepo             bool
	HasGo                     bool
	HasOtherGo                bool
	CannotRun                 bool
	RunCommandChan            chan string
	EntryBuf                  *ring.Entry
	FilterBuf                 *ring.Entry
	FilterFunc                func(ent *logi.Entry) bool
	FilterRoot                *Node
}

func NoFilter(_ *logi.Entry) bool { return true }

func (s *State) FilterOn(ent *logi.Entry) (out bool) {
	if s.Config.FilterNodes == nil || ent == nil {
		return true
	}
	if x, ok := s.Config.FilterNodes[ent.Package]; ok {
		if !x.Hidden {
			out = true
		}
	}
	cfgLevel := 0
	level := 0
	for i := range logi.Levels {
		if *s.Ctx.Config.LogLevel == logi.Levels[i] {
			cfgLevel = i
		}
		if ent.Level == logi.Levels[i] {
			level = i
		}
	}
	if level <= cfgLevel {
		out = true
	}
	return
}

func NewMonitor(cx *conte.Xt, gtx *layout.Context, rc *rcd.RcVar) (s *State) {
	s = &State{
		Ctx: cx,
		State: gui.State{
			Gtx:   gtx,
			Rc:    rc,
			Theme: gelook.NewDuoUItheme(),
		},
		ModesButtons:        make(map[string]*gel.Button),
		Config:              &Config{FilterNodes: make(map[string]*Node)},
		WindowWidth:         800,
		WindowHeight:        600,
		RunCommandChan:      make(chan string),
		EntryBuf:            ring.NewEntry(65536),
		FilterBuf:           ring.NewEntry(65536),
		FilterFunc:          NoFilter,
		FilterLevelsButtons: make([]gel.Button, 7),
		Buttons:             make(map[string]*gel.Button),
		Lists:               make(map[string]*layout.List),
	}
	modes := []string{
		"node", "wallet", "shell", "gui", "mon",
	}
	for i := range modes {
		s.ModesButtons[modes[i]] = new(gel.Button)
	}
	buttons := []string{
		"Close",
		"Restart",
		"Logo",
		"RunMenu",
		"StopMenu",
		"PauseMenu",
		"RestartMenu",
		"KillMenu",
		"RunModeFold",
		"SettingsFold",
		"SettingsClose",
		"SettingsZoom",
		"BuildFold",
		"BuildClose",
		"BuildZoom",
		"BuildTitleClose",
		"Filter",
		"FilterHeader",
		"FilterAll",
		"FilterHide",
		"FilterShow",
		"FilterNone",
		"FilterClear",
		"FilterSend",
		"RunningInRepo",
		"RunFromProfile",
		"UseBuiltinGo",
		"InstallNewGo",
	}
	for i := range buttons {
		s.Buttons[buttons[i]] = new(gel.Button)
	}
	checkboxes := []string{"FilterMode"}
	s.CheckBoxes = make(map[string]*gel.CheckBox)
	for i := range checkboxes {
		s.CheckBoxes[checkboxes[i]] = new(gel.CheckBox)
	}
	lists := []string{
		"Modes", "FilterLevel", "Groups", "Filter", "Log",
		"SettingsFields",
	}
	for i := range lists {
		s.Lists[lists[i]] = new(layout.List)
	}
	s.Lists = map[string]*layout.List{
		"Modes": {
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		},
		"Groups": {
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		},
		"SettingsFields": {
			Axis: layout.Vertical,
		},
		"Filter":      {},
		"FilterLevel": {},
		"Log":         {},
	}
	s.Config.RunMode = "node"
	s.Config.DarkTheme = true
	return
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
	FilterLevel    int
	FilterMode     bool
	ClickCommand   string
}

func (s *State) LoadConfig() (isNew bool) {
	Debug("loading config")
	var err error
	cnf := &Config{}
	filename := filepath.Join(*s.Ctx.Config.DataDir, ConfigFileName)
	if apputil.FileExists(filename) {
		var b []byte
		if b, err = ioutil.ReadFile(filename); !Check(err) {
			if err = json.Unmarshal(b, cnf); Check(err) {
				s.SaveConfig()
			}
			if s.Config.FilterNodes == nil {
				s.Config.FilterNodes = make(map[string]*Node)
			}
			for i := range cnf.FilterNodes {
				if s.Config.FilterNodes[i] == nil {
					s.Config.FilterNodes[i] = &Node{}
				}
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
			s.Config.FilterLevel = cnf.FilterLevel
			s.Config.FilterMode = cnf.FilterMode
			s.Config.ClickCommand = cnf.ClickCommand
			s.Config.FilterMode = cnf.FilterMode
			s.CommandEditor.SetText(s.Config.ClickCommand)
		}
	} else {
		Warn("creating new configuration")
		s.Config.UseBuiltinGo = s.HasGo
		s.Config.RunInRepo = s.RunningInRepo
		isNew = true
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
	return
}

func (s *State) SaveConfig() {
	Debug("saving monitor config")
	s.Config.Width = s.WindowWidth
	s.Config.Height = s.WindowHeight
	filename := filepath.Join(*s.Ctx.Config.DataDir, ConfigFileName)
	if cfgJSON, e := json.MarshalIndent(s.Config, "", "  "); !Check(e) {
		// Debug(string(cfgJSON))
		apputil.EnsureDir(filename)
		if e := ioutil.WriteFile(filename, cfgJSON, 0600); Check(e) {
			panic(e)
		}
	}
}

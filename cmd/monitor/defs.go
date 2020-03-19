package monitor

import (
	"encoding/json"
	"io/ioutil"
	"os"
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
	SettingsFields            *layout.List
	RunningInRepo             bool
	RunningInRepoButton       *gel.Button
	RunFromProfileButton      *gel.Button
	HasGo                     bool
	HasOtherGo                bool
	UseBuiltinGoButton        *gel.Button
	InstallNewGoButton        *gel.Button
	CannotRun                 bool
}

type Config struct {
	Width, Height int
	RunMode       string
	RunModeOpen   bool
	SettingsOpen  bool
	BuildOpen     bool
	DarkTheme     bool
	RunInRepo     bool
	UseBuiltinGo  bool
}

func (st *State) LoadConfig() {
	st.Config.Width, st.Config.Height = 800, 600
	filename := filepath.Join(*st.Ctx.Config.DataDir, ConfigFileName)
	if apputil.FileExists(filename) {
		b, err := ioutil.ReadFile(filename)
		if err == nil {
			err = json.Unmarshal(b, st.Config)
			if err != nil {
				L.Error("error unmarshalling config", err)
				os.Exit(1)
			}
		} else {
			L.Fatal("unexpected error reading configuration file:", err)
			os.Exit(1)
		}
	}
	st.SetTheme(st.Config.DarkTheme)
}

func (st *State) SaveConfig() {
	st.Config.Width, st.Config.Height = st.WindowWidth, st.WindowHeight
	filename := filepath.Join(*st.Ctx.Config.DataDir, ConfigFileName)
	if yp, e := json.MarshalIndent(st.Config, "", "  "); e == nil {
		apputil.EnsureDir(filename)
		if e := ioutil.WriteFile(filename, yp, 0600); e != nil {
			L.Error(e)
		}
	}
}

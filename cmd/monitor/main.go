// +build !headless

package monitor

import (
	"os"
	"os/exec"
	"strings"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gopkg.in/src-d/go-git.v4"

	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	log "github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/util/interrupt"
)

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
		SettingsFields: &layout.List{
			Axis: layout.Vertical,
		},
		RunningInRepoButton: new(gel.Button),
		RunFromProfileButton: new(gel.Button),
		UseBuiltinGoButton: new(gel.Button),
		InstallNewGoButton: new(gel.Button),
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
	_, _ = git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      "https://github.com/src-d/go-git",
		Progress: os.Stderr,
	})
	var cwd string
	if cwd, err = os.Getwd(); L.Check(err) {
	}
	var repo *git.Repository
	if repo, err = git.PlainOpen(cwd); L.Check(err) {
	}
	if repo != nil {
		L.Debug("running inside repo")
		mon.RunningInRepo = true
		L.Debug(repo.Remotes())
	}
	cmd := exec.Command("go", "version")
	var out []byte
	out, err = cmd.CombinedOutput()
	if !strings.HasPrefix("go version", string(out)) {
		mon.HasGo = true
	}
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

func (st *State) TopLevelLayout() {
	st.FlexV(
		st.DuoUIheader(),
		st.Body(),
		st.BottomBar(),
	)
}

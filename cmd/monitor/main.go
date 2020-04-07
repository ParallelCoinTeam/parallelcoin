package monitor

import (
	"gioui.org/app"
	"gioui.org/app/headless"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/clipboard"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Run(cx *conte.Xt, rc *rcd.RcVar) (err error) {
	clipboard.Start()

	mon := NewMonitor(cx, nil, rc)
	var lgs []string
	for i := range *logi.L.Packages {
		lgs = append(lgs, i)
	}
	//Debugs(lgs)
	mon.Loggers = mon.GetTree(lgs)
	isNew := mon.LoadConfig()
	_, _ = git.PlainClone("/tmp/foo", false,
		&git.CloneOptions{
			URL:      "https://github.com/src-d/go-git",
			Progress: os.Stderr,
		})
	var cwd string
	if cwd, err = os.Getwd(); Check(err) {
	}
	var repo *git.Repository
	if repo, err = git.PlainOpen(cwd); Check(err) {
	}
	if repo != nil {
		Debug("running inside repo")
		mon.RunningInRepo = true
		Debug(repo.Remotes())
		if isNew {
			mon.Config.RunInRepo = true
		}
	}
	cmd := exec.Command("go", "version")
	var out []byte
	out, err = cmd.CombinedOutput()
	if !strings.HasPrefix("go version", string(out)) {
		mon.HasGo = true
		if isNew {
			mon.Config.UseBuiltinGo = true
		}
	}
	mon.W = app.NewWindow(
		app.Size(unit.Dp(float32(mon.Config.Width)),
			unit.Dp(float32(mon.Config.Height))),
		app.Title("ParallelCoin Pod Monitor ["+*cx.Config.DataDir+"]"),
	)
	mon.Gtx = layout.NewContext(mon.W.Queue())
	go mon.Runner()
	if mon.Config.Running && !(mon.Config.RunMode == "m" ||
		mon.Config.RunMode == "mon" || mon.Config.RunMode == "monitor") {
		go func() {
			Debug("starting up as was running previously when shut down")
			time.Sleep(time.Second / 2)
			mon.Config.Running = false
			//mon.RunCommandChan <- "stop"
			mon.RunCommandChan <- "run"
			if mon.Config.Pausing {
				time.Sleep(time.Second / 2)
				mon.RunCommandChan <- "pause"
			}
		}()
	}
	//go mon.Consume()
	var prevH, prevW int
	lastChanged := time.Now()
	go func() {
		Debug("starting up GUI event loop")
	out:
		for {
			select {
			case <-cx.KillAll:
				Debug("kill signal received")
				mon.SaveConfig()
				mon.RunCommandChan <- "kill"
				break out
			case e := <-mon.W.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					Debug("destroy event received")
					mon.SaveConfig()
					close(mon.Ctx.KillAll)
				case system.FrameEvent:
					mon.Gtx.Reset(e.Config, e.Size)
					// update config and gui state for window so everything is
					// correctly sized (gui needs it internally and when the
					// app closes it saves this value for next run)
					cs := mon.Gtx.Constraints
					w, h := cs.Width.Max, cs.Height.Max
					mon.WindowWidth, mon.WindowHeight = w, h
					mon.Config.Width, mon.Config.Height = w, h
					mon.TopLevelLayout(false)
					e.Frame(mon.Gtx.Ops)
					if w != prevW || h != prevH {
						if time.Now().Sub(lastChanged) > time.Second {
							if mon.HW, err = headless.NewWindow(w, h); Check(err) {
								return
							}
							lastChanged = time.Now()
						}
					}
				}
			}
		}
		mon.SaveConfig()
		mon.RunCommandChan <- "kill"
		Debug("gui shut down")
		os.Exit(0)
	}()
	interrupt.AddHandler(func() {
		mon.SaveConfig()
		mon.RunCommandChan <- "kill"
		close(mon.Ctx.KillAll)
	})
	app.Main()
	return
}

func (s *State) TopLevelLayout(headless bool) {
	//if !s.ScreenShooting {
	s.FlexV(
		s.DuoUIheader(headless),
		gui.Flexed(1, func() {
			s.FlexH(
				s.Body(headless),
				s.Sidebar(headless),
			)
		}),
		s.BottomBar(headless),
	)
	//} else {
	//	s.FlexV(gui.Rigid(func(){
	//		cs := s.Gtx.Constraints
	//		s.Rectangle(cs.Width.Max, cs.Height.Max, s.Theme.Colors["White"],
	//			"ff")
	//	}))
	//}
}

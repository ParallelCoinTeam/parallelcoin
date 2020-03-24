// +build !headless

package monitor

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/util/interrupt"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Run(cx *conte.Xt, rc *rcd.RcVar) (err error) {
	mon := NewMonitor(cx, nil, rc)
	var lgs []string
	for i := range logi.Loggers {
		lgs = append(lgs, i)
	}
	//L.Debugs(mon.Loggers)
	mon.Loggers = mon.GetTree(lgs)
	mon.LoadConfig()
	_, _ = git.PlainClone("/tmp/foo", false,
		&git.CloneOptions{
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
	mon.W = app.NewWindow(
		app.Size(unit.Dp(float32(mon.Config.Width)),
			unit.Dp(float32(mon.Config.Height))),
		app.Title("ParallelCoin Pod Monitor ["+*cx.Config.DataDir+"]"),
	)
	mon.Gtx = layout.NewContext(mon.W.Queue())
	go mon.Runner()
	if mon.Config.Running {
		go func() {
			time.Sleep(time.Second)
			mon.RunCommandChan <- "restart"
			if mon.Config.Pausing {
				mon.RunCommandChan <- "pause"
			}
		}()
	}
	go func() {
		L.Debug("starting up GUI event loop")
	out:
		for {
			select {
			case <-cx.KillAll:
				L.Debug("kill signal received")
				mon.SaveConfig()
				break out
			case e := <-mon.W.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					L.Debug("destroy event received")
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
	interrupt.AddHandler(func() {
		close(cx.KillAll)
	})
	app.Main()
	return
}

func (s *State) TopLevelLayout() {
	s.FlexV(
		s.DuoUIheader(),
		Flexed(1, func() {
			s.FlexH(Flexed(1, func() {
				s.FlexV(Flexed(1, func() {
					s.FlexH(
						s.Body(),
					)
				}), s.BottomBar(),
				)
			}), s.Sidebar(),
			)
		}),
	)
}

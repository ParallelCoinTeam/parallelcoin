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
	log "github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func Run(cx *conte.Xt, rc *rcd.RcVar) (err error) {
	mon := NewMonitor(cx, nil, rc)
	var lgs []string
	for i := range log.Loggers {
		lgs = append(lgs, i)
	}
	mon.Loggers = GetTree(lgs)
	mon.LoadConfig()
	mon.W = app.NewWindow(
		app.Size(unit.Dp(float32(mon.Config.Width.Load())),
			unit.Dp(float32(mon.Config.Height.Load()))),
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
	mon.Gtx = layout.NewContext(mon.W.Queue())
	go mon.Runner()
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
	// w.Invalidate()
	interrupt.AddHandler(func() {
		close(cx.KillAll)
	})
	app.Main()
	return
}

func (s *State) TopLevelLayout() {
	s.FlexV(
		s.DuoUIheader(),
		s.Body(),
		s.BottomBar(),
	)
}

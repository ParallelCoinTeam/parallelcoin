package monitor

import (
	"gioui.org/layout"
	"github.com/stalker-loki/app/slog"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/consume"
	"go.uber.org/atomic"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func (s *State) RunControls() layout.FlexChild {
	return gui.Rigid(func() {
		if s.CannotRun || s.Config.RunModeOpen {
			return
		}
		if !s.Config.Running {
			b := s.Buttons["RunMenu"]
			s.ButtonArea(func() {
				s.Gtx.Constraints.Width.Max = 48
				s.Gtx.Constraints.Height.Max = 48
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Height.Max, "DocText")
				s.Inset(8, func() {
					s.Icon("Run", "ButtonBg", "DocText", 32)
				})
			}, b)
			for b.Clicked(s.Gtx) {
				slog.Debug("clicked run button")
				if !s.Config.RunModeOpen {
					s.RunCommandChan <- "run"
				}
			}
		} else {
			ic := "Pause"
			fg, bg := "PanelBg", "PanelText"
			if s.Config.Pausing {
				ic = "Run"
				fg, bg = "PanelText", "PanelBg"
			}
			s.FlexH(gui.Rigid(func() {
				b := s.Buttons["StopMenu"]
				s.ButtonArea(func() {
					s.Gtx.Constraints.Width.Max = 48
					s.Gtx.Constraints.Height.Max = 48
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, bg)
					s.Inset(8, func() {
						s.Icon("Stop", fg, bg, 32)
					})
				}, b)
				for b.Clicked(s.Gtx) {
					slog.Debug("clicked stop button")
					s.RunCommandChan <- "stop"
				}
			}), gui.Rigid(func() {
				b := s.Buttons["PauseMenu"]
				s.ButtonArea(func() {
					s.Gtx.Constraints.Width.Max = 48
					s.Gtx.Constraints.Height.Max = 48
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, bg)
					s.Inset(8, func() {
						s.Icon(ic, fg, bg, 32)
					})
				}, b)
				// s.IconButton(ic, fg, bg, b)
				for b.Clicked(s.Gtx) {
					if s.Config.Pausing {
						slog.Debug("clicked on resume button")
						s.RunCommandChan <- "resume"
					} else {
						slog.Debug("clicked pause button")
						s.RunCommandChan <- "pause"
					}
				}
			}), gui.Rigid(func() {
				b := s.Buttons["RestartMenu"]
				s.ButtonArea(func() {
					s.Gtx.Constraints.Width.Max = 48
					s.Gtx.Constraints.Height.Max = 48
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, bg)
					s.Inset(8, func() {
						s.Icon("Restart", fg, bg, 32)
					})
				}, b)
				// s.IconButton("Restart", "PanelBg", "PanelText", b)
				for b.Clicked(s.Gtx) {
					slog.Debug("clicked restart button")
					s.RunCommandChan <- "restart"
				}
			}),
			)
		}
	})
}

func (s *State) Build() (exePath string, err error) {
	var c *exec.Cmd
	gt := "goterm"
	if runtime.GOOS == "windows" {
		gt = ""
	}
	exePath = filepath.Join(*s.Ctx.Config.DataDir, "pod_mon")
	c = exec.Command("go", "build", "-v",
		"-tags", gt, "-o", exePath)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err = c.Run(); !slog.Check(err) {
	}
	return
}

func (s *State) Runner() {
	logi.L.SetLevel(*s.Ctx.Config.LogLevel, true, "pod")
	var err error
	var exePath string
	var quit chan struct{}
	run := &atomic.Bool{}
	run.Store(false)
	for cmd := range s.RunCommandChan {
		switch cmd {
		case "run":
			slog.Debug("run called")
			if s.HasGo && !s.Config.Running {
				if exePath, err = s.Build(); !slog.Check(err) {
					quit = make(chan struct{})
					s.Worker = consume.Log(quit, func(ent *logi.Entry) (
						err error) {
						slog.Debugf("KOPACH %s %s", ent.Text, ent.Level)
						s.EntryBuf.Add(ent)
						if s.FilterFunc(ent) {
							s.FilterBuf.Add(ent)
						}
						return
					}, func(pkg string) (out bool) {
						if x, ok := s.Config.FilterNodes[pkg]; ok {
							if x.Hidden {
								return true
							}
						}
						return false
					}, exePath, "-D", *s.Ctx.Config.DataDir, "--pipelog",
						s.Config.RunMode)
					consume.Start(s.Worker)
					s.Config.Running = true
					s.Config.Pausing = false
					consume.SetFilter(s.Worker, s.FilterRoot.GetPackages())
					s.W.Invalidate()
					go func() {
						// time.Sleep(time.Second/10)
						if err = s.Worker.Wait(); !slog.Check(err) {
							s.Config.Running = false
							s.Config.Pausing = false
							s.W.Invalidate()
						}
					}()
				}
			}
		case "stop":
			slog.Debug("stop called")
			if s.HasGo && s.Worker != nil && s.Config.Running {
				close(quit)
				if err = s.Worker.Interrupt(); !slog.Check(err) {
					s.Config.Running = false
				}
			}
		case "pause":
			slog.Debug("pause called")
			if s.HasGo && s.Worker != nil && s.Config.Running && !s.Config.
				Pausing {
				s.Config.Pausing = !s.Config.Pausing
				consume.Stop(s.Worker)
				if err = s.Worker.Pause(); slog.Check(err) {
				}
			}
		case "resume":
			slog.Debug("resume called")
			if s.HasGo && s.Worker != nil && s.Config.Running && s.Config.
				Pausing {
				s.Config.Pausing = !s.Config.Pausing
				if err = s.Worker.Resume(); slog.Check(err) {
				}
				consume.Start(s.Worker)
			}
		case "kill":
			slog.Debug("kill called")
			if s.HasGo && s.Worker != nil && s.Config.Running {
				// close(quit)
				if err = s.Worker.Interrupt(); !slog.Check(err) {
				}
			}
		case "restart":
			slog.Debug("restart called")
			if s.HasGo && s.Worker != nil {
				go func() {
					s.RunCommandChan <- "stop"
					time.Sleep(time.Second)
					s.RunCommandChan <- "run"
				}()
			}
		}
	}
	return
}

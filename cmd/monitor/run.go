package monitor

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/logi/consume"
	"github.com/p9c/pod/pkg/stdconn/worker"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func (s *State) RunControls() layout.FlexChild {
	return Rigid(func() {
		if s.CannotRun {
			return
		}
		if !s.Config.Running {
			s.IconButton("Run", "PanelBg", "PanelText", s.RunMenuButton)
			for s.RunMenuButton.Clicked(s.Gtx) {
				L.Debug("clicked run button")
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
			s.FlexH(Rigid(func() {
				s.IconButton("Stop", "PanelBg", "PanelText",
					s.StopMenuButton)
				for s.StopMenuButton.Clicked(s.Gtx) {
					L.Debug("clicked stop button")
					s.RunCommandChan <- "stop"
				}
			}), Rigid(func() {
				s.IconButton(ic, fg, bg, s.PauseMenuButton)
				for s.PauseMenuButton.Clicked(s.Gtx) {
					if s.Config.Pausing {
						L.Debug("clicked on resume button")
						s.RunCommandChan <- "resume"
					} else {
						L.Debug("clicked pause button")
						s.RunCommandChan <- "pause"
					}
				}
				//}), Rigid(func() {
				//	s.IconButton("Kill", "PanelBg", "PanelText",
				//		s.KillMenuButton)
				//	for s.KillMenuButton.Clicked(s.Gtx) {
				//		L.Debug("clicked kill button")
				//		s.RunCommandChan <- "kill"
				//	}
			}), Rigid(func() {
				s.IconButton("Restart", "PanelBg", "PanelText",
					s.RestartMenuButton)
				for s.RestartMenuButton.Clicked(s.Gtx) {
					L.Debug("clicked restart button")
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
	if err = c.Run(); !L.Check(err) {
	}
	return
}

func (s *State) Runner() {
	var err error
	var exePath string
	var w *worker.Worker
	var quit chan struct{}
	for cmd := range s.RunCommandChan {
		switch cmd {
		case "run":
			L.Debug("run called")
			if s.HasGo && !s.Config.Running {
				if exePath, err = s.Build(); !L.Check(err) {
					quit = make(chan struct{})
					w = consume.Log(quit, func(ent *logi.Entry) (err error) {
						L.Debugf("KOPACH %s %s", ent.Text, ent.Level)
						return
					}, exePath, "-D", *s.Ctx.Config.DataDir, s.Config.RunMode)
					consume.Start(w)
					s.Config.Running = true
					s.Config.Pausing = false
					s.W.Invalidate()
					go func() {
						if err = w.Wait(); !L.Check(err) {
							s.Config.Running = false
							s.Config.Pausing = false
							s.W.Invalidate()
						}
					}()
				}
			}
		case "stop":
			L.Debug("stop called")
			if s.HasGo && w != nil && s.Config.Running {
				close(quit)
				if err = w.Interrupt(); !L.Check(err) {
					s.Config.Running = false
				}
			}
		case "pause":
			L.Debug("pause called")
			if s.HasGo && w != nil && s.Config.Running && !s.Config.Pausing {
				s.Config.Pausing = !s.Config.Pausing
				consume.Stop(w)
				if err = w.Pause(); L.Check(err) {
				}
			}
		case "resume":
			L.Debug("resume called")
			if s.HasGo && w != nil && s.Config.Running && s.Config.Pausing {
				s.Config.Pausing = !s.Config.Pausing
				if err = w.Resume(); L.Check(err) {
				}
				consume.Start(w)
			}
		case "kill":
			L.Debug("kill called")
			if s.HasGo && w != nil && s.Config.Running {
				close(quit)
				if err = w.Interrupt(); !L.Check(err) {
				}
			}
		case "restart":
			L.Debug("restart called")
			if s.HasGo && w != nil {
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

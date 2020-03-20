package monitor

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"gioui.org/layout"
)

func (s *State) RunControls() layout.FlexChild {
	return Rigid(func() {
		if s.CannotRun {
			return
		}
		if !s.Running.Load() {
			s.IconButton("Run", "PanelBg", "PanelText", s.RunMenuButton)
			for s.RunMenuButton.Clicked(s.Gtx) {
				L.Debug("clicked run button")
				if !s.Config.RunModeOpen.Load() {
					s.RunCommandChan <- "run"
				}
			}
		} else {
			ic := "Pause"
			fg, bg := "PanelBg", "PanelText"
			if s.Pausing.Load() {
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
					if s.Pausing.Load() {
						L.Debug("clicked on resume button")
					} else {
						L.Debug("clicked pause button")
					}
					s.RunCommandChan <- "pause"
				}
			}), Rigid(func() {
				s.IconButton("Kill", "PanelBg", "PanelText",
					s.KillMenuButton)
				for s.KillMenuButton.Clicked(s.Gtx) {
					L.Debug("clicked kill button")
					s.RunCommandChan <- "kill"
				}
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

func (s *State) Runner() {
	var c *exec.Cmd
	var err error
	// var pid int
	for cmd := range s.RunCommandChan {
		switch cmd {
		case "run":
			L.Debug("run called")
			if s.HasGo && !s.Running.Load() {
				c = exec.Command("go", "build", "-o",
					filepath.Join(*s.Ctx.Config.DataDir, "pod"))
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Run()
				c = exec.Command(filepath.Join(*s.Ctx.Config.DataDir, "pod"),
					"-D", *s.Ctx.Config.DataDir, s.Config.RunMode.Load())
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if err = c.Start(); !L.Check(err) {
					s.Running.Store(true)
					s.Pausing.Store(false)
					// pid = c.Process.Pid
					// go func() {
					// 	if err := c.Wait(); !L.Check(err) {
					// 		s.Running.Store(false)
					// 	}
					// }()
				}
			}
			continue
		case "stop":
			L.Debug("stop called")
			if s.HasGo && c != nil && s.Running.Load() {
				// if err = c.Process.Signal(os.Interrupt); !L.Check(err) {
				// 	s.Running.Store(false)
				// 	L.Debug("interrupted")
				// }
				// if err = syscall.Kill(pid,
				// 	syscall.SIGINT); !L.Check(err) {
				// 	s.Running.Store(false)
				// 	L.Debug("killed")
				// }
				if err = c.Process.Kill(); !L.Check(err) {
					L.Debug("killing harder")
				}
				if err = c.Wait(); L.Check(err) {
				}
				if err = c.Process.Release(); L.Check(err) {
				}
				L.Debug("dead")
			}
		case "pause":
			L.Debug("pause called")
			if s.HasGo && c != nil && s.Running.Load() {
				s.Pausing.Toggle()
			}
			continue
		case "kill":
			L.Debug("kill called")
			if s.HasGo && c != nil && s.Running.Load() {
				var pgid int
				if pgid, err = syscall.Getpgid(c.Process.Pid); L.Check(err) {
					// if err = syscall.Kill(-pgid, 15); L.Check(err) {
					// }
					if err = syscall.Kill(-pgid, 9); L.Check(err) {
					}
				}
			}
			continue
		case "restart":
			L.Debug("restart called")
			if s.HasGo && c != nil {
				if err = c.Process.Signal(os.Interrupt); L.Check(err) {
				}
				done := make(chan struct{})
				go func() {
					if err := c.Wait(); L.Check(err) {
						close(done)
					}
				}()
				select {
				case <-done:
				case <-time.After(time.Second * 5):
					var pgid int
					if pgid, err = syscall.Getpgid(c.Process.Pid); L.Check(err) {
						if err = syscall.Kill(-pgid, 15); L.Check(err) {
						}
					}
					if err = c.Process.Kill(); L.Check(err) {
					}
				}
				c = exec.Command("go", "run", "main.go", "-D",
					*s.Ctx.Config.DataDir, s.Config.RunMode.Load())
				c.Stdin = os.Stdin
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if err = c.Start(); !L.Check(err) {
					s.Running.Store(false)
					s.Pausing.Store(false)
				}
			}
			continue
		}
	}
	return
}

package monitor

import (
	"gioui.org/layout"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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
						s.RunCommandChan <- "resume"
					} else {
						L.Debug("clicked pause button")
						s.RunCommandChan <- "pause"
					}
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
				exePath := filepath.Join(*s.Ctx.Config.DataDir, "pod_mon")
				c = exec.Command("go", "build", "-o",
					exePath)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if err = c.Run(); !L.Check(err) {
					c = exec.Command(exePath,
						"-D", *s.Ctx.Config.DataDir, s.Config.RunMode.Load())
					c.Stdout = os.Stdout
					c.Stderr = os.Stderr
					if err = c.Start(); !L.Check(err) {
						s.Running.Store(true)
						s.Pausing.Store(false)
					}
				}
			}
		case "stop":
			L.Debug("stop called")
			if s.HasGo && c != nil && s.Running.Load() {
				if err = c.Process.Signal(syscall.SIGINT); !L.Check(err) {
					s.Running.Store(false)
					L.Debug("interrupted")
				}
				if err = c.Wait(); L.Check(err) {
				}
				if err = c.Process.Release(); L.Check(err) {
				}
				L.Debug("stopped")
			}
		case "pause":
			L.Debug("pause called")
			if s.HasGo && c != nil && s.Running.Load() {
				s.Pausing.Toggle()
				if err = c.Process.Signal(syscall.SIGSTOP); !L.Check(err) {
					s.Pausing.Store(true)
					L.Debug("paused")
				}
			}
		case "resume":
			L.Debug("resume called")
			if s.HasGo && c != nil && s.Running.Load() {
				s.Pausing.Toggle()
				if err = c.Process.Signal(syscall.SIGCONT); !L.Check(err) {
					s.Pausing.Store(false)
					L.Debug("resumed")
				}
			}
		case "kill":
			L.Debug("kill called")
			if s.HasGo && c != nil && s.Running.Load() {
				if err = c.Process.Signal(syscall.SIGKILL); !L.Check(err) {
					s.Pausing.Store(false)
					s.Running.Store(false)
					L.Debug("killed")
				}
			}
		case "restart":
			L.Debug("restart called")
			if s.HasGo && c != nil {
				if err = c.Process.Signal(syscall.SIGINT); !L.Check(err) {
					s.Running.Store(false)
					L.Debug("interrupted")
				}
				if err = c.Wait(); L.Check(err) {
				}
				if err = c.Process.Release(); L.Check(err) {
				}
				L.Debug("stopped")
			}
			exePath := filepath.Join(*s.Ctx.Config.DataDir, "pod_mon")
			c = exec.Command("go", "build", "-o",
				exePath)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err = c.Run(); !L.Check(err) {
				c = exec.Command(exePath,
					"-D", *s.Ctx.Config.DataDir, s.Config.RunMode.Load())
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if err = c.Start(); !L.Check(err) {
					s.Running.Store(true)
					s.Pausing.Store(false)
				}
			}
		}
	}
	return
}

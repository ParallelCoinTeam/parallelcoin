package monitor

import (
	"gioui.org/layout"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func (s *State) BottomBar() layout.FlexChild {
	return Rigid(func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
		s.FlexV(
			s.SettingsPage(),
			s.BuildPage(),
			s.StatusBar(),
		)
	})
}

func (s *State) StatusBar() layout.FlexChild {
	return Rigid(func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
		s.FlexH(
			s.RunControls(),
			s.RunmodeButtons(),
			s.BuildButtons(),
			s.SettingsButtons(),
			Spacer(),
			s.RestartRunButton(),
		)
	})
}

func (s *State) RunmodeButtons() layout.FlexChild {
	return Rigid(func() {
		s.FlexH(Rigid(func() {
			if !s.Config.RunModeOpen.Load() {
				fg, bg := "ButtonText", "ButtonBg"
				if s.Config.Running.Load() {
					fg, bg = "DocBg", "DocText"
				}
				s.TextButton(s.Config.RunMode.Load(), "Secondary",
					23, fg, bg,
					s.RunModeFoldButton)
				for s.RunModeFoldButton.Clicked(s.Gtx) {
					if !s.Config.Running.Load() {
						s.Config.RunModeOpen.Store(true)
						s.SaveConfig()
					}
				}
			} else {
				modes := []string{
					"node", "wallet", "shell", "gui", "monitor",
				}
				s.ModesList.Layout(s.Gtx, len(modes), func(i int) {
					s.TextButton(modes[i], "Secondary",
						23, "ButtonText",
						"ButtonBg", s.ModesButtons[modes[i]])
					for s.ModesButtons[modes[i]].Clicked(s.Gtx) {
						L.Debug(modes[i], "clicked")
						if s.Config.RunModeOpen.Load() {
							s.Config.RunMode.Store(modes[i])
							s.Config.RunModeOpen.Store(false)
						}
						s.SaveConfig()
					}
				})
			}
		}),
		)
	})
}

func (s *State) RestartRunButton() layout.FlexChild {
	return Rigid(func() {
		var c *exec.Cmd
		var err error
		s.IconButton("Restart", "PanelText", "PanelBg",
			s.RestartButton)
		for s.RestartButton.Clicked(s.Gtx) {
			L.Debug("clicked restart button")
			s.SaveConfig()
			if s.HasGo {
				go func() {
					s.RunCommandChan <- "stop"
					exePath := filepath.Join(*s.Ctx.Config.DataDir, "mon")
					c = exec.Command("go", "build", "-v",
						"-tags", "goterm", "-o", exePath)
					c.Stderr = os.Stderr
					c.Stdout = os.Stdout
					time.Sleep(time.Second)
					if err = c.Run(); !L.Check(err) {
						if err = syscall.Exec(exePath, os.Args,
							os.Environ()); L.Check(err) {
						}
						os.Exit(0)
					}
				}()
			}
		}
	})
}

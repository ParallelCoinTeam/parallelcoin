package monitor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
)

func (s *State) DuoUIheader() layout.FlexChild {
	return Rigid(func() {
		s.Gtx.Constraints.Height.Max = 48
		s.Gtx.Constraints.Height.Min = 48
		s.FlexH(Rigid(func() {
			s.FlexH(Rigid(func() {
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Width.Max, "PanelBg", "ff")
				s.FlexH(Rigid(func() {
					fg, bg := "PanelText", "PanelBg"
					icon := "logo"
					s.Theme.DuoUIbutton("", "", "",
						s.Theme.Colors[bg], "", s.Theme.Colors[fg], icon,
						s.Theme.Colors[fg], 0, 40, 48, 48,
						6, 2, 2, 6).IconLayout(s.Gtx, &s.LogoButton)
					if s.LogoButton.Clicked(s.Gtx) {
						s.FlipTheme()
					}
				}))
			}), Rigid(func() {
				s.FlexV(Flexed(1, func() {
					s.Inset(4, func() {
						layout.W.Layout(s.Gtx, func() {
							t := s.Theme.DuoUIlabel(unit.Dp(float32(40)), "Monitor")
							t.Color = s.Theme.Colors["PanelText"]
							t.Font.Typeface = s.Theme.Fonts["Secondary"]
							t.Layout(s.Gtx)
						})
					})
				}))
			}),
			)
		}), Spacer(), Rigid(func() {
			t := s.Theme.DuoUIlabel(unit.Dp(float32(16)),
				fmt.Sprintf("%s %dx%d", *s.Ctx.Config.DataDir,
					s.WindowWidth, s.WindowHeight))
			t.Color = s.Theme.Colors["PanelText"]
			t.Font.Typeface = s.Theme.Fonts["Primary"]
			t.Layout(s.Gtx)
		}), s.RestartRunButton(),
			Rigid(func() {
				s.IconButton("closeIcon", "PanelText",
					"PanelBg", &s.CloseButton)
				for s.CloseButton.Clicked(s.Gtx) {
					Debug("close button clicked")
					s.SaveConfig()
					s.RunCommandChan <- "kill"
					close(s.Ctx.KillAll)
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
			&s.RestartButton)
		for s.RestartButton.Clicked(s.Gtx) {
			Debug("clicked restart button")
			s.SaveConfig()
			if s.HasGo {
				s.RunCommandChan <- "kill"
				go func() {
					exePath := filepath.Join(*s.Ctx.Config.DataDir, "mon")
					c = exec.Command("go", "build", "-v",
						"-o", exePath)
					c.Stderr = os.Stderr
					c.Stdout = os.Stdout
					time.Sleep(time.Second)
					if err = c.Run(); !Check(err) {
						if err = syscall.Exec(exePath, os.Args,
							os.Environ()); Check(err) {
						}
						close(s.Ctx.KillAll)
						//time.Sleep(time.Second/2)
						//os.Exit(0)
					}
				}()
			}
		}
	})
}

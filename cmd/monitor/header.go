package monitor

import (
	"fmt"
	"github.com/p9c/pod/pkg/gui"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
)

func (s *State) DuoUIheader(headless bool) layout.FlexChild {
	gtx := s.Gtx
	if headless {
		gtx = s.Htx
	}
	return gui.Rigid(func() {
		gtx.Constraints.Height.Max = 48
		gtx.Constraints.Height.Min = 48
		s.FlexH(gui.Rigid(func() {
			s.FlexH(gui.Rigid(func() {
				cs := gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Width.Max, "PanelBg", "ff")
				s.FlexH(gui.Rigid(func() {
					fg, bg := "PanelText", "PanelBg"
					icon := "logo"
					b := s.Buttons["Logo"]
					s.Theme.DuoUIbutton("", "", "",
						s.Theme.Colors[bg], "", s.Theme.Colors[fg], icon,
						s.Theme.Colors[fg], 0, 40, 48, 48,
						6, 2, 2, 6).IconLayout(gtx, b)
					if b.Clicked(gtx) {
						s.FlipTheme(&s.Config.DarkTheme, s.SaveConfig)
					}
				}))
			}), gui.Rigid(func() {
				s.FlexV(gui.Flexed(1, func() {
					s.Inset(8, func() {
						layout.W.Layout(gtx, func() {
							t := s.Theme.DuoUIlabel(unit.Dp(float32(40)), "Monitor")
							t.Color = s.Theme.Colors["PanelText"]
							t.Font.Typeface = s.Theme.Fonts["Secondary"]
							t.Layout(gtx)
						})
					})
				}))
			}),
			)
		}), s.Spacer("PanelBg"), gui.Rigid(func() {
			t := s.Theme.DuoUIlabel(unit.Dp(float32(16)),
				fmt.Sprintf("%s %dx%d", *s.Ctx.Config.DataDir,
					s.WindowWidth, s.WindowHeight))
			t.Color = s.Theme.Colors["PanelText"]
			t.Font.Typeface = s.Theme.Fonts["Primary"]
			t.Layout(gtx)
		}), s.RestartRunButton(headless),
			gui.Rigid(func() {
				b := s.Buttons["Close"]
				s.IconButton("closeIcon", "PanelText",
					"PanelBg", b)
				for b.Clicked(gtx) {
					Debug("close button clicked")
					s.SaveConfig()
					s.RunCommandChan <- "kill"
					close(s.Ctx.KillAll)
				}
			}),
		)
	})
}

func (s *State) RestartRunButton(headless bool) layout.FlexChild {
	gtx := s.Gtx
	if headless {
		gtx = s.Htx
	}
	return gui.Rigid(func() {
		var c *exec.Cmd
		var err error
		b := s.Buttons["Restart"]
		s.IconButton("Restart", "PanelText", "PanelBg", b)
		for b.Clicked(gtx) {
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

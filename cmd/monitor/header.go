package monitor

import (
	"fmt"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/gelook"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
)

func (s *State) DuoUIheader() layout.FlexChild {
	return gui.Rigid(func() {
		s.Gtx.Constraints.Height.Max = 48
		s.Gtx.Constraints.Height.Min = 48
		s.FlexH(gui.Rigid(func() {
			s.FlexH(gui.Rigid(func() {
				cs := s.Gtx.Constraints
				s.Rectangle(cs.Width.Max, cs.Width.Max, "PanelBg")
				s.FlexH(gui.Rigid(func() {
					fg, bg := "PanelText", "PanelBg"
					icon := "logo"
					b := s.Buttons["Logo"]
					s.Theme.DuoUIbutton(gelook.ButtonParams{
						BgColor:       s.Theme.Colors[bg],
						BgHoverColor:  s.Theme.Colors[fg],
						Icon:          icon,
						IconColor:     s.Theme.Colors[fg],
						IconSize:      40,
						Width:         48,
						Height:        48,
						PaddingTop:    6,
						PaddingRight:  2,
						PaddingBottom: 2,
						PaddingLeft:   6,
					}).IconLayout(s.Gtx, b)
					if b.Clicked(s.Gtx) {
						s.FlipTheme()
					}
				}))
			}), gui.Rigid(func() {
				s.FlexV(gui.Flexed(1, func() {
					s.Inset(8, func() {
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
		}), s.Spacer("PanelBg"), gui.Rigid(func() {
			t := s.Theme.DuoUIlabel(unit.Dp(float32(16)),
				fmt.Sprintf("%s %dx%d", *s.Ctx.Config.DataDir,
					s.WindowWidth, s.WindowHeight))
			t.Color = s.Theme.Colors["PanelText"]
			t.Font.Typeface = s.Theme.Fonts["Primary"]
			t.Layout(s.Gtx)
		}), s.RestartRunButton(),
			gui.Rigid(func() {
				b := s.Buttons["Close"]
				s.IconButton("closeIcon", "PanelText",
					"PanelBg", b)
				for b.Clicked(s.Gtx) {
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
	return gui.Rigid(func() {
		var c *exec.Cmd
		var err error
		b := s.Buttons["Restart"]
		s.IconButton("Restart", "PanelText", "PanelBg", b)
		for b.Clicked(s.Gtx) {
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

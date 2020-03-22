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

	"github.com/p9c/pod/pkg/gelook"
)

func (s *State) DuoUIheader() layout.FlexChild {
	return Rigid(func() {
		s.FlexH(Rigid(func() {
			cs := s.Gtx.Constraints
			s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
			var (
				textSize, iconSize       = 64, 64
				width, height            = 72, 72
				paddingV, paddingH       = 8, 8
				insetSize, textInsetSize = 16, 24
				closeInsetSize           = 4
			)
			if s.WindowWidth < 1024 || s.WindowHeight < 1280 {
				textSize, iconSize = 24, 24
				width, height = 32, 32
				paddingV, paddingH = 8, 8
				insetSize = 10
				textInsetSize = 16
				closeInsetSize = 4
			}
			s.FlexH(Rigid(func() {
				s.Inset(insetSize,
					func() {
						var logoMeniItem gelook.DuoUIbutton
						logoMeniItem = s.Theme.DuoUIbutton(
							"", "",
							"", s.Theme.Colors["PanelBg"],
							"", "",
							"logo", s.Theme.Colors["PanelText"],
							textSize, iconSize,
							width, height,
							paddingV, paddingH)
						for s.LogoButton.Clicked(s.Gtx) {
							s.FlipTheme()
						}
						logoMeniItem.IconLayout(s.Gtx, s.LogoButton)
					},
				)
			}), Rigid(func() {
				s.Inset(textInsetSize, func() {
					t := s.Theme.DuoUIlabel(unit.Dp(float32(
						textSize)),
						"monitor")
					t.Color = s.Theme.Colors["PanelText"]
					t.Layout(s.Gtx)
				},
				)
			}), Spacer(), Rigid(func() {
				s.Inset(closeInsetSize*2, func() {
					t := s.Theme.DuoUIlabel(unit.Dp(float32(24)),
						fmt.Sprintf("%dx%d",
							s.WindowWidth,
							s.WindowHeight))
					t.Color = s.Theme.Colors["PanelText"]
					t.Font.Typeface = s.Theme.Fonts["Primary"]
					t.Layout(s.Gtx)
				})
			}), s.RestartRunButton(),
				Rigid(func() {
					s.Inset(closeInsetSize, func() {
						s.IconButton("closeIcon", "PanelText",
							"PanelBg", s.CloseButton)
						for s.CloseButton.Clicked(s.Gtx) {
							L.Debug("close button clicked")
							s.SaveConfig()
							s.RunCommandChan <- "stop"
							close(s.Ctx.KillAll)
						}
					})
				}),
			)
		}),
		)
	})
}

func (s *State) RestartRunButton() layout.FlexChild {
	return Rigid(func() {
		s.Inset(4, func() {
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
	})
}

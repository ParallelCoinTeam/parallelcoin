// Package monitor is a log viewer and filter and configuration interface
//
// +build !headless

package pkg

import (
	"fmt"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/gelook"
	"github.com/p9c/pkg/app/slog"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
)

func (s *State) Header() layout.FlexChild {
	bg, fg := "PanelBg", "PanelText"
	return gui.Rigid(func() {
		s.Gtx.Constraints.Height.Max = 48
		// s.Gtx.Constraints.Height.Min = 48
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Width.Max, bg)
		s.FlexH(gui.Rigid(func() {
			s.FlexH(gui.Rigid(func() {
				s.FlexH(gui.Rigid(func() {
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
							t := s.Theme.DuoUILabel(unit.Dp(float32(40)), "Monitor")
							t.Color = s.Theme.Colors[fg]
							t.Font.Typeface = s.Theme.Fonts["Secondary"]
							t.Layout(s.Gtx)
						})
					})
				}))
			}),
			)
		}),
			gui.Flexed(1, func() {
				// cs := s.Gtx.Constraints
				// s.Rectangle(cs.Width.Max, cs.Width.Max, "Primary")
				layout.E.Layout(s.Gtx, func() {
					t := s.Theme.DuoUILabel(unit.Dp(float32(16)),
						fmt.Sprintf("p2p %s rpc %s ctl %s %s %dx%d",
							(*s.Ctx.Config.Listeners)[0],
							(*s.Ctx.Config.RPCListeners)[0],
							*s.Ctx.Config.Controller,
							*s.Ctx.Config.DataDir,
							s.WindowWidth, s.WindowHeight))
					t.Color = s.Theme.Colors[fg]
					t.Font.Typeface = s.Theme.Fonts["Primary"]
					t.Layout(s.Gtx)
				})
			}),
			s.RestartRunButton(fg, bg),
			gui.Rigid(func() {
				b := s.Buttons["Close"]
				s.IconButton("closeIcon", fg,
					bg, b)
				for b.Clicked(s.Gtx) {
					slog.Debug("close button clicked")
					s.SaveConfig()
					s.RunCommandChan <- "kill"
					close(s.Ctx.KillAll)
				}
			}),
		)
	})
}

func (s *State) RestartRunButton(fg, bg string) layout.FlexChild {
	return gui.Rigid(func() {
		var c *exec.Cmd
		var err error
		b := s.Buttons["Restart"]
		s.IconButton("Restart", fg, bg, b)
		for b.Clicked(s.Gtx) {
			slog.Debug("clicked restart button")
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
					if err = c.Run(); !slog.Check(err) {
						if err = syscall.Exec(exePath, os.Args,
							os.Environ()); slog.Check(err) {
						}
						close(s.Ctx.KillAll)
						// time.Sleep(time.Second/2)
						// os.Exit(0)
					}
				}()
			}
		}
	})
}

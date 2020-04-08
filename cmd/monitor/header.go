package monitor

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/ico/svg"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func (s *State) Header(hl bool) layout.FlexChild {
	gtx := s.Gtx
	if hl {
		gtx = s.Htx
	}
	return gui.Rigid(func() {
		s.Rectangle(gtx.Constraints.Width.Max,
			48, "PanelBg", hl)
		s.FlexH(hl,
			gui.Rigid(func() {
				b := s.Buttons["Logo"]
				//s.IconButton(svg.ParallelCoin, "PanelText",
				//	b, hl, func() {
				//	})
				s.ButtonArea(hl, func() {
					s.FlexH(hl,
						gui.Rigid(func() {
							s.Icon(svg.ParallelCoin,
								"PanelText", 32, 8, hl)()
							if b.Clicked(gtx) {
								s.FlipTheme(
									&s.Config.DarkTheme, s.SaveConfig)
							}
						}),
						gui.Rigid(func() {
							s.Text(hl, "Monitor", "PanelText",
								"Secondary", "h3", 48)
						}),
					)
				}, b)
			}),
			s.Spacer(hl),
			s.RestartRunButton(hl),
			gui.Rigid(func() {
				b := s.Buttons["Close"]
				s.IconButton(icons.NavigationClose, "PanelText",
					b, hl, func() {
						Debug("close button clicked")
						s.SaveConfig()
						s.RunCommandChan <- "kill"
						close(s.Ctx.KillAll)
					})
			}),
		)
	})
	//s.ScreenshotButton(hl),
	//)
}

//
//func (s *State) ScreenshotButton(headless bool) layout.FlexChild {
//	gtx := s.Gtx
//	if headless {
//		gtx = s.Htx
//	}
//	return gui.Rigid(func() {
//		b := s.Buttons["Screenshot"]
//		s.IconButton("Screenshot", "PanelText", "PanelBg", b)
//		for b.Clicked(gtx) {
//			Debug("clicked screenshot button")
//			if err := s.Screenshot(func() {
//				s.TopLevelLayout(true)
//			}, *s.Ctx.Config.DataDir+"screenshot.png"); Check(err) {
//			}
//		}
//	})
//}
//
func (s *State) RestartRunButton(hl bool) layout.FlexChild {
	return gui.Rigid(func() {
		var c *exec.Cmd
		var err error
		b := s.Buttons["Restart"]
		s.IconButton(icons.NavigationRefresh, "PanelText", b, hl, func() {
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
					}
				}()
			}
		})
	})
}

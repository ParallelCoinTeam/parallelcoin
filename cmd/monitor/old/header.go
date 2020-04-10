package old

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

func (s *State) ThemeButton() gui.WidgetFunc {
	var err error
	var logoIcon gui.IconFunc
	if logoIcon, err = s.IconSVGtoImage(svg.ParallelCoin, "PanelText",
		32); Check(err) {
	}
	themeWidget := func(hl bool) func() {
		return func() {
			s.FlexH(
				logoIcon(8)(hl),
				s.Text("Monitor", "PanelText", "Secondary", "h3", 48)(hl),
			)(hl)
		}
	}
	return func(hl bool) layout.FlexChild {
		return s.ButtonArea(themeWidget,
			s.FlipTheme(&s.Config.DarkTheme, s.SaveConfig),
			s.Buttons["Logo"])(hl)
	}
}

func (s *State) CloseButton() gui.WidgetFunc {
	var closeIcon gui.IconFunc
	var err error
	if closeIcon, err = s.IconSVGtoImage(icons.NavigationClose, "PanelText",
		32); Check(err) {
	}
	closeButtonWidget := func(hl bool) func() {
		return func() {
			s.FlexH(closeIcon(8)(hl))(hl)
		}
	}
	return func(hl bool) layout.FlexChild {
		return s.ButtonArea(closeButtonWidget, func() {
			Debug("close button clicked")
			s.SaveConfig()
			s.RunCommandChan <- "kill"
			close(s.Ctx.KillAll)
		}, s.Buttons["Close"])(hl)
	}
}

func (s *State) RebuildButton() gui.WidgetFunc {
	var rebuildIcon gui.IconFunc
	var err error
	if rebuildIcon, err = s.IconSVGtoImage(icons.NavigationRefresh, "PanelText",
		32); Check(err) {
	}
	rebuildButtonWidget := func(hl bool) func() {
		return func() {
			s.FlexH(rebuildIcon(8)(hl))(hl)
		}
	}
	return func(hl bool) layout.FlexChild {
		return s.ButtonArea(rebuildButtonWidget, func() {
			var c *exec.Cmd
			var err error
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
		}, s.Buttons["Restart"])(hl)
	}
}

func (s *State) ScreenshotButton() gui.WidgetFunc {
	var screenshotIcon gui.IconFunc
	var err error
	if screenshotIcon, err = s.IconSVGtoImage(icons.ImageCamera, "PanelText",
		32); Check(err) {
	}
	screenShotButtonWidget := func(hl bool) func() {
		return func() {
			s.FlexH(screenshotIcon(8)(hl))(hl)
		}
	}
	return func(hl bool) layout.FlexChild {
		return s.ButtonArea(screenShotButtonWidget, func() {
			Debug("clicked screenshot button")
			if err := s.Screenshot(func() {
				s.TopLevelLayout(true)
			}, *s.Ctx.Config.DataDir+"screenshot.png"); Check(err) {
			}
		}, s.Buttons["Screenshot"])(hl)
	}
}

func (s *State) Header(left, right []gui.WidgetFunc) gui.WidgetFunc {
	return func(hl bool) layout.FlexChild {
		gtx := s.Gtx
		if hl {
			gtx = s.Htx
		}
		var leftW, rightW []layout.FlexChild
		for i := range left {
			leftW = append(leftW, left[i](hl))
		}
		for i := range right {
			rightW = append(rightW, right[i](hl))
		}
		widgets := append(append(leftW, s.Spacer(hl)), rightW...)
		return gui.Rigid(func() {
			s.Rectangle(gtx.Constraints.Width.Max, 48, "PanelBg")(hl)()
			s.FlexH(widgets...)(hl)
		})
	}
}

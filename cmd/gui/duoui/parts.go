package duoui

import "github.com/p9c/pod/cmd/gui/mvc/theme"

func (ui *DuoUI) line(color string) func() {
	return func() {
		cs := ui.ly.Context.Constraints
		theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 1, color, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	}
}

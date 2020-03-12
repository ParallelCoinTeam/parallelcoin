package duoui

import (
	"gioui.org/layout"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/pkg/gelook"
)

var (
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
)

func (ui *DuoUI) DuoUIfooter() func() {
	return func() {
		cs := ui.ly.Context.Constraints
		gelook.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 64, ui.ly.Theme.Colors["Dark"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{Spacing: layout.SpaceBetween}.Layout(ui.ly.Context,
			layout.Rigid(component.FooterLeftMenu(ui.rc, ui.ly.Context, ui.ly.Theme, ui.ly.Pages)),
			layout.Rigid(component.FooterRightMenu(ui.rc, ui.ly.Context, ui.ly.Theme, ui.ly.Pages)),
		)
	}
}

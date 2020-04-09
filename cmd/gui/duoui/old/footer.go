package duoui

import (
	"gioui.org/layout"

	"github.com/p9c/pod/cmd/gui/component"
)

var (
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
)

func (ui *DuoUI) DuoUIfooter() func() {
	ctx := ui.ly.Context
	th := ui.ly.Theme
	return func() {
		footer := th.DuoUIcontainer(0, th.Colors["Dark"])
		footer.FullWidth = true
		footer.Layout(ctx, layout.N, func() {
			layout.Flex{Spacing: layout.SpaceBetween}.Layout(ctx,
				layout.Rigid(component.FooterLeftMenu(ui.rc, ctx, th,
					ui.ly.Pages)),
				layout.Flexed(1, func() {}),
				layout.Rigid(component.FooterRightMenu(ui.rc, ctx, th,
					ui.ly.Pages)),
			)
		})
	}
}

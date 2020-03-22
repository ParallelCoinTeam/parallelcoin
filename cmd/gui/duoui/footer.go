package duoui

import (
	"fmt"
	"gioui.org/layout"

	"github.com/p9c/pod/cmd/gui/component"
)

var (
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
)

func (ui *DuoUI) DuoUIfooter() func() {
	return func() {
		footer := ui.ly.Theme.DuoUIitem(0, ui.ly.Theme.Colors["Dark"])
		footer.FullWidth = true
		footer.Layout(ui.ly.Context, layout.N, func() {
			fmt.Println("footer")
			fmt.Println(ui.ly.Context.Constraints.Width.Max)
			layout.Flex{Spacing: layout.SpaceBetween}.Layout(ui.ly.Context,
				layout.Rigid(component.FooterLeftMenu(ui.rc, ui.ly.Context, ui.ly.Theme, ui.ly.Pages)),
				layout.Flexed(1, func() {}),
				layout.Rigid(component.FooterRightMenu(ui.rc, ui.ly.Context, ui.ly.Theme, ui.ly.Pages)),
			)
		})
	}
}

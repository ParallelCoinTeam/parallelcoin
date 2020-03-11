package duoui

import (
	"gioui.org/layout"
	"github.com/p9c/gel"

	"github.com/p9c/pod/cmd/gui/component"
)

var (
	navButtonOverview    = new(gel.Button)
	navButtonSend        = new(gel.Button)
	navButtonReceive     = new(gel.Button)
	navButtonAddressBook = new(gel.Button)
	navButtonHistory     = new(gel.Button)
	mainNav              = &layout.List{
		Axis: layout.Vertical,
	}

	navItemWidth             int = 96
	navItemHeight            int = 72
	navItemTextSize          int = 48
	navItemTconSize          int = 36
	navItemPaddingVertical   int = 8
	navItemPaddingHorizontal int = 0
)

func (ui *DuoUI) DuoUImenu() func() {
	return func() {
		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEvenly}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				layout.Flex{}.Layout(ui.ly.Context,
					layout.Rigid(component.MainNavigation(ui.rc, ui.ly.Context, ui.ly.Theme, ui.ly.Pages)),
				)
			}),
		)
	}
}

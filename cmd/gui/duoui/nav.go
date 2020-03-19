package duoui

import (
	"gioui.org/layout"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/pkg/gel"
)

var (
	navButtonOverview = new(gel.Button)
	navButtonSend     = new(gel.Button)
	// navButtonReceive     = new(gel.Button)
	navButtonAddressBook = new(gel.Button)
	navButtonHistory     = new(gel.Button)
	mainNav              = &layout.List{
		Axis: layout.Vertical,
	}

	navItemWidth             = 96
	navItemHeight            = 72
	navItemTextSize          = 48
	navItemTconSize          = 36
	navItemPaddingVertical   = 8
	navItemPaddingHorizontal = 0
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

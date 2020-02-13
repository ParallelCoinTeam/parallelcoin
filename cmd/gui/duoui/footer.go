package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	buttonTrace     = new(widget.Button)
	buttonCornerOne = new(widget.Button)
	buttonSettings  = new(widget.Button)
	buttonNetwork   = new(widget.Button)
	buttonBlocks    = new(widget.Button)
	buttonConsole   = new(widget.Button)
	buttonHelp      = new(widget.Button)
	cornerNav       = &layout.List{
		Axis: layout.Horizontal,
	}
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
)

func (duo *DuoUI) DuoUIfooter(rc *rcd.RcVar) func() {
	return func() {
		cs := duo.Model.DuoUIcontext.Constraints
		helpers.DuoUIdrawRectangle(duo.Model.DuoUIcontext, cs.Width.Max, 64, duo.Model.DuoUItheme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		var (
			width             float32 = 48
			height            float32 = 48
			iconSize          int     = 32
			paddingVertical   float32 = 8
			paddingHorizontal float32 = 8
		)
		settingsIcon, _ := parallel.NewDuoUIicon(icons.ActionSettings)
		blocksIcon, _ := parallel.NewDuoUIicon(icons.ActionExplore)
		networkIcon, _ := parallel.NewDuoUIicon(icons.ActionFingerprint)
		traceIcon, _ := parallel.NewDuoUIicon(icons.ActionTrackChanges)
		consoleIcon, _ := parallel.NewDuoUIicon(icons.ActionInput)
		helpIcon, _ := parallel.NewDuoUIicon(icons.NavigationArrowDropDown)
		layout.Flex{Spacing: layout.SpaceBetween}.Layout(duo.Model.DuoUIcontext,
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(0)).Layout(duo.Model.DuoUIcontext, func() {
					cornerButtons := []func(){
						func() {
							layout.UniformInset(unit.Dp(0)).Layout(duo.Model.DuoUIcontext, func() {
								var networkMeniItem parallel.DuoUIbutton
								networkMeniItem = duo.Model.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, networkIcon)
								for buttonCornerOne.Clicked(duo.Model.DuoUIcontext) {
									duo.Model.CurrentPage = "NETWORK"
								}
								networkMeniItem.Layout(duo.Model.DuoUIcontext, buttonCornerOne)
							})
						},
						func() {
							var settingsMenuItem parallel.DuoUIbutton
							settingsMenuItem = duo.Model.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, traceIcon)

							for buttonTrace.Clicked(duo.Model.DuoUIcontext) {
								duo.Model.CurrentPage = "TRACE"
							}
							settingsMenuItem.Layout(duo.Model.DuoUIcontext, buttonTrace)
						},
					}
					cornerNav.Layout(duo.Model.DuoUIcontext, len(cornerButtons), func(i int) {
						layout.UniformInset(unit.Dp(0)).Layout(duo.Model.DuoUIcontext, cornerButtons[i])
					})
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(0)).Layout(duo.Model.DuoUIcontext, func() {
					navButtons := []func(){
						func() {
							layout.UniformInset(unit.Dp(0)).Layout(duo.Model.DuoUIcontext, func() {
								var networkMeniItem parallel.DuoUIbutton
								networkMeniItem = duo.Model.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, networkIcon)
								for buttonNetwork.Clicked(duo.Model.DuoUIcontext) {
									duo.Model.CurrentPage = "NETWORK"
								}
								networkMeniItem.Layout(duo.Model.DuoUIcontext, buttonNetwork)
							})
						},
						func() {
							var blocksMenuItem parallel.DuoUIbutton
							blocksMenuItem = duo.Model.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, blocksIcon)
							for buttonBlocks.Clicked(duo.Model.DuoUIcontext) {
								//duo.CurrentPage = "EXPLORER"
								//rc.ShowToast = true
								toastAdd(duo, rc)
							}
							blocksMenuItem.Layout(duo.Model.DuoUIcontext, buttonBlocks)
						},
						func() {
							var helpMenuItem parallel.DuoUIbutton
							helpMenuItem = duo.Model.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, helpIcon)
							for buttonHelp.Clicked(duo.Model.DuoUIcontext) {
								rc.ShowDialog = true
							}
							helpMenuItem.Layout(duo.Model.DuoUIcontext, buttonHelp)
						},
						func() {
							layout.UniformInset(unit.Dp(0)).Layout(duo.Model.DuoUIcontext, func() {
								var consoleMenuItem parallel.DuoUIbutton
								consoleMenuItem = duo.Model.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, consoleIcon)
								for buttonConsole.Clicked(duo.Model.DuoUIcontext) {
									duo.Model.CurrentPage = "CONSOLE"
								}
								consoleMenuItem.Layout(duo.Model.DuoUIcontext, buttonConsole)
							})
						},
						func() {
							var settingsMenuItem parallel.DuoUIbutton
							settingsMenuItem = duo.Model.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, settingsIcon)

							for buttonSettings.Clicked(duo.Model.DuoUIcontext) {
								duo.Model.CurrentPage = "SETTINGS"
							}
							settingsMenuItem.Layout(duo.Model.DuoUIcontext, buttonSettings)
						},
					}
					footerNav.Layout(duo.Model.DuoUIcontext, len(navButtons), func(i int) {
						layout.UniformInset(unit.Dp(0)).Layout(duo.Model.DuoUIcontext, navButtons[i])
					})
				})
			}),
		)
	}
}

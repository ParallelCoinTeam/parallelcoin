package duoui

import (
	"github.com/p9c/pod/cmd/gui/components"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gio/widget"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	buttonCornerOne = new(widget.Button)
	buttonCornerTwo = new(widget.Button)
	buttonSettings  = new(widget.Button)
	buttonNetwork   = new(widget.Button)
	buttonBlocks    = new(widget.Button)
	buttonConsole   = new(widget.Button)
	buttonHelp      = new(widget.Button)
	cornerNav   = &layout.List{
		Axis: layout.Horizontal,
	}
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
)

func DuoUIfooter(duo *models.DuoUI, rc *rcd.RcVar) {
	cs := duo.DuoUIcontext.Constraints
	helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 64, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

	// Footer <<<
	var (
		width             float32 = 48
		height            float32 = 48
		iconSize          int     = 32
		paddingVertical   float32 = 8
		paddingHorizontal float32 = 8
	)
	settingsIcon, _ := components.NewDuoUIicon(icons.ActionSettings)
	blocksIcon, _ := components.NewDuoUIicon(icons.ActionExplore)
	networkIcon, _ := components.NewDuoUIicon(icons.ActionFingerprint)
	consoleIcon, _ := components.NewDuoUIicon(icons.ActionInput)
	helpIcon, _ := components.NewDuoUIicon(icons.NavigationArrowDropDown)

	//overviewButton :=
	//historyButton :=

	in := layout.UniformInset(unit.Dp(0))

	layout.Flex{Spacing: layout.SpaceBetween}.Layout(duo.DuoUIcontext,
		layout.Rigid(func() {
			in.Layout(duo.DuoUIcontext, func() {
				cornerButtons := []func(){
					func() {
						in.Layout(duo.DuoUIcontext, func() {
							var networkMeniItem components.DuoUIbutton
							networkMeniItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, networkIcon)
							for buttonCornerOne.Clicked(duo.DuoUIcontext) {
								duo.CurrentPage = "Network"
							}
							networkMeniItem.Layout(duo.DuoUIcontext, buttonCornerOne)
						})
					},
					func() {
						var settingsMenuItem components.DuoUIbutton
						settingsMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, settingsIcon)

						for buttonCornerTwo.Clicked(duo.DuoUIcontext) {
							duo.CurrentPage = "Settings"
						}
						settingsMenuItem.Layout(duo.DuoUIcontext, buttonCornerTwo)
					},
				}
				cornerNav.Layout(duo.DuoUIcontext, len(cornerButtons), func(i int) {
					layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, cornerButtons[i])
				})
			})
		}),
		layout.Rigid(func() {
			in.Layout(duo.DuoUIcontext, func() {
				navButtons := []func(){
					func() {
						in.Layout(duo.DuoUIcontext, func() {
							var networkMeniItem components.DuoUIbutton
							networkMeniItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, networkIcon)
							for buttonNetwork.Clicked(duo.DuoUIcontext) {
								duo.CurrentPage = "Network"
							}
							networkMeniItem.Layout(duo.DuoUIcontext, buttonNetwork)
						})
					},
					func() {
						var blocksMenuItem components.DuoUIbutton
						blocksMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, blocksIcon)
						for buttonBlocks.Clicked(duo.DuoUIcontext) {
							duo.CurrentPage = "Explorer"
						}
						blocksMenuItem.Layout(duo.DuoUIcontext, buttonBlocks)
					},
					func() {
						var helpMenuItem components.DuoUIbutton
						helpMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, helpIcon)
						for buttonHelp.Clicked(duo.DuoUIcontext) {
							duo.CurrentPage = "Help"
						}
						helpMenuItem.Layout(duo.DuoUIcontext, buttonHelp)
					},
					func() {
						in.Layout(duo.DuoUIcontext, func() {
							var consoleMenuItem components.DuoUIbutton
							consoleMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, consoleIcon)
							for buttonConsole.Clicked(duo.DuoUIcontext) {
								duo.CurrentPage = "Console"
							}
							consoleMenuItem.Layout(duo.DuoUIcontext, buttonConsole)
						})
					},
					func() {
						var settingsMenuItem components.DuoUIbutton
						settingsMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, settingsIcon)

						for buttonSettings.Clicked(duo.DuoUIcontext) {
							duo.CurrentPage = "Settings"
						}
						settingsMenuItem.Layout(duo.DuoUIcontext, buttonSettings)
					},
				}
				footerNav.Layout(duo.DuoUIcontext, len(navButtons), func(i int) {
					layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, navButtons[i])
				})
			})
		}),
	)
	//}),
	//)
	// Footer >>>

}

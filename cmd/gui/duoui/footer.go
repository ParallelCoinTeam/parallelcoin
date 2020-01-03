package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/theme"
	"github.com/p9c/pod/cmd/gui/widget"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	buttonSettings = new(widget.Button)
	buttonNetwork  = new(widget.Button)
	buttonBlocks   = new(widget.Button)
	buttonConsole  = new(widget.Button)
	buttonHelp     = new(widget.Button)
	footerNav      = &layout.List{
		Axis: layout.Horizontal,
	}
)

func DuoUIfooter(duo *models.DuoUI, rc *rcd.RcVar) {
	// Footer <<<
	var (
		width             float32 = 64
		height            float32 = 64
		iconSize          int     = 32
		paddingVertical   float32 = 8
		paddingHorizontal float32 = 8
	)
	settingsIcon, _ := theme.NewDuoUIicon(icons.ActionSettings)
	blocksIcon, _ := theme.NewDuoUIicon(icons.ActionExplore)
	networkIcon, _ := theme.NewDuoUIicon(icons.ActionFingerprint)
	consoleIcon, _ := theme.NewDuoUIicon(icons.ActionInput)
	helpIcon, _ := theme.NewDuoUIicon(icons.NavigationArrowDropDown)

	//overviewButton :=
	//historyButton :=

	in := layout.UniformInset(unit.Dp(0))

	duo.DuoUIcomponents.Footer.Layout.Layout(duo.DuoUIcontext,
		layout.Rigid(func() {
			layout.Flex{
				Alignment: layout.End,
			}.Layout(duo.DuoUIcontext,
				layout.Rigid(func() {
					in.Layout(duo.DuoUIcontext, func() {

						navButtons := []func(){
							func() {
								in.Layout(duo.DuoUIcontext, func() {
									var networkMeniItem theme.DuoUIbutton
									networkMeniItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, networkIcon)
									for buttonNetwork.Clicked(duo.DuoUIcontext) {
										duo.CurrentPage = "Network"
									}
									networkMeniItem.Layout(duo.DuoUIcontext, buttonNetwork)
								})
							},
							func() {
								var blocksMenuItem theme.DuoUIbutton
								blocksMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, blocksIcon)
								for buttonBlocks.Clicked(duo.DuoUIcontext) {
									duo.CurrentPage = "Blocks"
								}
								blocksMenuItem.Layout(duo.DuoUIcontext, buttonBlocks)
							},
							func() {
								var helpMenuItem theme.DuoUIbutton
								helpMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, helpIcon)
								for buttonHelp.Clicked(duo.DuoUIcontext) {
									duo.CurrentPage = "Help"
								}
								helpMenuItem.Layout(duo.DuoUIcontext, buttonHelp)
							},
							func() {
								in.Layout(duo.DuoUIcontext, func() {
									var consoleMenuItem theme.DuoUIbutton
									consoleMenuItem = duo.DuoUItheme.DuoUIbutton("", "",  "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, consoleIcon)
									for buttonConsole.Clicked(duo.DuoUIcontext) {
										duo.CurrentPage = "Console"
									}
									consoleMenuItem.Layout(duo.DuoUIcontext, buttonConsole)
								})
							},
							func() {
								var settingsMenuItem theme.DuoUIbutton
								settingsMenuItem = duo.DuoUItheme.DuoUIbutton("", "",  "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, settingsIcon)

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
		}),
	)
	// Footer >>>

}

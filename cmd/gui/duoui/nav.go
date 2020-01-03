package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/theme"
	"github.com/p9c/pod/cmd/gui/widget"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	buttonOverview    = new(widget.Button)
	buttonSend        = new(widget.Button)
	buttonReceive     = new(widget.Button)
	buttonAddressBook = new(widget.Button)
	buttonHistory     = new(widget.Button)
	mainNav           = &layout.List{
		Axis: layout.Vertical,
	}
)

func DuoUImenu(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	overviewIcon, _ := theme.NewDuoUIicon(icons.ActionHome)
	sendIcon, _ := theme.NewDuoUIicon(icons.NavigationArrowDropUp)
	receiveIcon, _ := theme.NewDuoUIicon(icons.NavigationArrowDropDown)
	addressBookIcon, _ := theme.NewDuoUIicon(icons.ActionBook)
	historyIcon, _ := theme.NewDuoUIicon(icons.ActionHistory)

	var (
		width             float32 = 96
		height            float32 = 96
		iconSize          int     = 48
		paddingVertical   float32 = 4
		paddingHorizontal float32 = 8
	)
	//overviewButton :=
	//historyButton :=

	in := layout.UniformInset(unit.Dp(0))

	duo.DuoUIcomponents.Menu.Layout.Layout(duo.DuoUIcontext,
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.DuoUIcontext,
				layout.Rigid(func() {
					in.Layout(duo.DuoUIcontext, func() {

						navButtons := []func(){
							func() {
								in.Layout(duo.DuoUIcontext, func() {
									var overviewMenuItem theme.DuoUIbutton
									overviewMenuItem = duo.DuoUItheme.DuoUIbutton("Overview", "ff303030",  "ff989898", "ff303030", iconSize, width, height, paddingVertical, paddingHorizontal, overviewIcon)
									for buttonOverview.Clicked(duo.DuoUIcontext) {
										duo.CurrentPage = "Overview"
									}
									overviewMenuItem.Layout(duo.DuoUIcontext, buttonOverview)
								})
							},
							func() {
								helpers.DuoUIdrawRectangle(duo.DuoUIcontext, int(width), 1, helpers.HexARGB("ff888888"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
							},
							func() {
								var sendMenuItem theme.DuoUIbutton
								sendMenuItem = duo.DuoUItheme.DuoUIbutton("Send", "ff303030",  "ff989898", "ff303030", iconSize, width, height, paddingVertical, paddingHorizontal, sendIcon)
								for buttonSend.Clicked(duo.DuoUIcontext) {
									duo.CurrentPage = "Send"
								}
								sendMenuItem.Layout(duo.DuoUIcontext, buttonSend)
							},
							func() {
								helpers.DuoUIdrawRectangle(duo.DuoUIcontext, int(width), 1, helpers.HexARGB("ff888888"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
							},
							func() {
								in.Layout(duo.DuoUIcontext, func() {
									var receiveMenuItem theme.DuoUIbutton
									receiveMenuItem = duo.DuoUItheme.DuoUIbutton("Receive", "ff303030",  "ff989898", "ff303030", iconSize, width, height, paddingVertical, paddingHorizontal, receiveIcon)
									for buttonReceive.Clicked(duo.DuoUIcontext) {
										duo.CurrentPage = "Receive"
									}
									receiveMenuItem.Layout(duo.DuoUIcontext, buttonReceive)
								})
							},
							func() {
								helpers.DuoUIdrawRectangle(duo.DuoUIcontext, int(width), 1, helpers.HexARGB("ff888888"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
							},
							func() {
								var addressBookMenuItem theme.DuoUIbutton
								addressBookMenuItem = duo.DuoUItheme.DuoUIbutton("Address Book", "ff303030",  "ff989898", "ff303030", iconSize, width, height, paddingVertical, paddingHorizontal, addressBookIcon)
								for buttonHistory.Clicked(duo.DuoUIcontext) {
									duo.CurrentPage = "AddressBook"
								}
								addressBookMenuItem.Layout(duo.DuoUIcontext, buttonAddressBook)
							},
							func() {
								helpers.DuoUIdrawRectangle(duo.DuoUIcontext, int(width), 1, helpers.HexARGB("ff888888"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
							},
							func() {
								var historyMenuItem theme.DuoUIbutton
								historyMenuItem = duo.DuoUItheme.DuoUIbutton("History", "ff303030",  "ff989898", "ff303030", iconSize, width, height, paddingVertical, paddingHorizontal, historyIcon)
								for buttonHistory.Clicked(duo.DuoUIcontext) {
									duo.CurrentPage = "History"
								}
								historyMenuItem.Layout(duo.DuoUIcontext, buttonHistory)
							},
						}
						mainNav.Layout(duo.DuoUIcontext, len(navButtons), func(i int) {
							layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, navButtons[i])
						})
					})
				}),
			)
		}),
	)
}

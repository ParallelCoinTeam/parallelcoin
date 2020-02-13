package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
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

func DuoUImenu(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) func(){
	return func() {
		overviewIcon, _ := parallel.NewDuoUIicon(icons.ActionHome)
		sendIcon, _ := parallel.NewDuoUIicon(icons.NavigationArrowDropUp)
		receiveIcon, _ := parallel.NewDuoUIicon(icons.NavigationArrowDropDown)
		addressBookIcon, _ := parallel.NewDuoUIicon(icons.ActionBook)
		historyIcon, _ := parallel.NewDuoUIicon(icons.ActionHistory)

		var (
			width             float32 = 96
			height            float32 = 72
			iconSize          int     = 48
			paddingVertical   float32 = 2
			paddingHorizontal float32 = 8
		)
		//overviewButton :=
		//historyButton :=

		in := layout.UniformInset(unit.Dp(0))

		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEvenly}.Layout(duo.Model.DuoUIcontext,
			layout.Rigid(func() {
				layout.Flex{}.Layout(duo.Model.DuoUIcontext,
					layout.Rigid(func() {
						in.Layout(duo.Model.DuoUIcontext, func() {

							navButtons := []func(){
								func() {
									in.Layout(duo.Model.DuoUIcontext, func() {
										var overviewMenuItem parallel.DuoUIbutton
										overviewMenuItem = duo.Model.DuoUItheme.DuoUIbutton("OVERVIEW", "ff303030", "ff989898", "ff303030", iconSize, width, height, paddingVertical, paddingHorizontal, overviewIcon)
										for buttonOverview.Clicked(duo.Model.DuoUIcontext) {
											duo.Model.CurrentPage = "OVERVIEW"
										}
										overviewMenuItem.Layout(duo.Model.DuoUIcontext, buttonOverview)
									})
								},
								func() {
									helpers.DuoUIdrawRectangle(duo.Model.DuoUIcontext, int(width), 1, helpers.HexARGB("ff888888"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									var sendMenuItem parallel.DuoUIbutton
									sendMenuItem = duo.Model.DuoUItheme.DuoUIbutton("SEND", "ff303030", "ff989898", "ff303030", iconSize, width, height, paddingVertical, paddingHorizontal, sendIcon)
									for buttonSend.Clicked(duo.Model.DuoUIcontext) {
										duo.Model.CurrentPage = "SEND"
									}
									sendMenuItem.Layout(duo.Model.DuoUIcontext, buttonSend)
								},
								func() {
									helpers.DuoUIdrawRectangle(duo.Model.DuoUIcontext, int(width), 1, helpers.HexARGB("ff888888"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									in.Layout(duo.Model.DuoUIcontext, func() {
										var receiveMenuItem parallel.DuoUIbutton
										receiveMenuItem = duo.Model.DuoUItheme.DuoUIbutton("RECEIVE", "ff303030", "ff989898", "ff303030", iconSize, width, height, paddingVertical, paddingHorizontal, receiveIcon)
										for buttonReceive.Clicked(duo.Model.DuoUIcontext) {
											duo.Model.CurrentPage = "RECEIVE"
										}
										receiveMenuItem.Layout(duo.Model.DuoUIcontext, buttonReceive)
									})
								},
								func() {
									helpers.DuoUIdrawRectangle(duo.Model.DuoUIcontext, int(width), 1, helpers.HexARGB("ff888888"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									var addressBookMenuItem parallel.DuoUIbutton
									addressBookMenuItem = duo.Model.DuoUItheme.DuoUIbutton("ADDRESS BOOK", "ff303030", "ff989898", "ff303030", iconSize, width, height, paddingVertical, paddingHorizontal, addressBookIcon)
									for buttonAddressBook.Clicked(duo.Model.DuoUIcontext) {
										duo.Model.CurrentPage = "ADDRESSBOOK"
									}
									addressBookMenuItem.Layout(duo.Model.DuoUIcontext, buttonAddressBook)
								},
								func() {
									helpers.DuoUIdrawRectangle(duo.Model.DuoUIcontext, int(width), 1, helpers.HexARGB("ff888888"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									var historyMenuItem parallel.DuoUIbutton
									historyMenuItem = duo.Model.DuoUItheme.DuoUIbutton("HISTORY", "ff303030", "ff989898", "ff303030", iconSize, width, height, paddingVertical, paddingHorizontal, historyIcon)
									for buttonHistory.Clicked(duo.Model.DuoUIcontext) {
										duo.Model.CurrentPage = "HISTORY"
									}
									historyMenuItem.Layout(duo.Model.DuoUIcontext, buttonHistory)
								},
							}
							mainNav.Layout(duo.Model.DuoUIcontext, len(navButtons), func(i int) {
								layout.UniformInset(unit.Dp(0)).Layout(duo.Model.DuoUIcontext, navButtons[i])
							})
						})
					}),
				)
			}),
		)
	}
}
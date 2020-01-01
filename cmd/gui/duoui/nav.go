package duoui

import (
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
	buttonOverview = new(widget.Button)
	buttonHistory  = new(widget.Button)
	navList        = &layout.List{
		Axis: layout.Vertical,
	}
)

func DuoUImenu(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	overviewIcon, _ := theme.NewDuoUIicon(icons.ActionHome)
	historyIcon, _ := theme.NewDuoUIicon(icons.ActionHistory)

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
									var overviewMeniItem  theme.DuoUIbutton

									overviewMeniItem = duo.DuoUItheme.DuoUIbutton("Overview", "","", 64, 48, 48, 0, overviewIcon)
									for buttonOverview.Clicked(duo.DuoUIcontext) {
										duo.CurrentPage = "Overview"
									}
									overviewMeniItem.Layout(duo.DuoUIcontext, buttonOverview)
								})
							},
							func() {
								var historyMenuItem theme.DuoUIbutton
								historyMenuItem = duo.DuoUItheme.DuoUIbutton("History", "","", 64, 48, 48, 0, historyIcon)

								for buttonHistory.Clicked(duo.DuoUIcontext) {
									duo.CurrentPage = "History"
								}
								historyMenuItem.Layout(duo.DuoUIcontext, buttonHistory)
							},
						}
						navList.Layout(duo.DuoUIcontext, len(navButtons), func(i int) {
							layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, navButtons[i])
						})
					})
				}),
			)
		}),


	)
}

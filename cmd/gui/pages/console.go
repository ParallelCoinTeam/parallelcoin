package pages

import (
	"time"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

var (
	consoleInputField = &gel.Editor{
		SingleLine: true,
		Submit:     true,
	}
	consoleOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
)

func Console(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "CONSOLE",
		Command:       func() {},
		Border:        4,
		BorderColor:   "ff000000",
		Header:        func() {},
		HeaderPadding: 0,
		Body:          consoleBody(rc, gtx, th),
		BodyBgColor:   "ff000000",
		BodyPadding:   4,
		Footer:        func() {},
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func consoleBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
			th.DuoUIitem(0, "ff000000").Layout(gtx, layout.N, func() {
				layout.Flex{}.Layout(gtx, layout.Flexed(1, func() {
					layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
						layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceAround}.Layout(gtx,
							layout.Flexed(1, func() {
								consoleOutputList.Layout(gtx,
									len(rc.ConsoleHistory.Commands), func(i int) {
										t := rc.ConsoleHistory.Commands[i]
										layout.Flex{
											Axis:      layout.Vertical,
											Alignment: layout.End,
										}.Layout(gtx,
											layout.Rigid(component.Label(
												gtx, th, th.Fonts["Mono"],
												12, th.Colors["Light"], "ds://"+t.ComID)),
											layout.Rigid(component.Label(
												gtx, th, th.Fonts["Mono"],
												12, th.Colors["Light"], t.Out)),
										)
									})
							}),
							layout.Rigid(
								component.ConsoleInput(gtx, th,
									consoleInputField, "Run command",
									func(e gel.SubmitEvent) {
										rc.ConsoleHistory.Commands = append(
											rc.ConsoleHistory.Commands,
											model.DuoUIconsoleCommand{
												ComID: e.Text,
												Time:  time.Time{},
												Out:   rc.ConsoleCmd(e.Text),
											})
									}),
							),
						)
					})
				}),
				)
			})
		})
	}
}

package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	"time"
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
	return th.DuoUIpage("CONSOLE", 0, func() {}, func() {}, consoleBody(rc, gtx, th), func() {})
}
func consoleBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Flexed(1, func() {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
					layout.Flex{
						Axis:    layout.Vertical,
						Spacing: layout.SpaceAround,
					}.Layout(gtx,
						layout.Flexed(1, func() {
							consoleOutputList.Layout(gtx, len(rc.ConsoleHistory.Commands), func(i int) {
								t := rc.ConsoleHistory.Commands[i]
								layout.Flex{
									Axis:      layout.Vertical,
									Alignment: layout.End,
								}.Layout(gtx,
									layout.Rigid(component.Label(gtx, th, th.Fonts["Mono"], 12, th.Colors["Dark"], "ds://"+t.ComID)),
									layout.Rigid(component.Label(gtx, th, th.Fonts["Mono"], 12, th.Colors["Dark"], t.Out)),
								)
							})
						}),
						layout.Rigid(
							component.ConsoleInput(gtx, th, consoleInputField, "Run command", func(e gel.SubmitEvent) {
								rc.ConsoleHistory.Commands = append(rc.ConsoleHistory.Commands, model.DuoUIconsoleCommand{
									ComID: e.Text,
									Time:  time.Time{},
									Out:   rc.ConsoleCmd(e.Text),
								})
							})))
				})
			}),
		)
	}
}

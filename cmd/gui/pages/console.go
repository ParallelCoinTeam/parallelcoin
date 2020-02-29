package pages

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
	"time"
)

var (
	testLabel         = "testtopLabel"
	consoleInputField = &controller.Editor{
		SingleLine: true,
		Submit:     true,
	}
	consoleOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
)

func Console(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("CONSOLE", 0, func() {}, func() {}, console(rc, gtx, th), func() {})
}
func console(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Flexed(1, func() {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
					layout.Flex{
						Axis:    layout.Vertical,
						Spacing: layout.SpaceAround,
					}.Layout(gtx,
						layout.Flexed(1, func() {
							consoleOutputList.Layout(gtx, len(rc.CommandsHistory.Commands), func(i int) {
								t := rc.CommandsHistory.Commands[i]
								layout.Flex{
									Axis:      layout.Vertical,
									Alignment: layout.End,
								}.Layout(gtx,
									layout.Rigid(func() {
										sat := th.Body1("ds://" + t.ComID)
										sat.Font.Typeface = th.Font.Mono
										sat.Color = theme.HexARGB(th.Color.Dark)
										sat.Layout(gtx)
									}),
									layout.Rigid(func() {
										sat := th.Body1(t.Out)
										sat.Font.Typeface = th.Font.Mono
										sat.Color = theme.HexARGB(th.Color.Dark)
										sat.Layout(gtx)
									}),
								)
							})
						}),
						layout.Rigid(func() {
							layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
								e := th.DuoUIeditor("Run command")
								e.Font.Typeface = th.Font.Mono
								e.Color = theme.HexARGB(th.Color.Dark)
								e.Font.Style = text.Regular
								e.Layout(gtx, consoleInputField)
								for _, e := range consoleInputField.Events(gtx) {
									if e, ok := e.(controller.SubmitEvent); ok {
										rc.CommandsHistory.Commands = append(rc.CommandsHistory.Commands, model.DuoUIcommand{
											ComID: e.Text,
											Time:  time.Time{},
											Out:   rc.ConsoleCmd(e.Text),
										})
										consoleInputField.SetText("")
									}
								}
							})
						}))
				})
			}),
		)
	}
}

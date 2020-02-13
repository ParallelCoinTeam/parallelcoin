package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
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

func (ui *DuoUI) DuoUIconsole() {
	layout.Flex{}.Layout(ui.ly.Context,
		layout.Flexed(1, func() {
			layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceAround,
				}.Layout(ui.ly.Context,
					layout.Flexed(1, func() {
						consoleOutputList.Layout(ui.ly.Context, len(ui.rc.CommandsHistory.Commands), func(i int) {
							t := ui.rc.CommandsHistory.Commands[i]
							layout.Flex{
								Alignment: layout.End,
							}.Layout(ui.ly.Context,
								layout.Rigid(func() {
									sat := ui.ly.Theme.Body1("ds://" + t.ComID)
									sat.Font.Size = unit.Dp(16)
									sat.Layout(ui.ly.Context)
								}),
							)
						})
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, func() {
							e := ui.ly.Theme.DuoUIeditor("Run command", "Run txt")
							e.Font.Style = text.Regular
							e.Font.Size = unit.Dp(16)
							e.Layout(ui.ly.Context, consoleInputField)
							for _, e := range consoleInputField.Events(ui.ly.Context) {
								if e, ok := e.(controller.SubmitEvent); ok {
									ui.rc.CommandsHistory.Commands = append(ui.rc.CommandsHistory.Commands, controller.Command{
										ComID: e.Text,
										Time:  time.Time{},
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

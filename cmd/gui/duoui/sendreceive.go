package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/clipboard"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"strconv"
)

var (
	address           string
	amount            float64
	passPharse        string
	addressLineEditor = &controller.Editor{
		SingleLine: true,
		Submit:     true,
	}
	amountLineEditor = &controller.Editor{
		SingleLine: true,
		Submit:     true,
	}
	buttonPasteAddress = new(controller.Button)
	buttonPasteAmount  = new(controller.Button)
	buttonSend         = new(controller.Button)
)

func (ui *DuoUI) DuoUIsend() func() {
	return func() {
		layout.Flex{}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				cs := ui.ly.Context.Constraints
				theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 180, ui.ly.Theme.Color.Bg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(ui.ly.Context,
					layout.Rigid(func() {

						layout.Flex{}.Layout(ui.ly.Context,
							layout.Flexed(1, func() {

								layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
									cs := ui.ly.Context.Constraints
									theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 32, "fff4f4f4", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
									layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
										theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 30, "ffffffff", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
										e := ui.ly.Theme.DuoUIeditor("DUO address", "DUO dva")
										e.Font.Typeface = ui.ly.Theme.Font.Primary
										e.Font.Style = text.Italic
										e.Font.Size = unit.Dp(24)
										e.Layout(ui.ly.Context, addressLineEditor)
										for _, e := range addressLineEditor.Events(ui.ly.Context) {
											if e, ok := e.(controller.SubmitEvent); ok {
												address = e.Text
												addressLineEditor.SetText("")
											}
										}
									})
								})

							}),
							layout.Rigid(func() {
								layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
									var pasteAddressButton theme.DuoUIbutton
									pasteAddressButton = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary, "PASTE ADDRESS", "ffcfcfcf", "ff303030", "", "ffcfcfcf", 0, 128, 48, 0, 0)

									for buttonPasteAddress.Clicked(ui.ly.Context) {
										addressLineEditor.SetText(clipboard.Get())
									}
									pasteAddressButton.Layout(ui.ly.Context, buttonPasteAddress)
								})

							}))

					}),
					layout.Rigid(func() {
						layout.Flex{}.Layout(ui.ly.Context,
							layout.Flexed(1, func() {
								layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
									cs := ui.ly.Context.Constraints
									theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 32, "fff4f4f4", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
									layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
										theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 30, "ffffffff", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
										e := ui.ly.Theme.DuoUIeditor("DUO Amount", "DUO dva")
										e.Font.Typeface = ui.ly.Theme.Font.Primary
										e.Font.Style = text.Italic
										e.Font.Size = unit.Dp(24)
										e.Layout(ui.ly.Context, amountLineEditor)
										for _, e := range amountLineEditor.Events(ui.ly.Context) {
											if e, ok := e.(controller.SubmitEvent); ok {
												f, err := strconv.ParseFloat(e.Text, 64)
												if err != nil {
													amount = f
													amountLineEditor.SetText("")
												}
											}
										}
									})
								})
							}),
							layout.Rigid(func() {
								layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
									var pasteAmountButton theme.DuoUIbutton
									pasteAmountButton = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary, "PASTE AMOUNT", "ffcfcfcf", "ff303030", "", "ffcfcfcf", 0, 128, 48, 0, 0)

									for buttonPasteAmount.Clicked(ui.ly.Context) {
										addressLineEditor.SetText(clipboard.Get())
									}
									pasteAmountButton.Layout(ui.ly.Context, buttonPasteAmount)
								})

							}))
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
							var sendButton theme.DuoUIbutton
							sendButton = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary, "SEND", "ffcfcfcf", "ff303030", "", "ffcfcfcf", 0, 128, 48, 0, 0)

							for buttonSend.Clicked(ui.ly.Context) {

							}
							sendButton.Layout(ui.ly.Context, buttonSend)
						})
					}))
			}),
		)
	}
}

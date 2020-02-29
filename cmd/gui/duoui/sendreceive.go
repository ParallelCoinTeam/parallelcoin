package duoui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/clipboard"
	"strconv"
)

var (
	layautList = &layout.List{
		Axis: layout.Vertical,
	}
	address           string
	amount            float64
	passPharse        string
	addressLineEditor = &controller.Editor{
		SingleLine: true,
	}
	amountLineEditor = &controller.Editor{
		SingleLine: true,
	}
	passLineEditor = &controller.Editor{
		SingleLine: true,
	}
	buttonPasteAddress = new(controller.Button)
	buttonPasteAmount  = new(controller.Button)
	buttonSend         = new(controller.Button)
)

func send(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(func() {
				cs := gtx.Constraints
				theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 180, th.Color.Light, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

				widgets := []func(){
					func() {

						layout.Flex{}.Layout(gtx,
							layout.Flexed(1, addressEditor(gtx, th)),
							layout.Rigid(pasteAddressButton(gtx, th)))

					},
					func() {
						layout.Flex{}.Layout(gtx,
							layout.Flexed(1, amountEditor(gtx, th)),
							layout.Rigid(pasteAmountButton(gtx, th)))

					},
					func() {
						layout.Flex{}.Layout(gtx,
							layout.Flexed(1, passwordEditor(gtx, th)))
					},
					func() {
						layout.Rigid(sendButton(rc, gtx, th))
					},
				}
				layautList.Layout(gtx, len(widgets), func(i int) {
					layout.UniformInset(unit.Dp(8)).Layout(gtx, widgets[i])
				})
			}))
	}
}

func pasteAddressButton(gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var pasteAddressButton theme.DuoUIbutton
			pasteAddressButton = th.DuoUIbutton(th.Font.Secondary, "PASTE ADDRESS", th.Color.Light, th.Color.Dark, "", th.Color.Light, 16, 0, 128, 48, 0, 0)
			for buttonPasteAddress.Clicked(gtx) {
				addressLineEditor.SetText(clipboard.Get())
			}
			pasteAddressButton.Layout(gtx, buttonPasteAddress)
		})
	}
}

func addressEditor(gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			cs := gtx.Constraints
			theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 32, "fff4f4f4", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 30, "ffffffff", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				e := th.DuoUIeditor("DUO address")
				e.Font.Typeface = th.Font.Primary
				e.Font.Style = text.Italic
				e.Layout(gtx, addressLineEditor)
				for _, e := range addressLineEditor.Events(gtx) {
					if e, ok := e.(controller.SubmitEvent); ok {
						address = e.Text
						addressLineEditor.SetText("")
					}
				}
			})
		})
	}
}

func pasteAmountButton(gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var pasteAmountButton theme.DuoUIbutton
			pasteAmountButton = th.DuoUIbutton(th.Font.Secondary, "PASTE AMOUNT", th.Color.Light, th.Color.Dark, "", th.Color.Light, 16, 0, 128, 48, 0, 0)

			for buttonPasteAmount.Clicked(gtx) {
				amountLineEditor.SetText(clipboard.Get())
			}
			pasteAmountButton.Layout(gtx, buttonPasteAmount)
		})

	}
}

func amountEditor(gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			cs := gtx.Constraints
			theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 32, "fff4f4f4", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 30, "ffffffff", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				e := th.DuoUIeditor("DUO Amount")
				e.Font.Typeface = th.Font.Primary
				e.Font.Style = text.Italic
				e.Layout(gtx, amountLineEditor)
				for _, e := range amountLineEditor.Events(gtx) {
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
	}
}

func passwordEditor(gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			cs := gtx.Constraints
			theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 32, "fff4f4f4", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 30, "ffffffff", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				e := th.DuoUIeditor("DUO Amount")
				e.Font.Typeface = th.Font.Primary
				e.Font.Style = text.Italic
				e.Layout(gtx, passLineEditor)
				for _, e := range passLineEditor.Events(gtx) {
					if e, ok := e.(controller.SubmitEvent); ok {
						passPharse = e.Text
						passLineEditor.SetText("")

					}
				}
			})
		})
	}
}

func sendButton(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var sendButton theme.DuoUIbutton
			sendButton = th.DuoUIbutton(th.Font.Secondary, "SEND", th.Color.Light, th.Color.Dark, "", th.Color.Light, 16, 0, 128, 48, 0, 0)
			for buttonSend.Clicked(gtx) {
				rc.Dialog.Show = true
				rc.Dialog = &model.DuoUIdialog{
					Show: true,
					Ok:   rc.DuoSend(passPharse, address, amount),
					Close: func() {

					},
					Cancel: func() { rc.Dialog.Show = false },
					Title:  "Are you sure?",
					Text:   "Confirm ParallelCoin send",
				}
			}
			sendButton.Layout(gtx, buttonSend)
		})
	}
}

package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
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
	addressLineEditor = &gel.Editor{
		SingleLine: true,
	}
	amountLineEditor = &gel.Editor{
		SingleLine: true,
	}
	passLineEditor = &gel.Editor{
		SingleLine: true,
	}
	buttonPasteAddress = new(gel.Button)
	buttonPasteAmount  = new(gel.Button)
	buttonSend         = new(gel.Button)
)

func Send(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	return th.DuoUIpage("SEND", 10, func() {}, func() {}, sendBody(rc, gtx, th), func() {})
}

func sendBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(func() {
				cs := gtx.Constraints
				gelook.DuoUIdrawRectangle(gtx, cs.Width.Max, 180, th.Colors["Light"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				widgets := []func(){
					func() {
						layout.Flex{}.Layout(gtx,
							layout.Flexed(1, component.Editor(gtx, th, addressLineEditor, address, "DUO address", func(e gel.EditorEvent) {
								address = addressLineEditor.Text()
							})),
							layout.Rigid(component.Button(gtx, th, buttonPasteAddress, th.Fonts["Primary"], 12, th.Colors["ButtonText"], th.Colors["ButtonBg"], "PASTE ADDRESS", func() {
								addressLineEditor.SetText(clipboard.Get())
							})))
					},
					func() {
						layout.Flex{}.Layout(gtx,
							layout.Flexed(1, component.Editor(gtx, th, amountLineEditor, fmt.Sprint(amount), "DUO Amount", func(e gel.EditorEvent) {
								f, err := strconv.ParseFloat(amountLineEditor.Text(), 64)
								if err != nil {
									amount = f
									amountLineEditor.SetText("")
								}
							})),
							layout.Rigid(component.Button(gtx, th, buttonPasteAmount, th.Fonts["Primary"], 12, th.Colors["ButtonText"], th.Colors["ButtonBg"], "PASTE AMOUNT", func() {
								amountLineEditor.SetText(clipboard.Get())
							})))
					},
					func() {
						layout.Flex{}.Layout(gtx,
							layout.Rigid(component.Button(gtx, th, buttonSend, th.Fonts["Primary"], 12, th.Colors["ButtonText"], th.Colors["ButtonBg"], "SEND", func() {
								rc.Dialog.Show = true
								rc.Dialog = &model.DuoUIdialog{
									Show: true,
									Ok:   rc.DuoSend(passPharse, address, 1),
									Close: func() {

									},
									CustomField: func() {
										layout.Flex{}.Layout(gtx,
											layout.Flexed(1, component.Editor(gtx, th, passLineEditor, passPharse, "Enter your password", func(e gel.EditorEvent) {
												passPharse = passLineEditor.Text()
											})))
									},
									Cancel: func() { rc.Dialog.Show = false },
									Title:  "Are you sure?",
									Text:   "Confirm ParallelCoin send",
								}
							})))
					},
				}
				layautList.Layout(gtx, len(widgets), func(i int) {
					layout.UniformInset(unit.Dp(8)).Layout(gtx, widgets[i])
				})
			}))
	}
}

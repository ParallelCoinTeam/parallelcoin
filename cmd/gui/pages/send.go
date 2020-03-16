package pages

import (
	"fmt"
	log "github.com/p9c/logi"
	"strconv"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/pkg/gui/clipboard"
)

var (
	layautList = &layout.List{
		Axis: layout.Vertical,
	}
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
	sendStruct         = new(send)
)

type send struct {
	address    string
	amount     float64
	passPharse string
}

func Send(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	return th.DuoUIpage("SEND", 10, func() {}, func() {}, sendBody(rc, gtx, th), func() {})
}

func sendBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(func() {
				cs := gtx.Constraints
				gelook.DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, th.Colors["Dark"], [4]float32{0, 0, 0, 0},
					[4]float32{0, 0, 0, 0})
				widgets := []func(){
					func() {
						layout.Flex{}.Layout(gtx,
							layout.Flexed(1, component.Editor(gtx, th, addressLineEditor, "DUO address",
								func(e gel.EditorEvent) {
									sendStruct.address = addressLineEditor.Text()
								})),
							layout.Rigid(component.Button(gtx, th, buttonPasteAddress, th.Fonts["Primary"], 12,
								th.Colors["ButtonText"], th.Colors["ButtonBg"], "PASTE ADDRESS", func() {
									addressLineEditor.SetText(clipboard.Get())
								})))
					},
					func() {
						layout.Flex{}.Layout(gtx,
							layout.Flexed(1, component.Editor(gtx, th, amountLineEditor,
								"DUO Amount", func(e gel.EditorEvent) {
									f, err := strconv.ParseFloat(amountLineEditor.Text(), 64)
									if err != nil {
									}
									sendStruct.amount = f
								})),
							layout.Rigid(component.Button(gtx, th, buttonPasteAmount, th.Fonts["Primary"], 12,
								th.Colors["ButtonText"], th.Colors["ButtonBg"], "PASTE AMOUNT", func() {
									amountLineEditor.SetText(clipboard.Get())
								})))
					},
					func() {
						layout.Flex{}.Layout(gtx,
							layout.Rigid(component.Button(gtx, th, buttonSend, th.Fonts["Primary"], 12,
								th.Colors["ButtonText"], th.Colors["ButtonBg"], "SEND", func() {
									log.L.Info("passPharse:" + sendStruct.passPharse)
									log.L.Info("address" + sendStruct.address)
									log.L.Info("amount:" + fmt.Sprint(sendStruct.amount))
									rc.Dialog.Show = true
									rc.Dialog = &model.DuoUIdialog{
										Show:       true,
										Green:      rc.DuoSend(sendStruct.passPharse, sendStruct.address, 11),
										GreenLabel: "SEND",
										CustomField: func() {
											layout.Flex{}.Layout(gtx,
												layout.Flexed(1, component.Editor(gtx, th, passLineEditor, "Enter your password",
													func(e gel.EditorEvent) {
														sendStruct.passPharse = passLineEditor.Text()
													})))
										},
										Red:      func() { rc.Dialog.Show = false },
										RedLabel: "CANCEL",
										Title:    "Are you sure?",
										Text:     "Confirm ParallelCoin send",
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

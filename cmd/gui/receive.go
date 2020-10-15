package gwallet

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/wallet/counter"
	"github.com/p9c/pod/pkg/gui/wallet/dap/box"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

var (
	receiveClearBtn   = new(widget.Clickable)
	receiveLabelInput = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	receiveMessageInput = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	amountCounter = &counter.Counter{
		Value:        1,
		OperateValue: 1,
		From:         1,
		To:           999,
		CounterInput: &widget.Editor{
			Alignment:  text.Middle,
			SingleLine: true,
			Submit:     true,
		},
		//PageFunction:    w.PrikazaniElementSumaRacunica(),
		CounterIncrease: new(widget.Clickable),
		CounterDecrease: new(widget.Clickable),
		CounterReset:    new(widget.Clickable),
	}
)

func (g *GioWallet) GetReceive() {

}

func (g *GioWallet) receiveHeader() func(gtx C) D {
	return box.BoxPanel(g.ui.Theme, func(gtx C) D {
		return lyt.Format(gtx, "vflex(start,r(inset(0dp0dp0dp0dp,_)),r(inset(0dp0dp4dp0dp,_)),r(inset(0dp0dp4dp0dp,_)),r(inset(0dp0dp4dp0dp,_)),r(inset(0dp0dp0dp0dp,_))))",
			labeledRow(g.ui.Theme, "",
				func(gtx C) D {
					title := theme.Body(g.ui.Theme, "Use this form to request payments. All fields are optional.")
					title.Alignment = text.Start
					return title.Layout(gtx)
				}),
			labeledRow(g.ui.Theme, "Label:",
				func(gtx C) D {
					return lyt.Format(gtx, "hflex(middle,f(1,inset(0dp0dp0dp0dp,_)))",
						box.BoxEditor(g.ui.Theme, func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							e := material.Editor(g.ui.Theme.T, receiveLabelInput, "Enter a label for this address to add it to the list of used addresses")
							return e.Layout(gtx)
						}),
					)
				}),
			labeledRow(g.ui.Theme, "Amount:",
				func(gtx C) D {
					return lyt.Format(gtx, "hflex(middle,r(_),r(_),r(_),r(_))",
						counter.CounterSt(g.ui.Theme, amountCounter).Layout(g.ui.Theme, fmt.Sprint(amountCounter.Value)),
						func(gtx C) D {
							btn := material.IconButton(g.ui.Theme.T, connectionsBtn, g.ui.Theme.Icons["networkIcon"])
							btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
							btn.Size = unit.Dp(21)
							btn.Background = helper.HexARGB(g.ui.Theme.Colors["Secondary"])
							for connectionsBtn.Clicked() {
								//ui.N.CurrentPage = "Welcome"
							}
							return btn.Layout(gtx)
						},
						func(gtx C) D {
							btn := material.IconButton(g.ui.Theme.T, connectionsBtn, g.ui.Theme.Icons["networkIcon"])
							btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
							btn.Size = unit.Dp(21)
							btn.Background = helper.HexARGB(g.ui.Theme.Colors["Secondary"])
							for connectionsBtn.Clicked() {
								//ui.N.CurrentPage = "Welcome"
							}
							return btn.Layout(gtx)
						},
						func(gtx C) D {
							btn := material.IconButton(g.ui.Theme.T, connectionsBtn, g.ui.Theme.Icons["networkIcon"])
							btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
							btn.Size = unit.Dp(21)
							btn.Background = helper.HexARGB(g.ui.Theme.Colors["Secondary"])
							for connectionsBtn.Clicked() {
								//ui.N.CurrentPage = "Welcome"
							}
							return btn.Layout(gtx)
						},
						func(gtx C) D {
							btn := material.IconButton(g.ui.Theme.T, connectionsBtn, g.ui.Theme.Icons["networkIcon"])
							btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
							btn.Size = unit.Dp(21)
							btn.Background = helper.HexARGB(g.ui.Theme.Colors["Secondary"])
							for connectionsBtn.Clicked() {
								//ui.N.CurrentPage = "Welcome"
							}
							return btn.Layout(gtx)
						},
					)
				}),
			labeledRow(g.ui.Theme, "Message:",
				func(gtx C) D {
					return lyt.Format(gtx, "hflex(middle,f(1,inset(0dp0dp0dp0dp,_)))",
						box.BoxEditor(g.ui.Theme, func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							e := material.Editor(g.ui.Theme.T, receiveMessageInput, "Enter a label for this address to add it to the list of used addresses")
							return e.Layout(gtx)
						}),
					)
				}),
			labeledRow(g.ui.Theme, ":",
				func(gtx C) D {
					return lyt.Format(gtx, "hflex(middle,r(inset(0dp0dp0dp0dp,_)),r(inset(0dp0dp0dp0dp,_)))",
						func(gtx C) D {
							btn := material.IconButton(g.ui.Theme.T, connectionsBtn, g.ui.Theme.Icons["networkIcon"])
							btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
							btn.Size = unit.Dp(21)
							btn.Background = helper.HexARGB(g.ui.Theme.Colors["Secondary"])
							for connectionsBtn.Clicked() {
								//ui.N.CurrentPage = "Welcome"
							}
							return btn.Layout(gtx)
						},
						func(gtx C) D {
							btn := material.IconButton(g.ui.Theme.T, connectionsBtn, g.ui.Theme.Icons["networkIcon"])
							btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
							btn.Size = unit.Dp(21)
							btn.Background = helper.HexARGB(g.ui.Theme.Colors["Secondary"])
							for connectionsBtn.Clicked() {
								//ui.N.CurrentPage = "Welcome"
							}
							return btn.Layout(gtx)
						},
					)
				}),
		)
	})
}

func (g *GioWallet) receiveBody() func(gtx C) D {
	return box.BoxPanel(g.ui.Theme, func(gtx C) D {
		return lyt.Format(gtx, "max(vflex(start,r(inset(0dp0dp0dp0dp,_)),r(inset(0dp0dp0dp0dp,_)))))",
			func(gtx C) D {
				title := theme.Body(g.ui.Theme, "Requested payments history")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			func(gtx C) D {
				title := theme.Body(g.ui.Theme, "Requested payments history")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
		)
	})
}

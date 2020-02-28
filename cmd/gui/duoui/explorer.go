package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/clipboard"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

var (
	blocksList = &layout.List{
		Axis: layout.Vertical,
	}
	blockNumber = &controller.DuoUIcounter{
		Value:           20,
		OperateValue:    1,
		From:            0,
		To:              50,
		CounterIncrease: new(controller.Button),
		CounterDecrease: new(controller.Button),
		CounterReset:    new(controller.Button),
	}
	blockFrom = &controller.DuoUIcounter{
		Value:           20,
		OperateValue:    1,
		From:            0,
		To:              50,
		CounterIncrease: new(controller.Button),
		CounterDecrease: new(controller.Button),
		CounterReset:    new(controller.Button),
	}
)

func bodyExplorer(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		in := layout.UniformInset(unit.Dp(0))
		in.Layout(gtx, func() {
			blocksList.Layout(gtx, len(rc.Blocks), func(i int) {
				blockRow(rc, gtx, th, i)
			})
		})
	}
}

func headerExplorer(gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
			Axis:    layout.Horizontal,
		}.Layout(gtx,
			//layout.Rigid(ui.txsFilter()),
			layout.Flexed(0.5, func() {
				th.DuoUIcounter().Layout(gtx, blockNumber)
			}),
			layout.Flexed(0.5, func() {
				th.DuoUIcounter().Layout(gtx, blockFrom)
			}),
		)
	}
}

func blockRow(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, i int) {
	b := rc.Blocks[i]
	line(gtx, th.Color.Dark)()
	layout.Flex{
		Spacing: layout.SpaceBetween,
	}.Layout(gtx,
		layout.Rigid(func() {
			layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Flexed(0.6, func() {
					var linkButton theme.DuoUIbutton
					linkButton = th.DuoUIbutton(th.Font.Mono, fmt.Sprint(b.Height), th.Color.Light, th.Color.Dark, "", th.Color.Light, 16, 0, 60, 24, 0, 0)
					for b.Link.Clicked(gtx) {
						clipboard.Set(b.BlockHash)
					}
					linkButton.Layout(gtx, b.Link)
				}),
				layout.Rigid(func() {
					amount := th.H5(fmt.Sprintf("%0.8f", b.Amount))
					amount.Font.Typeface = th.Font.Primary
					amount.Color = theme.HexARGB(th.Color.Dark)
					amount.Alignment = text.End
					amount.Font.Variant = "Mono"
					amount.Font.Weight = text.Bold
					amount.Layout(gtx)
				}),
				layout.Rigid(func() {
					sat := th.Body1(fmt.Sprint(b.TxNum))
					sat.Font.Typeface = th.Font.Primary
					sat.Color = theme.HexARGB(th.Color.Dark)
					sat.Layout(gtx)
				}),
				layout.Rigid(func() {
					sat := th.Body1(fmt.Sprint(b.BlockHash))
					sat.Font.Typeface = th.Font.Mono
					sat.Color = theme.HexARGB(th.Color.Dark)
					sat.Layout(gtx)
				}),
				layout.Rigid(func() {
					l := th.Body2(b.Time)
					l.Font.Typeface = th.Font.Primary
					l.Color = theme.HexARGB(th.Color.Dark)
					l.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func() {
			sat := th.Body1(fmt.Sprintf("%0.8f", b.Amount))
			sat.Color = theme.HexARGB(th.Color.Dark)
			sat.Layout(gtx)
		}))
}
